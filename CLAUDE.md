#### Contexto ############################

Panel de gestión de contenedores Docker construido con Go + SQLite,
desplegado dentro de un contenedor Docker en una Raspberry Pi.
Se accede desde la red local via IP:10000 para ver qué contenedores están corriendo, poder prenderlos o apagarlos,
y adicionalmente poder ver uso de consumo de memoria de la raspberry y tambien acceso fácil a los logs. Como fase
experimental y no se considerará como obligatorio, posibilidad de ejecutar comandos bash en los contenedores y 
comandos docker.

Proyecto
GITHUB: https://github.com/HShira101/Docker_viewer
HTTPS: https://github.com/HShira101/Docker_viewer.git

Desarrollador
Nombre: Javier Navarro
Apodo DEV: Shira

#### Herramientas ########################

- GO (front y backend)
- SQLite (`database/database.sqlite`)
- Vector, para colectar metricas y logs.
- Docker (imagen propia desde `Dockerfile`)
- Victoria metrics (DB para metricas y logs)
- Victoria Logs (Guardar logs)

#### Arrancar ############################

```bash
# Desarrollo (Air hot reload — cambios en .go/.html/.css se aplican solos)
docker compose -f docker-compose.dev.yml up --build   # primera vez
docker compose -f docker-compose.dev.yml up           # arranque normal

# Producción
docker compose up --build       # primera vez o tras cambios en Dockerfile/código
docker compose up               # arranque normal sin cambios
docker compose down -v          # parar y borrar volúmenes (rebuild limpio)
```

App accesible en `http://localhost:10000`.

#### Estructura del proyecto ##############

Carpeta proyecto
|-> Front.
|-> Back.
|-> Sqlite

#### Front ################################

- Lenguaje: GO
- Estructura:
Front
  |-> Main.go (rutas)
  |-> Login
  |-> Layout
  |-> Vistas (carpeta para vistas y html)
  |-> Public (carpetas para css y js)
  |-> Controladores (Carpeta para logica)

#### Back ################################

- Lenguaje: GO
- Estructura:
Back
  |-> Main.go
  |-> Autenticacion (lógica de autenticación)
  |-> Sesion (lógica de sesion)
  |-> Metricas (lógica para otorgar muestras)
  |-> Logs (lógica para otorgar logs)
  |-> Seguridad (lógica para seguridad)

#### Skills #############################

#### Skills - Documentación #############
Te daré la estructura de como documentar, puedes usar cualquier forma de código js, html, etc. {-- inicia comentario, --} termina.
recuerda adaptar el comentario al lenguaje. 

# Comentar funcionalidad ()

{---- Funcionalidad (agrega texto, cambia color, etc.) ----}
def función
|  {--- Recibe x desde función y como (entero, lista, tupla) ---}
|
|
end función

# Comentar linea

int x = 1; {-- ← Hace que la variable sea uno --}

# Comentar etiqueta HTML

{---- Barra de navegación (puede se cualquier elemento) ----}
<div> 
|
|   
|
</div>
{---- Fin Barra de navegación (puede se cualquier elemento) ----}

#### Skils - Documentación PDF #########
Cuando te pida documentación pdf debes darme un pdf con lo siguiente:
→ Portada
| Nombre del proyecto
| Rama github
| Versión del proyecto
→ Indice
| Indice de títulos PDF y subsecciones
→ Cuerpo
| Sección enumerada y subsecciones 
→ Detalle de versión
| Rama de desarrollo
| Especificaciones ténicas (Que softwares se usan)
| Cambios con versión anterior
→ Instalación (seccion windows y linux)
| Requerimintos (saparado por windows y linux)
| Como se instala (Separados por windows y linux)
→ Funcionalidad
| Como funciona
→ Creditos
| Creditos y referencias a mi github
| Creditos a softwares usados (Incluyeté), si es posible con links a las páginas de los softwares

#### Notas ##############################
- Este proyecto es principalmente de aprendizaje.
- Principalmente usa GO
- Uso casero
- IMPORTANTE: CSS para estilos, todo html debe usar los estilos en css. Funciones deben ser REUTILIZADOS, para reutilizar código.

#### Instrucciones para Claude ########## 

- Solo documentar código cuando el usuario lo pida explícitamente.
- Las memorias y preferencias van en `.claude/memory/` dentro del proyecto.
  El índice es `.claude/memory/MEMORY.md`. Usarlo siempre al inicio de conversación
  y actualizarlo cuando el usuario indique algo que deba recordarse entre sesiones.

#### Fases de desarollo ################

- Fase 1: [COMPLETADA] creación de carpetas, dockerfile, Front básico (Fase de armar esqueleto).
  |-> Contenido del front: Login (no funcional, simplemente entrar al hacer click), layout, vista de inicio de ejemplo (No funcional).
  |-> Layout
        |-> Navegador lateral izquierdo con resumen de nombre de usuario autenticado, Resumen Uso de recursos, Boton contenedores y Logs.

- Fase 2: [COMPLETADA] Vista de contenedores funcional.
  |-> 1: Lógica de autenticación y vista de login funcional (SHA256, cookie de sesión).
  |-> 2: Lógica de contenedores
          |-> Consulta Docker API via Unix socket (/var/run/docker.sock).
          |-> Upsert de contenedores en SQLite (id, nombre, imagen, estado, compose_project, puertos).
          |-> Sistema de migraciones en Back/Database/Migrations/.
          |-> Backend entrega JSON al frontend.
          |-> Frontend agrupa contenedores por proyecto Compose.
          |-> Tarjetas muestran: estado, imagen, puertos publicados (IPv4/IPv6 con link al servicio).

- Fase 3: [EN PROGRESO] Vista de logs funcional.
  |-> 1: [COMPLETADO] Colector de logs desde Docker API con ticker cada 5min + recolección bajo demanda.
          |-> Paquete Back/Logs/ con colector.go (ticker, Recolectar, StreamLogs) y handler.go.
          |-> enProceso mutex evita duplicados entre ticker y stream.
          |-> ultimo_log_guardado persiste en SQLite (migración 005).
          |-> Push a VictoriaLogs vía /insert/jsonline (JSONL).
  |-> 2: [COMPLETADO] Vista de logs en el frontend.
          |-> Componentes: tarjeta-log.html, grupo-log.html.
          |-> Vista logs.html agrupa contenedores por compose_project igual que contenedores.
          |-> JS carga últimos 10 logs por tarjeta desde VictoriaLogs al abrir.
          |-> Botón Actualizar: recolecta primero (síncrono), luego re-consulta VictoriaLogs.
  |-> 3: [COMPLETADO] Integración con Victoria Logs directamente (sin Vector).
          |-> Servicio victorialogs en docker-compose.yml, retención 7 días.
          |-> Query con LogsQL {container_id="X"} con URL encoding correcto.
  |-> 4: [PENDIENTE] Botón "Logs" en tarjeta de contenedor abre la vista filtrada por contenedor.

- Fase 4: Vista de Usuarios funcional.
- Fase 5: Vista de comandos.

Con el tiempo iré detallando las fases de desarrollo.