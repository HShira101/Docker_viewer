package contenedores

import "time"

// ---- Puerto publicado por un contenedor Docker ----
type Puerto struct {
	PublicPort  int    `json:"PublicPort"`
	PrivatePort int    `json:"PrivatePort"`
	Type        string `json:"Type"`
}

// ---- Representa un contenedor Docker almacenado en SQLite ----
type Contenedor struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Imagen         string    `json:"imagen"`
	Estado         string    `json:"estado"`
	UltimaConsulta time.Time `json:"ultima_consulta"`
	ComposeProject string    `json:"compose_project"`
	Puertos        []Puerto  `json:"puertos"`
}

// ---- Respuesta estándar para acciones sobre contenedores ----
type Respuesta struct {
	OK      bool   `json:"ok"`
	Mensaje string `json:"mensaje"`
}
