# Backend Productos

API REST + GraphQL en Go usando arquitectura hexagonal (puertos y adaptadores), Gin, gqlgen y MongoDB.

<div style="display:flex; flex-wrap:wrap; gap:12px; align-items:center;">

<svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 50 50">
<path d="M 22.462891 15.003906 C 21.863562 15.016969 21.255125 15.082359 20.640625 15.193359 C 17.634625 15.742359 15.176312 17.284797 13.320312 19.716797 C 11.595313 21.964797 10.759609 24.526953 11.099609 27.376953 C 11.387609 29.781953 12.485922 31.715688 14.419922 33.179688 C 16.510922 34.748688 18.889172 35.244297 21.451172 34.904297 C 24.562172 34.486297 27.019344 32.943938 28.902344 30.460938 C 28.950344 30.397938 28.98525 30.328625 29.03125 30.265625 C 29.51625 31.165625 30.168375 31.974688 30.984375 32.679688 C 32.831375 34.262687 35.020875 34.948953 37.421875 35.001953 C 38.107875 34.922953 38.81825 34.896625 39.53125 34.765625 C 41.98425 34.264625 44.121281 33.156672 45.863281 31.388672 C 48.316281 28.905672 49.346437 26.002406 48.898438 22.441406 C 48.555437 19.908406 47.261734 17.983594 45.177734 16.558594 C 42.882734 15.001594 40.349203 14.739844 37.658203 15.214844 C 34.536203 15.765844 32.224688 17.075078 30.304688 19.580078 C 29.658687 18.328078 28.768359 17.267609 27.568359 16.474609 C 25.980609 15.396109 24.260875 14.964719 22.462891 15.003906 z M 4 20 A 1.0001 1.0001 0 1 0 4 22 L 8 22 A 1.0001 1.0001 0 1 0 8 20 L 4 20 z M 22.134766 20.011719 C 22.734971 20.038039 23.345266 20.193922 23.962891 20.498047 C 24.511891 20.759047 24.799406 21.048188 25.191406 21.492188 C 25.531406 21.884187 25.557219 21.858906 25.949219 21.753906 C 27.262219 21.414906 28.261844 21.140125 29.464844 20.828125 C 29.007844 21.595125 28.760625 22.155422 28.515625 22.982422 L 21.798828 22.982422 C 21.406828 22.982422 21.223531 23.243391 21.144531 23.400391 C 20.778531 24.080391 20.1525 25.44 19.8125 26.25 C 19.6305 26.695 19.759594 27.035156 20.308594 27.035156 L 25.353516 27.035156 C 25.092516 27.401156 24.770156 27.948328 24.535156 28.236328 C 23.359156 29.569328 21.868797 30.198891 20.091797 29.962891 C 18.026797 29.674891 16.5885 27.948422 16.5625 25.857422 C 16.5365 23.739422 17.452469 22.040625 19.230469 20.890625 C 20.161719 20.28625 21.134424 19.967852 22.134766 20.0... (truncated)
</svg>

<svg xmlns="http://www.w3.org/2000/svg" width="100" height="140.6" viewBox="0 0 100 140.6"><path style="fill:#000" d="M49.002 140.466c-12.056 -0.613 -19.734 -2.105 -26.165 -5.086 -0.767 -0.355 -0.675 -0.159 -1.628 -3.49 -1.613 -5.638 -3.241 -12.462 -4.596 -19.266 -0.79 -3.97 -0.899 -4.459 -2.762 -12.4 -2.988 -12.743 -3.749 -16.578 -4.708 -23.75 -0.623 -4.662 0.716 -7.703 1.613 -3.663 3.528 15.898 6.582 27.504 9.81 37.293 1.113 3.375 1.026 3.235 1.96 3.145 0.642 -0.062 0.643 -0.057 0.163 1.06 -0.655 1.525 -0.335 2.337 1.706 4.339 1.529 1.5 1.544 1.588 0.238 1.398 -1.019 -0.149 -0.969 -0.254 -0.599 1.252 0.173 0.702 0.473 2.019 0.668 2.927 1.548 7.211 3.569 10.737 6.512 11.366 8.599 1.838 21.088 2.636 32.086 2.051 12.776 -0.68 15.833 -1.165 19.8 -3.144 2.683 -1.339 4.994 -3.847 3.578 -3.884 -1.043 -0.028 -1.141 -0.272 -0.256 -0.639 1.927 -0.799 2.741 -1.649 2.583 -2.7 -0.047 -0.314 0.013 -0.641 0.247 -1.358 0.722 -2.203 0.842 -3.995 0.658 -9.744 -0.086 -2.661 -0.079 -2.852 0.123 -3.85 0.533 -2.63 0.773 -4.618 1.165 -9.7 0.278 -3.587 0.367 -4.485 0.654 -6.55 0.288 -2.069 0.38 -3.242 0.544 -6.95 0.321 -7.252 0.679 -10.14 1.339 -10.8 1.104 -1.104 1.098 -1.692 -0.035 -3.298 -1.291 -1.832 -1.287 -2.791 0.018 -4.164 0.636 -0.669 0.663 -1.417 0.229 -6.238 -0.112 -1.238 -0.231 -3.094 -0.265 -4.125 -0.043 -1.271 -0.099 -1.875 -0.174 -1.875 -0.24 0 -1.791 1.325 -5.208 4.451 -3.885 3.553 -4.9 4.391 -7.079 5.846 -4.493 3 -8.808 4.469 -16.271 5.541 -7.533 1.082 -16.401 0.826 -22.207 -0.64 -3.912 -0.987 -6.177 -2.054 -10.59 -4.985 -4.133 -2.745 -3.908 -7.189 0.463 -9.149 1.969 -0.882 0.445 -1.844 -3.506 -2.213 -2.972 -0.277 -1.461 -0.855 2.34 -0.895 4.618 -0.049 7.9 -1.689 12.143 -6.069 1.776 -1.834 2.191 -2.194 4.558 -3.964 3.763 -2.816 5.091 -4.038 5.692 -5.239 0.454 -0.908 0.428 -1.805 -0.056 -1.989 -0.159 -0.06 -0.237 -0.17 -0.237 -0.332 0 -0.302 0.92 -0.259 -15 -0.714 -16.796 -0.481 -29.51 -1.24 -35.999 -2.149 -1.268 -0.178 -0.664 0.946 0.883 1.644 0.489 0.22 0.535 0.418 0.115... (truncated)

<svg stroke-miterlimit="10" style="fill-rule:nonzero;clip-rule:evenodd;stroke-linecap:round;stroke-linejoin:round" viewBox="0 0 100 34.615" xml:space="preserve" xmlns="http://www.w3.org/2000/svg" xmlns:vectornator="http://vectornator.io" width="100" height="34.615"><g vectornator:layerName="Layer 1"><g vectornator:layerName="LOGO1"><g vectornator:layerName="Gopher"><g vectornator:layerName="query"><g vectornator:layerName="Group 3"><path d="M39.382 16.846c.189-.627.859-.98 1.495-.787s.999.856.81 1.483c-.19.627-.859.98-1.495.787s-.999-.856-.81-1.483" fill="#8fd2f9" fill-rule="evenodd" vectornator:layerName="circle"/></g><g vectornator:layerName="g 1"><path d="m41.64 18.682-.31-.569a1.2 1.2 0 0 0 .369-.56 1.204 1.204 0 0 0-.815-1.493 1.203 1.203 0 0 0-1.505.792 1.203 1.203 0 0 0 .815 1.493 1.2 1.2 0 0 0 .68.008l.31.569a.263.263 0 0 0 .352.106.253.253 0 0 0 .105-.346m-1.422-.418c-.595-.18-.933-.8-.757-1.386a1.117 1.117 0 0 1 1.397-.736c.595.18.933.8.757 1.386a1.117 1.117 0 0 1-1.397.736m1.276.689a.175.175 0 0 1-.234-.071l-.303-.555a1.2 1.2 0 0 0 .304-.16l.303.554a.17.17 0 0 1-.07.231" fill="#38b6ff" vectornator:layerName="path"/><path d="m41.587 18.749-.329-.621a1 1 0 0 1-.329.179l.329.621a.18.18 0 0 0 .11.091c.045.014.097.01.142-.015a.19.19 0 0 0 .077-.255" fill="#38b6ff" vectornator:layerName="path"/><path d="M40.808 16.285a.957.957 0 0 0-1.197.63c-.151.502.139 1.033.648 1.187s1.045-.128 1.196-.63a.955.955 0 0 0-.648-1.187m-.526 1.741c-.467-.141-.733-.628-.594-1.088s.63-.718 1.097-.577a.876.876 0 0 1 .594 1.088.876.876 0 0 1-1.097.577" fill="#fff" vectornator:layerName="path"/><path d="M40.64 16.565a.04.04 0 0 0-.05.026.04.04 0 0 0 .027.049c.318.096.5.428.405.742-.007.021.006.043.027.049s.043-.005.05-.026a.68.68 0 0 0-.459-.841" fill="#fff" vectornator:layerName="path"/></g></g><g vectornator:layerName="g"><g vectornator:layerName="g" fill-rule="evenodd"><path d="M50.807 14.335c.023-.14.311-.355.43-.404 1.452-.591.943 1.074.098 1.008" fill="#46bbff" vectornator:layer... (truncated)

<svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 50 50">
<path d="M 35.199219 2.101563 C 33.699219 2.101563 32.398438 2.398438 31.199219 2.699219 C 33.300781 3.597656 34.601563 4.699219 35.300781 5.199219 C 36.902344 6.597656 37.800781 8 39.402344 10.300781 C 39.699219 10.800781 40.199219 11.5 40.5 12.597656 C 40.800781 13.597656 40.800781 14.398438 40.800781 15.199219 C 40.800781 16.300781 40.699219 17.199219 40.597656 18.097656 C 40.5 18.800781 40.5 19.101563 40.402344 19.300781 C 40.402344 19.402344 40.402344 19.5 40.300781 19.699219 C 40.300781 20.199219 40.300781 20.402344 40.402344 20.800781 C 40.402344 21.199219 40.5 21.601563 40.5 22.300781 C 40.601563 23.601563 40.601563 24.5 40.402344 25.597656 L 40.402344 26 C 40.199219 26.898438 40 27.800781 39.5 28.597656 C 39.601563 28.800781 39.699219 28.898438 39.800781 29.097656 C 40.300781 28.398438 40.699219 27.699219 41.097656 26.902344 C 42.300781 24.699219 43 22.800781 43.5 21.402344 C 44.398438 18.800781 44.898438 16.898438 45.199219 15.5 C 45.898438 12.5 46 11.101563 45.699219 9.5 C 45.699219 9 45.5 8.097656 45 7.199219 C 43.898438 5.199219 42.199219 4.300781 41 3.699219 C 40.199219 3.300781 38.097656 2.199219 35.199219 2.101563 Z M 13.535156 2.542969 C 12.382813 2.519531 10.976563 2.648438 9.398438 3.398438 C 8.898438 3.601563 7.398438 4.300781 6.199219 5.898438 C 5.398438 6.898438 4.800781 8.398438 4.5 10.097656 C 4.199219 11.597656 4.097656 13.402344 4.699219 16.800781 C 5.097656 19.199219 5.5 20.800781 6.300781 24.097656 C 6.402344 24.5 7 26.300781 8.300781 30.300781 L 8.398438 30.5 C 8.601563 31.199219 9.199219 32.699219 10.5 34.199219 C 11.398438 35.199219 12.199219 35.800781 12.902344 35.800781 L 13.097656 35.800781 C 14.398438 35.800781 15.300781 34.800781 16.097656 34 C 16.097656 33.898438 18 31.601563 18.699219 30.800781 C 18.597656 30.699219 18.402344 30.699219 18.300781 30.597656 C 17.101563 29.898438 16.199219 28.800781 15.5 27.597656 C 14.300781 25.398438 14.398438 23.199219 14.597656 22.097656 L 14.699219 20.402344 C 14.300781 17.699219 14.402344 15.... (truncated)
</svg>

<svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 50 50">
<path d="M 20 9 L 20 14 L 10 14 L 10 19 L 5 19 L 5 24 L 1.125 24 C 0.640625 24 0.242188 24.335938 0.15625 24.8125 C 0.148438 24.847656 0 25.683594 0 26.75 C 0 27.449219 0.0664063 28.210938 0.1875 28.96875 C 1.332031 28.695313 3.429688 28.285156 3.0625 26.9375 C 5.035156 29.222656 9.769531 28.53125 10.96875 27.40625 C 12.308594 29.347656 20.113281 28.605469 20.65625 27.09375 C 22.335938 29.0625 27.542969 29.0625 29.21875 27.09375 C 29.761719 28.605469 37.535156 29.347656 38.875 27.40625 C 39.300781 27.804688 40.1875 28.136719 41.21875 28.3125 C 41.566406 27.652344 41.886719 26.988281 42.1875 26.28125 C 48.539063 26.203125 49.910156 21.636719 49.96875 21.4375 C 50.078125 21.054688 49.929688 20.660156 49.625 20.40625 C 49.519531 20.316406 47.175781 18.414063 43.375 19.0625 C 42.308594 15.589844 39.5625 14.007813 39.4375 13.9375 C 39.078125 13.734375 38.632813 13.765625 38.3125 14.03125 C 38.210938 14.113281 35.847656 16.117188 36.21875 20.21875 C 36.3125 21.25 36.582031 22.160156 37 22.96875 C 36.179688 23.425781 34.769531 24 32.5 24 L 32 24 L 32 19 L 27 19 L 27 9 Z M 41.21875 28.3125 C 41.097656 28.546875 40.941406 28.773438 40.8125 29 L 49.84375 29 C 48.757813 28.726563 46.425781 28.359375 46.8125 26.9375 C 45.535156 28.414063 43.109375 28.632813 41.21875 28.3125 Z M 40.8125 29 L 0.1875 29 C 0.429688 30.46875 0.929688 32.007813 1.6875 33.5 C 7.117188 34.777344 12.816406 32.832031 12.875 32.8125 C 13.398438 32.628906 13.945313 32.917969 14.125 33.4375 C 14.308594 33.957031 14.050781 34.539063 13.53125 34.71875 C 13.339844 34.785156 9.90625 35.9375 5.6875 35.9375 C 4.851563 35.9375 3.972656 35.890625 3.09375 35.78125 C 5.71875 39.261719 10.167969 42 17 42 C 27.804688 42 36.113281 37.410156 40.8125 29 Z M 0.1875 29 C 0.183594 28.984375 0.191406 28.984375 0.1875 28.96875 C 0.121094 28.984375 0.0585938 28.984375 0 29 Z M 22 11 L 25 11 L 25 14 L 22 14 Z M 12 16 L 15 16 L 15 19 L 12 19 Z M 17 16 L 20 16 L 20 19 L 17 19 Z M 22 16 L 25 16 L 25 19 L 22 19 Z M 7 21 L 10 21 L 10... (truncated)
</svg>

</div>

## Configuracion

La aplicacion carga variables de entorno desde el archivo `.env`.

Variables usadas:

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `PORT_APP` | Puerto de la API | `8082` |
| `MONGO_HOST` | Host de MongoDB | `localhost` |
| `MONGO_PORT` | Puerto de MongoDB | `27017` |
| `MONGO_HOST_PORT` | Puerto expuesto en host | `27017` |
| `MONGO_DB` | Nombre de la base de datos | `productos_db` |
| `MONGO_USER` | Usuario de MongoDB | `admin` |
| `MONGO_PASSWORD` | Password de MongoDB | `change_me` |
| `RATE_LIMIT_READ_PER_MINUTE` | Limite lectura por IP/min | `30` |
| `RATE_LIMIT_WRITE_PER_MINUTE` | Limite escritura por IP/min | `10` |

Ejemplo (`.env.example`):

```dotenv
PORT_APP=8082
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_HOST_PORT=27017
MONGO_DB=productos_db
MONGO_USER=admin
MONGO_PASSWORD=change_me
RATE_LIMIT_READ_PER_MINUTE=30
RATE_LIMIT_WRITE_PER_MINUTE=10
```

## Base de datos (Docker)

El archivo `docker-compose.yml` levanta MongoDB 7 y mapea el puerto:

- `${MONGO_HOST_PORT:-27017}:27017`

Comandos utiles:

```bash
docker compose up -d mongodb
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

## Documentacion Swagger / OpenAPI

La API expone Swagger UI en:

```bash
http://localhost:${PORT_APP}/swagger/index.html
```

### Regenerar documentacion

```bash
go run github.com/swaggo/swag/cmd/swag@latest init
```

Esto actualiza `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.

Notas:

- `docs/swagger.json` y `docs/swagger.yaml` estan en `.gitignore`.
- `docs/docs.go` se mantiene en el repositorio.

## Endpoints REST

Prefijo base: `/api/v1`

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| `GET` | `/api/v1/productos` | Lista productos (paginado) |
| `GET` | `/api/v1/productos/:id` | Obtiene un producto por UUID |
| `POST` | `/api/v1/productos` | Crea un producto |
| `PUT` | `/api/v1/productos/:id` | Actualiza un producto |
| `DELETE` | `/api/v1/productos/:id` | Elimina un producto (soft delete) |

Ejemplos:

```bash
# Listar productos
curl http://localhost:8082/api/v1/productos
curl http://localhost:8082/api/v1/productos?page=1&limit=10

# Obtener por ID
curl http://localhost:8082/api/v1/productos/<uuid>

# Crear producto
curl -X POST http://localhost:8082/api/v1/productos \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Laptop Dell","descripcion":"XPS 15","precio":1299.99}'

# Actualizar producto
curl -X PUT http://localhost:8082/api/v1/productos/<uuid> \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Laptop Dell XPS","precio":1199.99}'

# Eliminar producto
curl -X DELETE http://localhost:8082/api/v1/productos/<uuid>
```

## GraphQL

La API expone un endpoint GraphQL en `/api/v1/graphql` con Playground integrado.

### Acceso

- **Playground UI**: `http://localhost:8082/api/v1/graphql` (GET)
- **Endpoint**: `http://localhost:8082/api/v1/graphql` (POST)

### Schema

El schema define 3 entidades con relaciones:

| Tipo | Campos |
|------|--------|
| `Producto` | `id`, `nombre`, `descripcion`, `precio`, `createdAt`, `updatedAt`, `categoria`, `inventario` |
| `Categoria` | `id`, `nombre`, `descripcion`, `productoIds`, `productos` |
| `Inventario` | `id`, `productoId`, `cantidad`, `almacen` |

### Queries

```graphql
# Listar productos con relaciones
query {
  productos(page: 1, limit: 10) {
    id
    nombre
    precio
    categoria {
      nombre
    }
    inventario {
      cantidad
      almacen
    }
  }
}

# Obtener un producto
query {
  producto(id: "550e8400-e29b-41d4-a716-446655440000") {
    id
    nombre
    precio
  }
}

# Listar categorias con sus productos
query {
  categorias {
    id
    nombre
    productos {
      id
      nombre
    }
  }
}

# Obtener inventario por producto
query {
  inventarioByProducto(productoId: "550e8400-e29b-41d4-a716-446655440000") {
    cantidad
    almacen
  }
}
```

### Mutations

```graphql
# Crear producto
mutation {
  crearProducto(input: {
    nombre: "Monitor 4K"
    descripcion: "Samsung 27 pulgadas"
    precio: 399.99
  }) {
    id
    nombre
    precio
  }
}

# Actualizar producto
mutation {
  actualizarProducto(
    id: "550e8400-e29b-41d4-a716-446655440000"
    input: { nombre: "Monitor 4K Samsung", precio: 349.99 }
  ) {
    id
    nombre
    precio
  }
}

# Eliminar producto
mutation {
  eliminarProducto(id: "550e8400-e29b-41d4-a716-446655440000")
}

# Crear categoria
mutation {
  crearCategoria(input: {
    nombre: "Electronics"
    descripcion: "Electronic devices"
    productoIds: ["550e8400-e29b-41d4-a716-446655440000"]
  }) {
    id
    nombre
    productoIds
  }
}

# Crear inventario
mutation {
  crearInventario(input: {
    productoId: "550e8400-e29b-41d4-a716-446655440000"
    cantidad: 50
    almacen: "Central"
  }) {
    id
    cantidad
    almacen
  }
}
```

## Rate limiting

Los limites de solicitudes por IP se configuran mediante variables de entorno:

- `RATE_LIMIT_READ_PER_MINUTE`: limite para endpoints de lectura (`GET`). Default: `30`
- `RATE_LIMIT_WRITE_PER_MINUTE`: limite para endpoints de escritura (`POST`, `PUT`, `DELETE`). Default: `10`

Ventana de tiempo: 1 minuto.

Los headers `X-RateLimit-Limit` y `X-RateLimit-Remaining` se incluyen en cada respuesta. Al exceder el limite se devuelve `429` con header `Retry-After`.

## Catalogo de codigos de error

Fuente de verdad en codigo:

- `api/errors/codes.go`
- `api/response/response.go`

Codigos actuales:

| Codigo | HTTP | Descripcion |
|--------|------|-------------|
| `DB_NOT_INITIALIZED` | 500 | Conexion a DB no inicializada |
| `DATABASE_ERROR` | 500 | Error en operacion de DB |
| `INVALID_PARAMETER` | 400 | Parametro invalido |
| `VALIDATION_ERROR` | 400 | Payload invalido |
| `RESOURCE_NOT_FOUND` | 404 | Recurso no existe |
| `RATE_LIMIT_EXCEEDED` | 429 | Limite de solicitudes excedido |

## Arquitectura Hexagonal

Este proyecto utiliza **arquitectura hexagonal** (puertos y adaptadores).

### Estructura de capas

```
┌─────────────────────────────────────────────────────────┐
│                    Capa de Dominio                       │
│  internal/domain/model/   → Entidades (3 modelos)       │
│  internal/domain/port/    → Interfaces (3 puertos)       │
│  internal/domain/errors.go → ErrNotFound canónico       │
└─────────────────────────────────────────────────────────┘
                           ↕
┌─────────────────────────────────────────────────────────┐
│                 Capa de Aplicacion                       │
│  internal/application/service/ → 3 casos de uso          │
└─────────────────────────────────────────────────────────┘
                           ↕
┌─────────────────────────────────────────────────────────┐
│              Capa de Infraestructura                     │
│  internal/infrastructure/http/        → Handlers REST    │
│  internal/infrastructure/persistence/mongodb/ → 3 repos  │
│  graph/                               → Resolvers GraphQL│
└─────────────────────────────────────────────────────────┘
```

### Flujo de una solicitud

1. **HTTP/GraphQL** → Handler/Resolver recibe la request, valida y parsea
2. **Service** → Servicio ejecuta la logica de negocio
3. **Repository** → Repositorio MongoDB accede a la base de datos
4. **Response** → El handler formatea y devuelve la respuesta

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

// GraphQL resolver
resolver := &graph.Resolver{
    ProductoService:   productoSvc,
    CategoriaService:  categoriaSvc,
    InventarioService: inventarioSvc,
}
```

## Resumen de testing

El proyecto cuenta con **35 tests** distribuidos en 2 capas:

| Capa | Archivo | Tests | Cobertura |
|------|---------|-------|-----------|
| **Service** | `internal/application/service/producto_service_test.go` | 13 | 100% |
| **Handler** | `internal/infrastructure/http/producto_handler_test.go` | 22 | ~97% |

Casos cubiertos por endpoint:

- **Crear producto**: exito, validacion fallida, error DB
- **Obtener productos**: exito, paginacion valida, paginacion invalida, error DB
- **Obtener producto por ID**: exito, ID invalido, no encontrado, error DB
- **Actualizar producto**: exito, ID invalido, no encontrado, payload invalido, error DB lookup, error DB update
- **Eliminar producto**: exito, ID invalido, no encontrado, error DB lookup, error DB delete
- **Service**: GetAll, GetByID, Create, Update, Delete (exito, not found, error)

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
go test ./internal/application/service -v      # Tests del service
go test ./internal/infrastructure/http -v       # Tests del handler
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

## Docker: Desarrollo y Produccion

### Opcion 1: Desarrollo con Docker Compose (Hot Reload)

Levantar la aplicacion y MongoDB con recarga automatica usando Air:

```bash
docker compose --env-file .env.docker up --build
```

Esto ejecuta:
- **MongoDB**: en el puerto del host definido por `MONGO_HOST_PORT`
- **App con Air** (Dockerfile.dev): en puerto `8082`, recarga automatica en cambios de codigo

Ver logs en vivo:

```bash
docker compose logs -f app
```

Para detener:

```bash
docker compose down
```

### Opcion 2: Produccion - Imagen Optimizada

Compilar la imagen multietapa de produccion:

```bash
# Build
docker build -t backend-productos:latest .

# Run (requiere MongoDB externo)
docker run \
  -e MONGO_HOST=mongodb \
  -e MONGO_USER=admin \
  -e MONGO_PASSWORD=secret \
  -e MONGO_DB=productos_db \
  -p 8082:8082 \
  backend-productos:latest
```

La imagen final de produccion:
- Tamaño: ~5-10 MB (solo binario compilado)
- Sin codigo fuente ni herramientas de desarrollo
- Segura para despliegue en produccion

### Estructura de Dockerfiles

- **Dockerfile**: Multietapa para produccion
  - Etapa 1: golang:1.26.2-alpine3.22, compila el binario
  - Etapa 2: alpine:3.18, solo binario, ~5-10 MB

- **Dockerfile.dev**: Para desarrollo con Air
  - golang:1.26.2-alpine3.22 + Air + golangci-lint preinstalados
  - Recarga automatica en cambios de codigo
  - Volumen montado del codigo fuente

### .env requerido para Docker

Para desarrollo local, usa `.env`:

```env
MONGO_USER=admin
MONGO_PASSWORD=P4ssW0rS3cr3t
MONGO_DB=productos_db
PORT_APP=8082
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_HOST_PORT=27017
RATE_LIMIT_READ_PER_MINUTE=30
RATE_LIMIT_WRITE_PER_MINUTE=10
```

Para desarrollo en Docker, existe `.env.docker` (cargado automaticamente en contenedor):

```env
MONGO_USER=admin
MONGO_PASSWORD=P4ssW0rS3cr3t
MONGO_DB=productos_db
PORT_APP=8082
MONGO_HOST=mongodb
MONGO_PORT=27017
MONGO_HOST_PORT=27017
RATE_LIMIT_READ_PER_MINUTE=30
RATE_LIMIT_WRITE_PER_MINUTE=10
```

Ver `.env.example` para mas detalles.

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
| `infrastructure` | `domain`, `application`, `infrastructure`, `api`, `mongo-driver` |
| `api` | `api` |
| `config` | `config`, `mongo-driver` |
| `middleware` | `middleware`, `api`, `gin` |
| `graph` | `domain`, `application`, `graph`, `gqlgen` |

### Dentro del contenedor de desarrollo

```bash
docker compose exec app arch-go
```

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

## Regenerar codigo GraphQL

Si modificas el schema en `graph/schema/*.graphqls`, regenera el codigo:

```bash
~/go/bin/gqlgen generate
```

Esto actualiza:

- `graph/generated.go` — servidor y marshaling
- `graph/model/models_gen.go` — tipos input auto-generados
- `graph/*.resolvers.go` — stubs de resolvers (preserva codigo existente)

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
├── gqlgen.yml
├── plan.md
├── README.md
├── api/
│   ├── errors/
│   │   └── codes.go
│   └── response/
│       └── response.go
├── config/
│   └── mongo.go
├── docs/
│   ├── docs.go
│   ├── error-codes.md
│   ├── swagger.json
│   └── swagger.yaml
├── graph/
│   ├── schema/
│   │   └── producto.graphqls
│   ├── resolver.go
│   ├── producto.resolvers.go
│   ├── generated.go
│   └── model/
│       └── models_gen.go
├── internal/
│   ├── application/
│   │   └── service/
│   │       ├── producto_service.go
│   │       ├── producto_service_test.go
│   │       ├── categoria_service.go
│   │       └── inventario_service.go
│   ├── domain/
│   │   ├── errors.go
│   │   ├── model/
│   │   │   ├── producto.go
│   │   │   ├── categoria.go
│   │   │   └── inventario.go
│   │   └── port/
│   │       ├── producto_repository.go
│   │       ├── categoria_repository.go
│   │       └── inventario_repository.go
│   └── infrastructure/
│       ├── http/
│       │   ├── producto_handler.go
│       │   └── producto_handler_test.go
│       └── persistence/
│           ├── mongodb/
│           │   ├── producto_repository.go
│           │   ├── categoria_repository.go
│           │   └── inventario_repository.go
│           └── postgres/          # Legacy (referencia)
│               └── producto_repository.go
├── middleware/
│   └── ratelimit.go
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.dev
├── main.go
├── go.mod
├── go.sum
├── routes.go
├── tools.go
└── swagger_test.go
```

## Notas tecnicas

- Se usa `uuid` como clave primaria en todas las entidades.
- MongoDB usa `_id` como string (UUID string) para los IDs.
- Soft delete via campo `deleted_at` (timestamp) en todas las colecciones.
- Las relaciones se resuelven en GraphQL via field resolvers (`Producto.categoria`, `Producto.inventario`, `Categoria.productos`).
- REST se mantiene como fallback junto a GraphQL.
- **Arquitectura hexagonal**: las dependencias fluyen hacia adentro (infraestructura → dominio).
- **Inyeccion de dependencias**: `main.go` instancia `repo → service → resolver/handler` y los inyecta.
- **Paginacion REST**: `GET /productos` acepta `?page=N&limit=N` (default: page=1, limit=20, max=100).
- **Paginacion GraphQL**: `productos(page: Int, limit: Int)` con mismos defaults.
- **Colecciones MongoDB**: `productos`, `categorias`, `inventarios`.
