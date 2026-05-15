package controladores

import "html/template"

// ---- Agrupa todas las plantillas HTML de la aplicación ----
type Plantillas struct {
	Login        *template.Template // ← vista de login
	Inicio       *template.Template // ← vista de inicio
	Contenedores *template.Template // ← vista de contenedores
	Logs         *template.Template // ← vista de logs
}

// ---- Controlador principal: accede a las plantillas para renderizar vistas ----
type Controlador struct {
	plantillas *Plantillas
}

// ---- Crea e inicializa el controlador con las plantillas recibidas ----
func Nuevo(plantillas *Plantillas) *Controlador {
	// --- Recibe plantillas como puntero a Plantillas ---
	return &Controlador{plantillas: plantillas}
}
