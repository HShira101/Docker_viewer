package controladores

// ---- Datos para la vista de login ----
type DatosLogin struct {
	Error string // ← mensaje de error a mostrar bajo el formulario (vacío si no hay error)
}

// ---- Datos que se pasan a todas las vistas que usan el layout ----
type DatosLayout struct {
	NombreUsuario string // ← nombre del usuario autenticado
	PaginaActual  string // ← marca el enlace activo en el sidebar
	CSS           string // ← nombre del archivo CSS extra de la vista (opcional)
}

// ---- Puerto publicado por un contenedor Docker ----
type Puerto struct {
	PublicPort  int    `json:"PublicPort"`
	PrivatePort int    `json:"PrivatePort"`
	Type        string `json:"Type"`
}

// ---- Representa un contenedor Docker deserializado desde el backend ----
type Contenedor struct {
	ID             string   `json:"id"`
	Nombre         string   `json:"nombre"`
	Imagen         string   `json:"imagen"`
	Estado         string   `json:"estado"`
	ComposeProject string   `json:"compose_project"`
	Puertos        []Puerto `json:"puertos"`
}

// ---- Agrupa contenedores bajo un proyecto Compose (o sin proyecto) ----
type Grupo struct {
	Nombre       string
	Contenedores []Contenedor
}

// ---- Datos para la vista de contenedores: extiende DatosLayout ----
type DatosContenedores struct {
	DatosLayout        // ← hereda NombreUsuario, PaginaActual y CSS
	Grupos      []Grupo // ← contenedores agrupados por compose_project
}
