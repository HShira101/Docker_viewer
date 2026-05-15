package controladores

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	sesion "docker_viewer/front/Sesion"
)

type respuestaAuth struct {
	OK      bool   `json:"ok"`
	Nombre  string `json:"nombre"`
	Mensaje string `json:"mensaje"`
}

// ---- Renderiza la vista del formulario de login ----
func (c *Controlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	// --- Redirige al inicio si ya hay sesión activa ---
	if sesion.EstaAutenticado(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c.plantillas.Login.ExecuteTemplate(w, "login.html", DatosLogin{})
}

// ---- Procesa el formulario de login: valida con el back y crea sesión ----
func (c *Controlador) EntrarLogin(w http.ResponseWriter, r *http.Request) {
	// --- Recibe credenciales vía POST desde login.html ---
	nombre := strings.TrimSpace(r.FormValue("identificador"))
	password := r.FormValue("password")

	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}

	// --- Serializa las credenciales y las envía al backend ---
	cuerpo, _ := json.Marshal(map[string]string{"nombre": nombre, "password": password})
	resp, err := http.Post(backURL+"/api/auth/login", "application/json", bytes.NewReader(cuerpo))
	if err != nil {
		c.plantillas.Login.ExecuteTemplate(w, "login.html", DatosLogin{Error: "Backend no disponible"})
		return
	}
	defer resp.Body.Close()

	var resultado respuestaAuth
	json.NewDecoder(resp.Body).Decode(&resultado)

	if !resultado.OK {
		c.plantillas.Login.ExecuteTemplate(w, "login.html", DatosLogin{Error: "Usuario o contraseña incorrectos"})
		return
	}

	sesion.Establecer(w, resultado.Nombre) // ← crea la cookie firmada con el nombre real del usuario
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---- Cierra la sesión eliminando la cookie y redirige al login ----
func (c *Controlador) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	sesion.Eliminar(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
