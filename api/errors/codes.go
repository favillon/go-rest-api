package errors

import "net/http"

// Code define los codigos estables que la API devuelve en errores.
type Code string

const (
	DBNotInitialized Code = "DB_NOT_INITIALIZED"
	Database         Code = "DATABASE_ERROR"
	InvalidParam     Code = "INVALID_PARAMETER"
	Validation       Code = "VALIDATION_ERROR"
	NotFound         Code = "RESOURCE_NOT_FOUND"
)

// Metadata describe cada codigo del catalogo.
type Metadata struct {
	HTTPStatus  int
	Description string
}

// Catalog es el catalogo oficial de codigos de error de la API.
var Catalog = map[Code]Metadata{
	DBNotInitialized: {
		HTTPStatus:  http.StatusInternalServerError,
		Description: "La conexion a base de datos no fue inicializada en el servidor.",
	},
	Database: {
		HTTPStatus:  http.StatusInternalServerError,
		Description: "Error al ejecutar una operacion de base de datos.",
	},
	InvalidParam: {
		HTTPStatus:  http.StatusBadRequest,
		Description: "Parametro invalido en la ruta o query string.",
	},
	Validation: {
		HTTPStatus:  http.StatusBadRequest,
		Description: "El payload no cumple las reglas de validacion.",
	},
	NotFound: {
		HTTPStatus:  http.StatusNotFound,
		Description: "El recurso solicitado no existe.",
	},
}
