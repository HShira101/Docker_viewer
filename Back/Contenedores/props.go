package contenedores

import "time"

type Contenedor struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Imagen         string    `json:"imagen"`
	Estado         string    `json:"estado"`
	UltimaConsulta time.Time `json:"ultima_consulta"`
}
