package controladores

import (
	"html/template"
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

type DatosLayout struct {
	NombreUsuario string
	PaginaActual  string
}

func (c *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	c.plantillas.Login.ExecuteTemplate(w, "login.html", nil)
}

func (c *Controlador) EntrarLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controlador) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *Controlador) MostrarInicio(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "inicio"}
	c.plantillas.Inicio.ExecuteTemplate(w, "layout.html", datos)
}

func (c *Controlador) MostrarContenedores(w http.ResponseWriter, r *http.Request) {
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "contenedores"}
	c.plantillas.Contenedores.ExecuteTemplate(w, "layout.html", datos)
}

func (c *Controlador) MostrarLogs(w http.ResponseWriter, r *http.Request) {
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "logs"}
	c.plantillas.Logs.ExecuteTemplate(w, "layout.html", datos)
}
