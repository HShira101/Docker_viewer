package controladores

// ---- Datos que se pasan a todas las vistas que usan el layout ----
type DatosLayout struct {
	NombreUsuario string // ← nombre del usuario autenticado
	PaginaActual  string // ← marca el enlace activo en el sidebar
	CSS           string // ← nombre del archivo CSS extra de la vista (opcional)
}

// ---- Representa un contenedor Docker deserializado desde el backend ----
type Contenedor struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Imagen string `json:"imagen"`
	Estado string `json:"estado"`
}

// ---- Datos para la vista de contenedores: extiende DatosLayout ----
type DatosContenedores struct {
	DatosLayout              // ← hereda NombreUsuario, PaginaActual y CSS
	Contenedores []Contenedor // ← lista de contenedores a renderizar
}
