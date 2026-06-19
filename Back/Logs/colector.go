package logs

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	database "docker_viewer/back/Database"
)

// enProceso evita que el ticker y el streaming recolecten el mismo contenedor simultáneamente
var (
	enProceso   = map[string]bool{}
	muEnProceso sync.Mutex
)

// ---- IniciarTicker arranca el recolector en background, disparando cada 5 minutos ----
func IniciarTicker() {
	go func() {
		recolectarTodos() // primera recolección al arrancar

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			recolectarTodos()
		}
	}()
}

// ---- Recolecta logs de todos los contenedores de forma inmediata ----
func RecolectarAhora() {
	go recolectarTodos()
}

// ---- Construye el cliente HTTP para el socket Unix de Docker (con timeout) ----
func clienteDocker() *http.Client {
	return &http.Client{
		Transport: transporteDocker(),
		Timeout:   30 * time.Second,
	}
}

// ---- Construye el cliente HTTP para streaming (sin timeout, conexión larga) ----
func clienteDockerStream() *http.Client {
	return &http.Client{
		Transport: transporteDocker(),
	}
}

// ---- Transport compartido para el socket Unix de Docker ----
func transporteDocker() *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	}
}

// ---- Devuelve la URL base de VictoriaLogs desde env ----
func urlVlogs() string {
	if u := os.Getenv("VLOGS_URL"); u != "" {
		return u
	}
	return "http://victorialogs:9428"
}

// ---- Consulta SQLite y recolecta logs de todos los contenedores ----
func recolectarTodos() {
	rows, err := database.DB.Query(`SELECT id, nombre FROM contenedores`)
	if err != nil {
		log.Println("Ticker logs DB:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, nombre string
		rows.Scan(&id, &nombre)
		Recolectar(id, nombre)
	}
}

// ---- leerDesde lee ultimo_log_guardado de SQLite; si es NULL devuelve hace 1 hora ----
func leerDesde(id string) time.Time {
	var t time.Time
	err := database.DB.QueryRow(
		`SELECT ultimo_log_guardado FROM contenedores WHERE id = ?`, id,
	).Scan(&t)
	if err != nil || t.IsZero() {
		return time.Now().Add(-1 * time.Hour)
	}
	return t
}

// ---- guardarDesde persiste ultimo_log_guardado en SQLite ----
func guardarDesde(id string, t time.Time) {
	database.DB.Exec(
		`UPDATE contenedores SET ultimo_log_guardado = ? WHERE id = ?`, t, id,
	)
}

// ---- Recolectar obtiene logs de un contenedor desde Docker y los envía a VictoriaLogs ----
func Recolectar(id, nombre string) {
	// Si ya hay una recolección en curso para este contenedor, salir para evitar duplicados
	muEnProceso.Lock()
	if enProceso[id] {
		muEnProceso.Unlock()
		return
	}
	enProceso[id] = true
	muEnProceso.Unlock()

	desde := leerDesde(id)
	ahora := time.Now()

	defer func() {
		muEnProceso.Lock()
		enProceso[id] = false
		muEnProceso.Unlock()
	}()

	ruta := fmt.Sprintf(
		"/containers/%s/logs?stdout=true&stderr=true&timestamps=true&since=%d",
		id, desde.Unix(),
	)

	resp, err := clienteDocker().Get("http://localhost" + ruta)
	if err != nil {
		log.Printf("Logs Docker [%s]: %v", nombre, err)
		return
	}
	defer resp.Body.Close()

	lineas := parsearLogs(resp.Body, id, nombre)
	guardarDesde(id, ahora)

	if len(lineas) == 0 {
		return
	}

	if err := enviarAVictoriaLogs(lineas); err != nil {
		log.Printf("VictoriaLogs push [%s]: %v", nombre, err)
		return
	}

	log.Printf("Logs [%s]: %d líneas enviadas", nombre, len(lineas))
}

// ---- StreamLogs abre un stream SSE hacia el cliente con los logs en tiempo real ----
// Mientras el stream está activo, el ticker y RecolectarAhora omiten este contenedor
// porque el stream ya está pusheando a VictoriaLogs.
// Se detiene cuando el cliente cierra la conexión (ctx cancelado).
func StreamLogs(ctx context.Context, id, nombre string, w http.ResponseWriter) {
	// Marca el contenedor como en proceso para que el ticker y on-demand lo omitan
	muEnProceso.Lock()
	if enProceso[id] {
		muEnProceso.Unlock()
		http.Error(w, "ya hay una recolección en curso", http.StatusConflict)
		return
	}
	enProceso[id] = true
	muEnProceso.Unlock()

	defer func() {
		guardarDesde(id, time.Now()) // el ticker arranca desde aquí, no duplica lo que el stream ya envió
		muEnProceso.Lock()
		enProceso[id] = false
		muEnProceso.Unlock()
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming no soportado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ruta := fmt.Sprintf(
		"/containers/%s/logs?stdout=true&stderr=true&timestamps=true&follow=true",
		id,
	)

	// Usa el contexto del request: cuando el browser cierra la conexión, cancela el stream de Docker
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost"+ruta, nil)
	if err != nil {
		return
	}

	resp, err := clienteDockerStream().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	for {
		header := make([]byte, 8)
		if _, err := io.ReadFull(resp.Body, header); err != nil {
			break
		}

		streamTipo := header[0]
		frameSize := binary.BigEndian.Uint32(header[4:8])

		frame := make([]byte, frameSize)
		if _, err := io.ReadFull(resp.Body, frame); err != nil {
			break
		}

		stream := "stdout"
		if streamTipo == 2 {
			stream = "stderr"
		}

		var loteFrame []LineaLog

		scanner := bufio.NewScanner(bytes.NewReader(frame))
		for scanner.Scan() {
			linea := scanner.Text()
			if linea == "" {
				continue
			}

			tiempo := time.Now().UTC().Format(time.RFC3339Nano)
			mensaje := linea

			if idx := strings.Index(linea, " "); idx > 0 {
				candidato := linea[:idx]
				if strings.Contains(candidato, "T") && strings.Contains(candidato, "Z") {
					tiempo = candidato
					mensaje = linea[idx+1:]
				}
			}

			ll := LineaLog{
				Tiempo:  tiempo,
				Mensaje: mensaje,
				ID:      id,
				Nombre:  nombre,
				Stream:  stream,
			}

			// Envía la línea al browser como evento SSE
			data, _ := json.Marshal(ll)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			loteFrame = append(loteFrame, ll)
		}

		// Pushea el frame completo a VictoriaLogs en background sin bloquear el stream
		if len(loteFrame) > 0 {
			lote := loteFrame
			go enviarAVictoriaLogs(lote)
		}
	}
}

// ---- Envía un slice de líneas a VictoriaLogs en formato JSONL ----
func enviarAVictoriaLogs(lineas []LineaLog) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, l := range lineas {
		enc.Encode(l)
	}

	url := urlVlogs() + "/insert/jsonline?_stream_fields=container_id,container_name,stream"
	resp, err := http.Post(url, "application/stream+json", &buf)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ---- Parsea el stream multiplexado de Docker y devuelve líneas de log ----
func parsearLogs(r io.Reader, id, nombre string) []LineaLog {
	// Formato de frame Docker: [tipo(1)][pad(3)][size(4 big-endian)][datos(size)]
	// tipo: 1=stdout, 2=stderr. Sin TTY siempre usa este formato.
	var lineas []LineaLog

	for {
		header := make([]byte, 8)
		if _, err := io.ReadFull(r, header); err != nil {
			break
		}

		streamTipo := header[0]
		frameSize := binary.BigEndian.Uint32(header[4:8])

		frame := make([]byte, frameSize)
		if _, err := io.ReadFull(r, frame); err != nil {
			break
		}

		stream := "stdout"
		if streamTipo == 2 {
			stream = "stderr"
		}

		scanner := bufio.NewScanner(bytes.NewReader(frame))
		for scanner.Scan() {
			linea := scanner.Text()
			if linea == "" {
				continue
			}

			tiempo := time.Now().UTC().Format(time.RFC3339Nano)
			mensaje := linea

			if idx := strings.Index(linea, " "); idx > 0 {
				candidato := linea[:idx]
				if strings.Contains(candidato, "T") && strings.Contains(candidato, "Z") {
					tiempo = candidato
					mensaje = linea[idx+1:]
				}
			}

			lineas = append(lineas, LineaLog{
				Tiempo:  tiempo,
				Mensaje: mensaje,
				ID:      id,
				Nombre:  nombre,
				Stream:  stream,
			})
		}
	}

	return lineas
}
