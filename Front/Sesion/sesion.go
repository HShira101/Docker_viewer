package sesion

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"
)

const nombreCookie = "docker_viewer_session"

// ---- Devuelve la clave secreta para firmar cookies desde variable de entorno ----
func secreto() []byte {
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		s = "fallback-inseguro-solo-dev" // ← reemplazar con SESSION_SECRET en producción
	}
	return []byte(s)
}

// ---- Codifica el nombre en base64url y firma con HMAC-SHA256 ----
func firmar(nombre string) string {
	val := base64.RawURLEncoding.EncodeToString([]byte(nombre))
	mac := hmac.New(sha256.New, secreto())
	mac.Write([]byte(val))
	firma := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return val + "." + firma
}

// ---- Verifica la firma HMAC y devuelve el nombre de usuario ----
func verificarFirma(firmado string) (string, bool) {
	// --- Divide en valor + firma usando el punto como separador ---
	partes := strings.SplitN(firmado, ".", 2)
	if len(partes) != 2 {
		return "", false
	}
	val, firmaRecibida := partes[0], partes[1]

	mac := hmac.New(sha256.New, secreto())
	mac.Write([]byte(val))
	firmaEsperada := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(firmaEsperada), []byte(firmaRecibida)) {
		return "", false
	}

	nombre, err := base64.RawURLEncoding.DecodeString(val)
	if err != nil {
		return "", false
	}
	return string(nombre), true
}

// ---- Establece la cookie de sesión firmada con el nombre de usuario ----
func Establecer(w http.ResponseWriter, nombre string) {
	http.SetCookie(w, &http.Cookie{
		Name:     nombreCookie,
		Value:    firmar(nombre),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

// ---- Invalida la cookie de sesión ----
func Eliminar(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    nombreCookie,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
}

// ---- Devuelve el nombre de usuario si la sesión es válida, o "" si no ----
func ObtenerUsuario(r *http.Request) string {
	cookie, err := r.Cookie(nombreCookie)
	if err != nil {
		return ""
	}
	nombre, ok := verificarFirma(cookie.Value)
	if !ok {
		return ""
	}
	return nombre
}

// ---- Indica si la request tiene una sesión válida ----
func EstaAutenticado(r *http.Request) bool {
	return ObtenerUsuario(r) != ""
}
