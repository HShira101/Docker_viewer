package contenedores

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	database "docker_viewer/back/Database"
)

// ---- Hace una petición GET a la API de Docker vía Unix socket ----
func dockerGet(ruta string, destino any) error {
	// --- Recibe ruta como string con el endpoint Docker y destino como puntero a struct ---
	resp, err := clienteDocker().Get("http://localhost" + ruta)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(destino)
}

// ---- Hace una petición POST a la API de Docker vía Unix socket y devuelve el código HTTP ----
func dockerPost(ruta string) (int, error) {
	// --- Recibe ruta como string con el endpoint Docker ---
	resp, err := clienteDocker().Post("http://localhost"+ruta, "application/json", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// ---- Construye el cliente HTTP configurado para conectar al socket Unix de Docker ----
func clienteDocker() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock") // ← socket Unix de Docker
			},
		},
	}
}

// ---- Consulta Docker y actualiza o inserta los contenedores en SQLite ----
func Actualizar() {
	var raw []struct {
		Id     string            `json:"Id"`
		Names  []string          `json:"Names"`
		Image  string            `json:"Image"`
		State  string            `json:"State"`
		Labels map[string]string `json:"Labels"`
		Ports  []struct {
			IP          string `json:"IP"`
			PublicPort  int    `json:"PublicPort"`
			PrivatePort int    `json:"PrivatePort"`
			Type        string `json:"Type"`
		} `json:"Ports"`
	}

	if err := dockerGet("/containers/json?all=true", &raw); err != nil {
		log.Println("Docker socket:", err)
		return
	}

	ahora := time.Now()
	for _, r := range raw {
		nombre := r.Id[:12]
		if len(r.Names) > 0 {
			nombre = r.Names[0][1:]
		}

		composeProject := r.Labels["com.docker.compose.project"]

		// serializa solo los puertos con mapeo al host (PublicPort > 0)
		puertos := []Puerto{}
		for _, p := range r.Ports {
			if p.PublicPort > 0 {
				puertos = append(puertos, Puerto{
					IP:          p.IP,
					PublicPort:  p.PublicPort,
					PrivatePort: p.PrivatePort,
					Type:        p.Type,
				})
			}
		}
		puertosJSON, _ := json.Marshal(puertos)

		database.DB.Exec(`
			INSERT INTO contenedores_running (id, nombre, imagen, estado, ultima_consulta, compose_project, puertos)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				nombre          = excluded.nombre,
				imagen          = excluded.imagen,
				estado          = excluded.estado,
				ultima_consulta = excluded.ultima_consulta,
				compose_project = excluded.compose_project,
				puertos         = excluded.puertos
		`, r.Id, nombre, r.Image, r.State, ahora, composeProject, string(puertosJSON))
	}
}
