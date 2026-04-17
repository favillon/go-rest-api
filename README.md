# Backend Productos

API REST en Go usando Gin, GORM y PostgreSQL.

## Configuracion

La aplicacion carga variables de entorno desde el archivo `.env`.

Variables usadas:

- `PORT_APP`: puerto donde corre la API. Por defecto: `8082`
- `POSTGRES_HOST`: host de PostgreSQL. Por defecto: `localhost`
- `POSTGRES_PORT`: puerto de PostgreSQL. Por defecto: `5432`
- `POSTGRES_DB`: nombre de la base de datos
- `POSTGRES_USER`: usuario de PostgreSQL
- `POSTGRES_PASSWORD`: contraseña de PostgreSQL

Ejemplo de archivo `.env`:

```dotenv
PORT_APP=8082
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=restapi_go_db
POSTGRES_USER=user_go
POSTGRES_PASSWORD=P4ssW0rS3cr3t
```

## Docker Compose

El archivo `docker-compose.yml` levanta un servicio de PostgreSQL 18 con volumen persistente.

Para iniciar la base de datos:

```bash
docker compose up -d postgres
```

## Ejecucion de la aplicacion

1. Levantar PostgreSQL con Docker.
2. Ejecutar la API:

```bash
go run .
```

## Endpoint expuesto

La API expone actualmente este endpoint:

- `GET /api/v1/productos` - Lista todos los productos

La aplicacion escucha en el puerto definido por `PORT_APP`.

Ejemplo:

```bash
http://localhost:8082/api/v1/productos
```

## Estructura del proyecto

```text
.
├── config/             # Conexion a la base de datos
│   └── db.go
├── controllers/        # Manejadores de rutas
│   └── producto_controller.go
├── models/             # Entidades GORM
│   └── producto.go
├── docker-compose.yml  # Servicio PostgreSQL 18
├── .env                # Variables de entorno locales
├── .env.example        # Plantilla de variables de entorno
├── .gitignore
├── main.go             # Punto de entrada de la API
├── go.mod
└── go.sum
```

## Notas

- El proyecto usa GORM para conectar con PostgreSQL.
- Gin maneja las rutas HTTP.
- Si quieres agregar nuevos endpoints, el archivo principal para rutas es `main.go`.