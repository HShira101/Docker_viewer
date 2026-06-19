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

// ultimaRecoleccion guarda la hora de la última recolección por container ID
var (
	ultimaRecoleccion   = map[string]time.Time{}
	muUltimaRecoleccion sync.Mutex
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

// ---- Construye el cliente HTTP para el socket Unix de Docker ----
func clienteDocker() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
			},
		},
		Timeout: 30 * time.Second,
	}
}

// ---- Devuelve la URL base de VictoriaLogs desde env ----
func urlVlogs() string {
	if u := os.Getenv("VLOGS_URL"); u != "" {
		return u
	}
	return "http://victorialogs:9428"
}

// ---- Consulta SQLite y recolecta logs de todos los contenedores (running y detenidos) ----
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

// ---- Recolectar obtiene logs de un contenedor desde Docker y los envía a VictoriaLogs ----
func Recolectar(id, nombre string) {
	muUltimaRecoleccion.Lock()
	desde, existe := ultimaRecoleccion[id]
	if !existe {
		desde = time.Now().Add(-1 * time.Hour) // primera vez: última hora
	}
	ahora := time.Now()
	muUltimaRecoleccion.Unlock()

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

	muUltimaRecoleccion.Lock()
	ultimaRecoleccion[id] = ahora
	muUltimaRecoleccion.Unlock()

	if len(lineas) == 0 {
		return
	}

	if err := enviarAVictoriaLogs(lineas); err != nil {
		log.Printf("VictoriaLogs push [%s]: %v", nombre, err)
		return
	}

	log.Printf("Logs [%s]: %d líneas enviadas", nombre, len(lineas))
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

			// Docker con timestamps: "2024-01-01T00:00:00.000000000Z mensaje"
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
