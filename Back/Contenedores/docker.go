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

func dockerGet(ruta string, destino any) error {
	cliente := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
			},
		},
	}
	resp, err := cliente.Get("http://localhost" + ruta)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(destino)
}

func Actualizar() {
	var raw []struct {
		Id    string   `json:"Id"`
		Names []string `json:"Names"`
		Image string   `json:"Image"`
		State string   `json:"State"`
	}
	if err := dockerGet("/containers/json?all=true", &raw); err != nil {
		log.Println("Docker socket:", err)
		return
	}

	ahora := time.Now()
	for _, r := range raw {
		nombre := r.Id[:12]
		if len(r.Names) > 0 {
			nombre = r.Names[0][1:] // Docker prefija con "/"
		}
		database.DB.Exec(`
			INSERT INTO contenedores_running (id, nombre, imagen, estado, ultima_consulta)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				nombre          = excluded.nombre,
				imagen          = excluded.imagen,
				estado          = excluded.estado,
				ultima_consulta = excluded.ultima_consulta
		`, r.Id, nombre, r.Image, r.State, ahora)
	}
}
