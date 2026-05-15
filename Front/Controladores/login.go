package controladores

import "net/http"

// ---- Renderiza la vista del formulario de login ----
func (c *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	c.plantillas.Login.ExecuteTemplate(w, "login.html", nil)
}

// ---- Procesa el formulario de login y redirige al inicio ----
func (c *Controlador) EntrarLogin(w http.ResponseWriter, r *http.Request) {
	// --- Recibe credenciales vía POST desde login.html ---
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---- Cierra la sesión y redirige al login ----
func (c *Controlador) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
