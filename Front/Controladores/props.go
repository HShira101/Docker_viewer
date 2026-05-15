package controladores

type DatosLayout struct {
	NombreUsuario string
	PaginaActual  string
	CSS           string
}

type Contenedor struct {
	ID     string
	Nombre string
	Imagen string
	Estado string
}

type DatosContenedores struct {
	DatosLayout
	Contenedores []Contenedor
}
