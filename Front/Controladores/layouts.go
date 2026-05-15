package controladores

import "net/http"

func (c *Controlador) MostrarInicio(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "inicio"}
	c.plantillas.Inicio.ExecuteTemplate(w, "layout.html", datos)
}

func (c *Controlador) MostrarContenedores(w http.ResponseWriter, r *http.Request) {
	datos := DatosContenedores{
		DatosLayout:  DatosLayout{NombreUsuario: "Shira", PaginaActual: "contenedores", CSS: "contenedores.css"},
		Contenedores: obtenerContenedores(),
	}
	c.plantillas.Contenedores.ExecuteTemplate(w, "layout.html", datos)
}

func (c *Controlador) MostrarLogs(w http.ResponseWriter, r *http.Request) {
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "logs"}
	c.plantillas.Logs.ExecuteTemplate(w, "layout.html", datos)
}
