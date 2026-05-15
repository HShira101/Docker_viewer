package controladores

type DatosLayout struct {
	NombreUsuario string
	PaginaActual  string
	CSS           string
}

type Contenedor struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Imagen string `json:"imagen"`
	Estado string `json:"estado"`
}

type DatosContenedores struct {
	DatosLayout
	Contenedores []Contenedor
}
