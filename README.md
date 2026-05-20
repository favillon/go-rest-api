# Backend Productos

**Repositorio:** https://github.com/favillon/go-rest-api/tree/feature/gRPC

API gRPC en Go usando arquitectura hexagonal (puertos y adaptadores), Protocol Buffers y MongoDB.

---

## Configuracion

La aplicacion carga variables de entorno desde el archivo `.env`.

Variables usadas:

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `PORT_GRPC` | Puerto del servidor gRPC | `50051` |
| `MONGO_HOST` | Host de MongoDB | `localhost` |
| `MONGO_PORT` | Puerto de MongoDB | `27017` |
| `MONGO_HOST_PORT` | Puerto expuesto en host | `27017` |
| `MONGO_DB` | Nombre de la base de datos | `productos_db` |
| `MONGO_USER` | Usuario de MongoDB | `admin` |
| `MONGO_PASSWORD` | Password de MongoDB | `change_me` |

Ejemplo (`.env.example`):

```dotenv
PORT_GRPC=50051
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_HOST_PORT=27017
MONGO_DB=productos_db
MONGO_USER=admin
MONGO_PASSWORD=change_me
```

---

## Base de datos (Docker)

El archivo `docker-compose.yml` levanta MongoDB 7 y mapea el puerto:

- `${MONGO_HOST_PORT:-27017}:27017`

Comandos utiles:

```bash
docker compose up -d mongodb
docker compose ps
docker compose down
```

---

## Ejecucion de la API

### Modo normal

```bash
go run .
```

### Modo desarrollo con recarga (Air)

```bash
air
```

El servidor gRPC queda disponible en:

```bash
localhost:${PORT_GRPC}
```

---

## gRPC

La API expone 3 servicios gRPC definidos en `proto/`.

### Servicios y metodos

| Servicio | Metodo | Descripcion |
|----------|--------|-------------|
| **ProductoService** | `CreateProducto` | Crear producto |
| | `GetProducto` | Obtener producto por ID |
| | `UpdateProducto` | Actualizar producto |
| | `DeleteProducto` | Eliminar producto (soft delete) |
| | `ListProductos` | Listar productos paginados |
| **CategoriaService** | `CreateCategoria` | Crear categoria |
| | `GetCategoria` | Obtener categoria por ID |
| | `UpdateCategoria` | Actualizar categoria |
| | `DeleteCategoria` | Eliminar categoria (soft delete) |
| | `ListCategorias` | Listar todas las categorias |
| **InventarioService** | `CreateInventario` | Crear inventario |
| | `GetInventario` | Obtener inventario por ID |
| | `GetInventarioByProductoId` | Obtener inventario por producto ID |
| | `UpdateInventario` | Actualizar inventario |
| | `DeleteInventario` | Eliminar inventario (soft delete) |
| | `ListInventarios` | Listar todos los inventarios |

### Mensajes principales

**Producto:**
```protobuf
message Producto {
  string id = 1;
  string nombre = 2;
  string descripcion = 3;
  double precio = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}
```

**Categoria:**
```protobuf
message Categoria {
  string id = 1;
  string nombre = 2;
  string descripcion = 3;
  repeated string producto_ids = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}
```

**Inventario:**
```protobuf
message Inventario {
  string id = 1;
  string producto_id = 2;
  int32 cantidad = 3;
  string almacen = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}
```

---

## Probar con grpcurl

### Instalacion

```bash
brew install grpcurl
```

### Opcion A: Especificar protos individuales (rapido para pruebas)

Listar servicios:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  localhost:50051 list
```

Describir un servicio:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  localhost:50051 describe productos.ProductoService
```

Crear un producto:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "nombre": "Monitor 4K",
    "descripcion": "Samsung 27 pulgadas",
    "precio": 399.99
  }' localhost:50051 productos.ProductoService/CreateProducto
```

Obtener un producto:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }' localhost:50051 productos.ProductoService/GetProducto
```

Listar productos:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "page": 1,
    "limit": 10
  }' localhost:50051 productos.ProductoService/ListProductos
```

Crear una categoria:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "nombre": "Electronics",
    "descripcion": "Electronic devices"
  }' localhost:50051 productos.CategoriaService/CreateCategoria
```

Crear un inventario:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "producto_id": "550e8400-e29b-41d4-a716-446655440000",
    "cantidad": 50,
    "almacen": "Central"
  }' localhost:50051 productos.InventarioService/CreateInventario
```

Obtener inventario por producto:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "producto_id": "550e8400-e29b-41d4-a716-446655440000"
  }' localhost:50051 productos.InventarioService/GetInventarioByProductoId
```

Actualizar un producto:

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nombre": "Monitor 4K Samsung",
    "descripcion": "Samsung 27 pulgadas actualizado",
    "precio": 349.99
  }' localhost:50051 productos.ProductoService/UpdateProducto
```

Eliminar un producto (soft delete):

```bash
grpcurl -plaintext \
  -proto proto/producto.proto \
  -proto proto/categoria.proto \
  -proto proto/inventario.proto \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }' localhost:50051 productos.ProductoService/DeleteProducto
```

### Opcion B: Generar un archivo descriptor (protoset) — recomendado para uso frecuente

Si no quieres repetir los 3 flags `-proto` en cada llamada, genera un `protoset`:

**1. Generar el archivo descriptor:**

```bash
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --proto_path=. \
  --descriptor_set_out=productos.protoset \
  --include_imports \
  proto/producto.proto \
  proto/categoria.proto \
  proto/inventario.proto
```

**2. Usar grpcurl con el protoset:**

```bash
# Listar servicios
grpcurl -plaintext -protoset productos.protoset localhost:50051 list

# Crear producto
grpcurl -plaintext -protoset productos.protoset \
  -d '{"nombre":"Monitor 4K","precio":399.99}' \
  localhost:50051 productos.ProductoService/CreateProducto

# Obtener producto
grpcurl -plaintext -protoset productos.protoset \
  -d '{"id":"550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 productos.ProductoService/GetProducto
```

El archivo `productos.protoset` es un descriptor binario que contiene toda la informacion de los 3 protos. Solo necesitas regenerarlo cuando modifiques los `.proto`.

---

## Arquitectura Hexagonal

Este proyecto utiliza **arquitectura hexagonal** (puertos y adaptadores).

### Estructura de capas

```
┌─────────────────────────────────────────────────────────┐
│                    Capa de Dominio                       │
│  internal/domain/model/   → Entidades (3 modelos)     │
│  internal/domain/port/    → Interfaces (3 puertos)     │
│  internal/domain/errors.go → ErrNotFound canonico      │
└─────────────────────────────────────────────────────────┘
                            ↕
┌─────────────────────────────────────────────────────────┐
│                 Capa de Aplicacion                       │
│  internal/application/service/ → 3 casos de uso        │
└─────────────────────────────────────────────────────────┘
                            ↕
┌─────────────────────────────────────────────────────────┐
│              Capa de Infraestructura                     │
│  internal/infrastructure/grpc/        → Servidor gRPC    │
│  internal/infrastructure/grpc/handler/ → Handlers gRPC  │
│  internal/infrastructure/persistence/mongodb/ → 3 repos│
└─────────────────────────────────────────────────────────┘
```

### Flujo de una solicitud

1. **gRPC** → Handler recibe la request, traduce protobuf → domain
2. **Service** → Servicio ejecuta la logica de negocio
3. **Repository** → Repositorio MongoDB accede a la base de datos
4. **Response** → El handler traduce domain → protobuf y devuelve la respuesta

### Flujo de inyeccion de dependencias (main.go)

```go
// MongoDB repositories
productoRepo   := mongodb.NewProductoRepository(config.MongoDatabase)
categoriaRepo  := mongodb.NewCategoriaRepository(config.MongoDatabase)
inventarioRepo := mongodb.NewInventarioRepository(config.MongoDatabase)

// Application services
productoSvc   := service.NewProductoService(productoRepo)
categoriaSvc  := service.NewCategoriaService(categoriaRepo)
inventarioSvc := service.NewInventarioService(inventarioRepo)

// gRPC server
server := grpcserver.NewServer(productoSvc, categoriaSvc, inventarioSvc)
```

---

## Resumen de testing

El proyecto cuenta con **37 tests**:

| Capa | Archivo | Tests |
|------|---------|-------|
| **Service** | `internal/application/service/producto_service_test.go` | 13 |
| **Handler gRPC** | `internal/infrastructure/grpc/handler/producto_handler_test.go` | 10 |
| **Handler gRPC** | `internal/infrastructure/grpc/handler/categoria_handler_test.go` | 7 |
| **Handler gRPC** | `internal/infrastructure/grpc/handler/inventario_handler_test.go` | 7 |

Casos cubiertos:

- **Create**: exito
- **Get**: exito, not found
- **Update**: exito, not found
- **Delete**: exito, not found
- **List**: exito, error DB
- **GetByProductoId** (Inventario): exito

---

## Comandos de testing y coverage

Ejecutar todos los tests:

```bash
go test ./...
```

Ejecutar tests con detalle:

```bash
go test ./... -v
```

Ejecutar tests por capa:

```bash
go test ./internal/application/service -v           # Tests del service
go test ./internal/infrastructure/grpc/handler -v    # Tests de handlers gRPC
```

Generar reporte de coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

### Dentro del contenedor de desarrollo

```bash
docker compose exec app go test ./...
docker compose exec app go test ./... -coverprofile=coverage.out
docker compose exec app go tool cover -html=coverage.out
```

---

## Docker: Desarrollo y Produccion

Los puertos de la aplicacion y de MongoDB se configuran desde el archivo `.env` (o `.env.docker`).

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `PORT_GRPC` | Puerto expuesto por el servidor gRPC Go | `50051` |
| `MONGO_HOST_PORT` | Puerto expuesto por MongoDB en el host | `27017` |

### Opcion 1: Desarrollo con Docker Compose (Hot Reload)

Levantar la aplicacion y MongoDB con recarga automatica usando Air:

```bash
docker compose --env-file .env.docker up --build
```

Esto ejecuta:
- **MongoDB**: en el puerto del host definido por `MONGO_HOST_PORT`
- **App con Air** (`Dockerfile.dev`): en puerto `${PORT_GRPC}`, recarga automatica en cambios de codigo

Ver logs en vivo:

```bash
docker compose logs -f app
```

Para detener:

```bash
docker compose down
```

### Opcion 2: Produccion - Docker Compose Prod

Levantar la aplicacion con la imagen de produccion optimizada (multietapa, sin Air, sin codigo fuente):

```bash
docker compose -f docker-compose.prod.yml --env-file .env.docker up --build -d
```

Esto ejecuta:
- **MongoDB**: en el puerto del host definido por `MONGO_HOST_PORT`
- **App** (`Dockerfile`): binario compilado, ~5-10 MB, puerto `${PORT_GRPC}`

Ver logs:

```bash
docker compose -f docker-compose.prod.yml logs -f app
```

Para detener:

```bash
docker compose -f docker-compose.prod.yml down
```

### Opcion 3: Build manual de la imagen de produccion

Si prefieres construir y correr la imagen manualmente:

```bash
# Build con puerto personalizado
docker build \
  --build-arg PORT_GRPC=50051 \
  -t backend-productos:latest .

# Run (requiere MongoDB externo)
docker run \
  -e MONGO_HOST=mongodb \
  -e MONGO_USER=admin \
  -e MONGO_PASSWORD=secret \
  -e MONGO_DB=productos_db \
  -e PORT_GRPC=50051 \
  -p 50051:50051 \
  backend-productos:latest
```

### Estructura de Dockerfiles

- **Dockerfile**: Multietapa para produccion
  - Acepta `ARG PORT_GRPC=50051` en build time
  - Etapa 1: golang:1.26.2-alpine3.22, compila el binario
  - Etapa 2: alpine:3.18, solo binario, ~5-10 MB
  - Healthcheck usa `nc -z` al puerto gRPC

- **Dockerfile.dev**: Para desarrollo con Air
  - Acepta `ARG PORT_GRPC=50051` en build time
  - golang:1.26.2-alpine3.22 + Air + golangci-lint preinstalados
  - Recarga automatica en cambios de codigo
  - Volumen montado del codigo fuente

### Archivos Compose

| Archivo | Proposito | Dockerfile usado |
|---------|-----------|------------------|
| `docker-compose.yml` | Desarrollo local con hot-reload | `Dockerfile.dev` |
| `docker-compose.prod.yml` | Produccion, imagen optimizada | `Dockerfile` |

### .env requerido para Docker

Crea un archivo `.env` (o `.env.docker`) con las variables necesarias:

```env
PORT_GRPC=50051
MONGO_HOST_PORT=27017
MONGO_USER=admin
MONGO_PASSWORD=P4ssW0rS3cr3t
MONGO_DB=productos_db
```

Para desarrollo en Docker, anade tambien:

```env
MONGO_HOST=mongodb
MONGO_PORT=27017
```

Para desarrollo local (fuera de Docker):

```env
MONGO_HOST=localhost
MONGO_PORT=27017
```

Ver `.env.example` para mas detalles.

---

## Verificacion de arquitectura (Arch-Go)

El proyecto usa [Arch-Go](https://github.com/arch-go/arch-go) para validar que las dependencias entre capas respetan la arquitectura hexagonal.

### Instalacion

```bash
go install -v github.com/arch-go/arch-go/v2@latest
```

### Ejecucion

```bash
arch-go
```

Salida esperada:

```
Compliance:      100% (PASS)
Coverage:        100% (PASS)
```

### Opciones utiles

```bash
arch-go --verbose          # Muestra detalle por paquete
arch-go --json             # Genera reporte en .arch-go/report.json
arch-go --html             # Genera reporte en .arch-go/report.html
arch-go describe           # Describe las reglas configuradas
```

### Reglas configuradas (`arch-go.yml`)

| Capa | Solo puede depender de |
|------|----------------------|
| `domain` | `domain`, librerias externas (`uuid`) |
| `application` | `domain`, `application`, librerias externas |
| `infrastructure` | `domain`, `application`, `infrastructure`, `proto`, `mongo-driver`, `grpc` |
| `proto` | `proto`, `protobuf` |
| `config` | `config`, `mongo-driver` |

### Dentro del contenedor de desarrollo

```bash
docker compose exec app arch-go
```

---

## Comandos para linters

Formatear codigo con gofmt:

```bash
gofmt -w .
```

Analisis estatico con go vet:

```bash
go vet ./...
```

Instalar golangci-lint (si no esta instalado):

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Ejecutar golangci-lint:

```bash
golangci-lint run ./...
```

### Dentro del contenedor de desarrollo

```bash
docker compose exec app gofmt -w .
docker compose exec app go vet ./...
docker compose exec app golangci-lint run ./...
```

---

## Regenerar codigo Protobuf

Si modificas los archivos `.proto`, regenera los stubs:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/*.proto
```

Esto actualiza:

- `proto/*_grpc.pb.go` — interfaces y clientes gRPC
- `proto/*.pb.go` — structs y serializacion

**Recomendacion:** Despues de regenerar, ejecuta:

```bash
go build ./...
go test ./...
```

---

## Estructura del proyecto

```text
.
├── .air.toml
├── .dockerignore
├── .env
├── .env.docker
├── .env.example
├── .gitignore
├── .golangci.yml
├── arch-go.yml
├── plan.md                     # Legacy (GraphQL)
├── plan_gRPC.md               # Plan actual
├── README.md
├── config/
│   └── mongo.go               # Conexion MongoDB
├── proto/
│   ├── producto.proto         # Schema Producto
│   ├── categoria.proto        # Schema Categoria
│   ├── inventario.proto       # Schema Inventario
│   ├── producto.pb.go         # Stub generado
│   ├── producto_grpc.pb.go    # Stub generado
│   ├── categoria.pb.go        # Stub generado
│   ├── categoria_grpc.pb.go   # Stub generado
│   ├── inventario.pb.go       # Stub generado
│   └── inventario_grpc.pb.go  # Stub generado
├── internal/
│   ├── application/
│   │   └── service/
│   │       ├── producto_service.go        # Caso de uso Producto
│   │       ├── producto_service_test.go   # 13 tests
│   │       ├── categoria_service.go       # Caso de uso Categoria
│   │       └── inventario_service.go      # Caso de uso Inventario
│   ├── domain/
│   │   ├── errors.go                      # ErrNotFound
│   │   ├── model/
│   │   │   ├── producto.go                # Entidad Producto
│   │   │   ├── categoria.go               # Entidad Categoria
│   │   │   └── inventario.go              # Entidad Inventario
│   │   └── port/
│   │       ├── producto_repository.go     # Puerto Producto
│   │       ├── categoria_repository.go    # Puerto Categoria
│   │       └── inventario_repository.go   # Puerto Inventario
│   └── infrastructure/
│       ├── grpc/
│       │   ├── server.go                  # Servidor gRPC
│       │   ├── interceptor/
│       │   │   ├── recovery.go           # Recovery interceptor
│       │   │   └── logging.go             # Logging interceptor
│       │   └── handler/
│       │       ├── producto_handler.go     # Handler Producto
│       │       ├── producto_handler_test.go # 10 tests
│       │       ├── categoria_handler.go    # Handler Categoria
│       │       ├── categoria_handler_test.go # 7 tests
│       │       ├── inventario_handler.go   # Handler Inventario
│       │       └── inventario_handler_test.go # 7 tests
│       └── persistence/
│           ├── mongodb/
│           │   ├── producto_repository.go   # Repo Producto
│           │   ├── categoria_repository.go  # Repo Categoria
│           │   └── inventario_repository.go # Repo Inventario
│           └── postgres/                   # Legacy (referencia)
│               └── producto_repository.go
├── docker-compose.yml
├── docker-compose.prod.yml
├── Dockerfile
├── Dockerfile.dev
├── main.go                         # Entrypoint gRPC-only
├── go.mod
└── go.sum
```

---

## Notas tecnicas

- Se usa `uuid` como clave primaria en todas las entidades.
- MongoDB almacena `uuid.UUID` como `Binary` (subtype 4) en el campo `_id`.
- Soft delete via campo `deleted_at` (timestamp) en todas las colecciones.
- **Arquitectura hexagonal**: las dependencias fluyen hacia adentro (infraestructura → dominio).
- **Inyeccion de dependencias**: `main.go` instancia `repo → service → gRPC server` y los inyecta.
- **Paginacion gRPC**: `ListProductos` usa `page` y `limit` (default: page=1, limit=10).
- **Interceptores gRPC**: recovery (evita panics) y logging (registra duracion y errores).
- **Colecciones MongoDB**: `productos`, `categorias`, `inventarios`.
- **gRPC reflection**: No esta habilitado. Se usa `-proto` o `-protoset` con `grpcurl`.
