---
name: project-branches
description: Estrategia de ramas y variantes del proyecto Docker Viewer
metadata:
  type: project
---

El proyecto tiene 3 ramas principales con distintos niveles de stack:

- **main-lite** — versión actual. Go + SQLite. Sin Victoria Metrics, sin Victoria Logs, sin MariaDB.
- **main-vmetrics** — añade Victoria Metrics y Victoria Logs para métricas y logs.
- **main-db** — incluye MariaDB para gestión de usuarios + Victoria Metrics.

**Why:** separar complejidad de infraestructura; el usuario quiere poder elegir qué servicios montar.

**How to apply:** al sugerir features que requieran VM/VLogs/MariaDB, aclarar en qué rama aplican. No asumir que el stack completo está disponible en main-lite.

---

Idea futura registrada: en el docker-compose, hacer los servicios de VM/VLogs/MariaDB opcionales (comentados o con profiles). Al primer arranque, mostrar una pantalla de configuración tipo "instalación" donde el usuario indique puerto de VM, si usa logs nativos o VLogs, etc.
