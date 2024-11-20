// test trial
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mngmnt/internal/handlers"
	"mngmnt/internal/models"
)

func TestCreateProduct(t *testing.T) {
	
	product := models.Product{
		UserID:          1,
		ProductName:     "Test Product",
		ProductDescription: "Test Description",
		ProductImages:   []string{"http://example.com/image1.jpg"},
		ProductPrice:    29.99,
	}

	body, err := json.Marshal(product)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewReader(body))
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/products", handlers.CreateProduct)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
