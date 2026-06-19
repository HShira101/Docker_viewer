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
	IP          string `json:"IP"`
	PublicPort  int    `json:"PublicPort"`
	PrivatePort int    `json:"PrivatePort"`
	Type        string `json:"Type"`
}

// ---- Representa un contenedor Docker deserializado desde el backend ----
type Contenedor struct {
	ID                string   `json:"id"`
	Nombre            string   `json:"nombre"`
	Imagen            string   `json:"imagen"`
	Estado            string   `json:"estado"`
	ComposeProject    string   `json:"compose_project"`
	Puertos           []Puerto `json:"puertos"`
	UltimaLogGuardado string   `json:"ultimo_log_guardado"`
}

// ---- Agrupa contenedores bajo un proyecto Compose (o sin proyecto) ----
type Grupo struct {
	Nombre       string
	Contenedores []Contenedor
}

// ---- Datos para la vista de contenedores: extiende DatosLayout ----
type DatosContenedores struct {
	DatosLayout
	Grupos []Grupo
}

// ---- Datos para la vista de logs: misma estructura que contenedores ----
type DatosLogs struct {
	DatosLayout
	Grupos []Grupo
}
