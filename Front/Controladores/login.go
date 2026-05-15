package controladores

import "net/http"

func (c *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	c.plantillas.Login.ExecuteTemplate(w, "login.html", nil)
}

func (c *Controlador) EntrarLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controlador) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
