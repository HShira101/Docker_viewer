package autenticacion

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	database "docker_viewer/back/Database"
)

type solicitudLogin struct {
	Nombre   string `json:"nombre"`
	Password string `json:"password"`
}

type respuestaLogin struct {
	OK      bool   `json:"ok"`
	Nombre  string `json:"nombre,omitempty"`
	Mensaje string `json:"mensaje,omitempty"`
}

// ---- Valida credenciales contra la tabla usuarios y responde con JSON ----
func Login(w http.ResponseWriter, r *http.Request) {
	// --- Recibe JSON con nombre y password desde el frontend ---
	var sol solicitudLogin
	if err := json.NewDecoder(r.Body).Decode(&sol); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuestaLogin{OK: false, Mensaje: "Solicitud inválida"})
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(sol.Password))) // ← compara hash, nunca la contraseña plana

	var nombreDB string
	err := database.DB.QueryRow(
		`SELECT nombre FROM usuarios WHERE nombre = ? AND password = ?`,
		sol.Nombre, hash,
	).Scan(&nombreDB)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(respuestaLogin{OK: false, Mensaje: "Credenciales incorrectas"})
		return
	}

	json.NewEncoder(w).Encode(respuestaLogin{OK: true, Nombre: nombreDB})
}
