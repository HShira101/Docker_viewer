package controladores

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
)

type Plantillas struct {
	Login        *template.Template
	Inicio       *template.Template
	Contenedores *template.Template
	Logs         *template.Template
}

type Controlador struct {
	plantillas *Plantillas
}

func Nuevo(plantillas *Plantillas) *Controlador {
	return &Controlador{plantillas: plantillas}
}

// dockerGet hace una petición GET a la API de Docker vía Unix socket.
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

func obtenerContenedores() []Contenedor {
	var raw []struct {
		Id    string   `json:"Id"`
		Names []string `json:"Names"`
		Image string   `json:"Image"`
		State string   `json:"State"`
	}

	if err := dockerGet("/containers/json?all=true", &raw); err != nil {
		log.Println("Docker socket:", err)
		return nil
	}

	lista := make([]Contenedor, 0, len(raw))
	for _, r := range raw {
		nombre := r.Id[:12]
		if len(r.Names) > 0 {
			nombre = r.Names[0][1:] // Docker prefija con "/"
		}
		lista = append(lista, Contenedor{
			ID:     r.Id[:12],
			Nombre: nombre,
			Imagen: r.Image,
			Estado: r.State,
		})
	}
	return lista
}
