# Intelfy API

Guia rapida para trabajar con este proyecto.

## Requisitos

- Go instalado
- Docker y Docker Compose
- `sqlc` instalado y disponible en tu `PATH`

Si instalaste `sqlc` con `go install`, normalmente queda en `~/go/bin`.

## Variables de entorno

El proyecto usa estas variables:

- `DB_URL` para conectarse a PostgreSQL
- `PORT` para el puerto del servidor HTTP
## Requisitos

- Go instalado
- Docker y Docker Compose
- `sqlc` instalado
- `migrate` (golang-migrate) instalado

## Comandos utiles

### Base de Datos

- `make db-up`: Levanta PostgreSQL en Docker.
- `make db-down`: Detiene el contenedor de base de datos.

### Migraciones (Database Migrations)

El proyecto utiliza `golang-migrate` para gestionar los cambios en la base de datos de forma secuencial.

- `make migrate-create`: Crea un nuevo par de archivos (up/down) para una nueva migración. Te pedirá el nombre de la migración.
- `make migrate-up`: Aplica **todas** las migraciones pendientes en la carpeta `db/migrations/`.
- `make migrate-down`: Revierte la **última** migración aplicada.
- `make migrate-version`: Muestra la versión actual de la base de datos.
- `make migrate-reset`: Borra todo el esquema y aplica todas las migraciones desde cero.

Note: the songs table uses the column `duration_seconds` (not `duration_ms`).

### Generar codigo con sqlc

```bash
make sqlc-generate
```

Genera el código Go en `internal/repository/` basado en el schema y las queries.

### Ejecutar la API

```bash
go run ./cmd/api
```

## Flujo de trabajo recomendado

1. Levanta PostgreSQL con `make db-up`.
2. Si es la primera vez o hay cambios, aplica las migraciones con `make migrate-up`.
3. Si cambiaste las queries en `db/queries/`, genera el código con `make sqlc-generate`.
4. Ejecuta la API con `go run ./cmd/api`.

## Notas Importantes

- **Estructura SQLC:** Las queries deben estar en `db/queries/` y el schema en `db/schema.sql`.
- **Orden de Migraciones:** Las migraciones se ejecutan en orden numérico (000001, 000002, etc.). No modifiques archivos de migración ya aplicados; en su lugar, crea una nueva migración con `make migrate-create`.

