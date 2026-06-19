package logs

// ---- Línea de log de un contenedor, formato compatible con VictoriaLogs jsonline ----
type LineaLog struct {
	Tiempo  string `json:"_time"`
	Mensaje string `json:"_msg"`
	ID      string `json:"container_id"`
	Nombre  string `json:"container_name"`
	Stream  string `json:"stream"` // stdout | stderr
}
