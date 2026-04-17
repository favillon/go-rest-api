# Backend Productos

API REST en Go usando Gin, GORM y PostgreSQL.

## Configuracion

La aplicacion carga variables de entorno desde el archivo `.env`.

Variables usadas:

- `PORT_APP`: puerto donde corre la API (default: `8082`)
- `POSTGRES_HOST`: host de PostgreSQL (default: `localhost`)
- `POSTGRES_PORT`: puerto de PostgreSQL (default: `5432`)
- `POSTGRES_DB`: nombre de la base de datos
- `POSTGRES_USER`: usuario de PostgreSQL
- `POSTGRES_PASSWORD`: password de PostgreSQL

Ejemplo (`.env.example`):

```dotenv
PORT_APP=8082
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=restapi_go_db
POSTGRES_USER=user_go
POSTGRES_PASSWORD=change_me
```

## Base de datos (Docker)

El archivo `docker-compose.yml` levanta PostgreSQL 18 y mapea el puerto:

- `5432:${POSTGRES_PORT}`

Comandos utiles:

```bash
docker compose up -d postgres
docker compose ps
docker compose down
```

## Ejecucion de la API

### Modo normal

```bash
go run .
```

### Modo desarrollo con recarga (Air)

```bash
air
```

La API queda disponible en:

```bash
http://localhost:${PORT_APP}
```

## Endpoints expuestos

Prefijo base: `/api/v1`

- `GET /api/v1/productos` - Lista productos
- `GET /api/v1/productos/:id` - Obtiene un producto por UUID
- `POST /api/v1/productos` - Crea un producto
- `PUT /api/v1/productos/:id` - Actualiza un producto por UUID
- `DELETE /api/v1/productos/:id` - Elimina un producto por UUID

Ejemplos:

```bash
curl http://localhost:8082/api/v1/productos
curl http://localhost:8082/api/v1/productos/<uuid>
```

## Estructura del proyecto

```text
.
├── .air.toml
├── .env
├── .env.example
├── .gitignore
├── README.md
├── config/
│   └── db.go
├── controllers/
│   └── producto_controller.go
├── models/
│   └── producto.go
├── docker-compose.yml
├── main.go
├── go.mod
├── go.sum
└── tmp/
```

## Notas tecnicas

- Se usa `uuid` como clave primaria en `Producto`.
- El endpoint `GET /productos/:id` valida UUID y maneja `record not found`.
- La conexion a DB se inicializa en `main.go` con `config.InitDB()`.