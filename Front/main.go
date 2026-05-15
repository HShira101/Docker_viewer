package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	controladores "docker_viewer/front/Controladores"
)

//go:embed Layout Login Vistas Public Componentes
var archivos embed.FS

func main() {

	plantillas := &controladores.Plantillas{
		// Inicialización de las plantillas
		Login:        template.Must(template.ParseFS(archivos, "Login/login.html")),
		Inicio:       template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/inicio.html")),
		Contenedores: template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/contenedores.html", "Componentes/tarjeta.html")),
		Logs:         template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/logs.html")),
	}

	publicFS, err := fs.Sub(archivos, "Public")
	// Manejo de errores al acceder al sistema de archivos
	if err != nil {
		log.Fatal(err)
	}

	// Configuración del enrutador y manejo de rutas
	rute := http.NewServeMux()
	rute.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.FS(publicFS))))

	controlador := controladores.Nuevo(plantillas)
	rute.HandleFunc("GET /login", controlador.MostrarLogin)
	rute.HandleFunc("POST /login", controlador.EntrarLogin)
	rute.HandleFunc("POST /logout", controlador.CerrarSesion)
	rute.HandleFunc("GET /contenedores", controlador.MostrarContenedores)
	rute.HandleFunc("GET /logs", controlador.MostrarLogs)
	rute.HandleFunc("GET /", controlador.MostrarInicio)

	// Logs de inicio del servidor
	log.Println("Frontend en :10000 → http://localhost:10000")
	log.Fatal(http.ListenAndServe(":10000", rute))
}
