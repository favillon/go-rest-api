# Catalogo de codigos de error

Este documento define los codigos de error estables de la API para consumo de clientes.

## Formato de error

```json
{
  "status": "error",
  "message": "No fue posible obtener el producto",
  "error": {
    "code": "DATABASE_ERROR",
    "message": "detalle tecnico"
  }
}
```

## Codigos

| Code | HTTP Status | Descripcion |
| --- | --- | --- |
| DB_NOT_INITIALIZED | 500 | La conexion a base de datos no fue inicializada en el servidor. |
| DATABASE_ERROR | 500 | Error al ejecutar una operacion de base de datos. |
| INVALID_PARAMETER | 400 | Parametro invalido en la ruta o query string. |
| VALIDATION_ERROR | 400 | El payload no cumple las reglas de validacion. |
| RESOURCE_NOT_FOUND | 404 | El recurso solicitado no existe. |

## Fuente de verdad en codigo

La fuente de verdad del catalogo esta en:

- `api/errors/codes.go`
