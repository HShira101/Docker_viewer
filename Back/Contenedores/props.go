package contenedores

import "time"

// ---- Representa un contenedor Docker almacenado en SQLite ----
type Contenedor struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Imagen         string    `json:"imagen"`
	Estado         string    `json:"estado"`
	UltimaConsulta time.Time `json:"ultima_consulta"` // ← fecha/hora de la última sincronización con Docker
}
