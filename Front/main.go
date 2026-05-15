package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	controladores "docker_viewer/front/Controladores"
)

// ---- Embebe carpetas estáticas en el binario compilado ----
//go:embed Layout Login Vistas Public Componentes
var archivos embed.FS

// ---- Punto de entrada del servidor frontend ----
func main() {

	// ---- Carga y asocia cada plantilla HTML con su vista ----
	plantillas := &controladores.Plantillas{
		Login:        template.Must(template.ParseFS(archivos, "Login/login.html")),
		Inicio:       template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/inicio.html")),
		Contenedores: template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/contenedores.html", "Componentes/tarjeta.html")),
		Logs:         template.Must(template.ParseFS(archivos, "Layout/layout.html", "Vistas/logs.html")),
	}

	// ---- Expone la carpeta Public como archivos estáticos en /public/ ----
	publicFS, err := fs.Sub(archivos, "Public")
	if err != nil {
		log.Fatal(err)
	}

	// ---- Registra rutas HTTP y asocia cada una a su controlador ----
	rute := http.NewServeMux()
	rute.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.FS(publicFS))))

	controlador := controladores.Nuevo(plantillas)
	rute.HandleFunc("GET /login",            controlador.MostrarLogin)
	rute.HandleFunc("POST /login",           controlador.EntrarLogin)
	rute.HandleFunc("POST /logout",          controlador.CerrarSesion)
	rute.HandleFunc("GET /contenedores",     controlador.MostrarContenedores)
	rute.HandleFunc("GET /api/contenedores", controlador.APIContenedores) // ← endpoint para el fetch del botón actualizar
	rute.HandleFunc("GET /logs",             controlador.MostrarLogs)
	rute.HandleFunc("GET /",                 controlador.MostrarInicio)

	log.Println("Frontend en :10000 → http://localhost:10000")
	log.Fatal(http.ListenAndServe(":10000", rute))
}
