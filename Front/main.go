package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	controladores "docker_viewer/front/Controladores"
	sesion "docker_viewer/front/Sesion"
)

// ---- Embebe carpetas estáticas en el binario compilado ----
//go:embed Layout Login Vistas Public Componentes
var archivos embed.FS

// ---- Middleware que redirige al login si no hay sesión activa ----
func proteger(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sesion.EstaAutenticado(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		siguiente(w, r)
	}
}

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

	// ---- Rutas públicas: login y assets ----
	rute.HandleFunc("GET /login",  controlador.MostrarLogin)
	rute.HandleFunc("POST /login", controlador.EntrarLogin)

	// ---- Rutas protegidas: requieren sesión activa ----
	rute.HandleFunc("POST /logout",                        proteger(controlador.CerrarSesion))
	rute.HandleFunc("GET /contenedores",                   proteger(controlador.MostrarContenedores))
	rute.HandleFunc("GET /api/contenedores",               proteger(controlador.APIContenedores))
	rute.HandleFunc("POST /api/contenedores/{id}/iniciar", proteger(controlador.APIIniciar))
	rute.HandleFunc("POST /api/contenedores/{id}/detener", proteger(controlador.APIDetener))
	rute.HandleFunc("GET /logs",                           proteger(controlador.MostrarLogs))
	rute.HandleFunc("GET /",                               proteger(controlador.MostrarInicio))

	log.Println("Frontend en :10000 → http://localhost:10000")
	log.Fatal(http.ListenAndServe(":10000", rute))
}
