package controladores

import "html/template"

type Plantillas struct {
	Login        *template.Template
	Inicio       *template.Template
	Contenedores *template.Template
	Logs         *template.Template
}

type Controlador struct {
	plantillas *Plantillas
}

func Nuevo(plantillas *Plantillas) *Controlador {
	return &Controlador{plantillas: plantillas}
}
