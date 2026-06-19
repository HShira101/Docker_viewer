package contenedores

import "time"

// ---- Representa un contenedor Docker almacenado en SQLite ----
type Contenedor struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Imagen         string    `json:"imagen"`
	Estado         string    `json:"estado"`
	UltimaConsulta time.Time `json:"ultima_consulta"`
	ComposeProject string    `json:"compose_project"` // ← valor de com.docker.compose.project, vacío si no es Compose
}

// ---- Respuesta estándar para acciones sobre contenedores ----
type Respuesta struct {
	OK      bool   `json:"ok"`
	Mensaje string `json:"mensaje"`
}
