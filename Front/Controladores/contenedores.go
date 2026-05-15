package controladores

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func obtenerContenedores() []Contenedor {
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}
	resp, err := http.Get(backURL + "/api/contenedores")
	if err != nil {
		log.Println("Back API:", err)
		return nil
	}
	defer resp.Body.Close()

	var lista []Contenedor
	if err := json.NewDecoder(resp.Body).Decode(&lista); err != nil {
		log.Println("Decode contenedores:", err)
		return nil
	}
	return lista
}

func (c *Controlador) MostrarContenedores(w http.ResponseWriter, r *http.Request) {
	datos := DatosContenedores{
		DatosLayout:  DatosLayout{NombreUsuario: "Shira", PaginaActual: "contenedores", CSS: "contenedores.css"},
		Contenedores: obtenerContenedores(),
	}
	c.plantillas.Contenedores.ExecuteTemplate(w, "layout.html", datos)
}

func (c *Controlador) APIContenedores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obtenerContenedores())
}
