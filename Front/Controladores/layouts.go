package controladores

import "net/http"

// ---- Renderiza la vista de inicio ----
func (c *Controlador) MostrarInicio(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound) // ← redirige cualquier ruta no registrada al inicio
		return
	}
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "inicio"}
	c.plantillas.Inicio.ExecuteTemplate(w, "layout.html", datos)
}

// ---- Renderiza la vista de logs ----
func (c *Controlador) MostrarLogs(w http.ResponseWriter, r *http.Request) {
	datos := DatosLayout{NombreUsuario: "Shira", PaginaActual: "logs"}
	c.plantillas.Logs.ExecuteTemplate(w, "layout.html", datos)
}
