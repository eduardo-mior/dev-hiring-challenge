package router

import (
	"desafio_ateliware/backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareNotFound(t *testing.T) {

	router := gin.Default()
	setupNotFoundEndPointMiddleware(router)

	server := httptest.NewServer(router)
	defer server.Close()

	request, err := http.NewRequest("GET", "http://localhost:8080/end-point-inexistente ", nil)
	assert.Equal(t, nil, err)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	responseBody := models.Error{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.Equal(t, nil, err)

	assert.Equal(t, 400, response.Code)
	assert.Equal(t, models.Error{Error: "Recurso não encontrado"}, responseBody)

}

func TestPanicRecovery(t *testing.T) {

	router := gin.Default()
	setupPanicRecoveryMiddleware(router)

	server := httptest.NewServer(router)
	defer server.Close()

	router.GET("/panic-recovery", func(c *gin.Context) {
		panic("teste")
	})

	request, err := http.NewRequest("GET", "http://localhost:8080/panic-recovery", nil)
	assert.Equal(t, nil, err)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	responseBody := models.Error{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.Equal(t, nil, err)

	assert.Equal(t, 500, response.Code)
	assert.Equal(t, models.Error{Error: "Ocorreu um erro interno ao tentar realizar a operação"}, responseBody)

}
