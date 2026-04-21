package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed docs/swagger.json
var swaggerJSON []byte

func TestSwaggerDocumentation(t *testing.T) {
	// Parse swagger.json to validate it's valid JSON
	var spec map[string]interface{}
	err := json.Unmarshal(swaggerJSON, &spec)
	assert.NoError(t, err, "swagger.json should be valid JSON")

	// Verify required top-level fields
	assert.NotNil(t, spec["swagger"], "swagger.json should have 'swagger' field")
	assert.NotNil(t, spec["info"], "swagger.json should have 'info' field")
	assert.NotNil(t, spec["paths"], "swagger.json should have 'paths' field")

	// Verify all endpoints are documented
	paths, ok := spec["paths"].(map[string]interface{})
	assert.True(t, ok, "paths should be a map")

	expectedEndpoints := []string{
		"/api/v1/productos",
		"/api/v1/productos/{id}",
	}

	for _, endpoint := range expectedEndpoints {
		assert.Contains(t, paths, endpoint, "endpoint %s should be documented", endpoint)
	}

	// Verify GET /api/v1/productos is documented
	productsList := paths["/api/v1/productos"].(map[string]interface{})
	assert.Contains(t, productsList, "get", "GET /api/v1/productos should be documented")

	// Verify GET /api/v1/productos/{id} is documented
	productByID := paths["/api/v1/productos/{id}"].(map[string]interface{})
	assert.Contains(t, productByID, "get", "GET /api/v1/productos/{id} should be documented")
	assert.Contains(t, productByID, "put", "PUT /api/v1/productos/{id} should be documented")
	assert.Contains(t, productByID, "delete", "DELETE /api/v1/productos/{id} should be documented")

	// Verify Producto model is documented
	definitions, ok := spec["definitions"].(map[string]interface{})
	assert.True(t, ok, "definitions should exist")
	assert.Contains(t, definitions, "models.Producto", "Producto model should be documented")

	t.Logf("✓ Swagger documentation is complete with %d endpoints", len(expectedEndpoints))
}

func TestSwaggerJSONNotEmpty(t *testing.T) {
	assert.Greater(t, len(bytes.TrimSpace(swaggerJSON)), 0, "swagger.json should not be empty")
}
