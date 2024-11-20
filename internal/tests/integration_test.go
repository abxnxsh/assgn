
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"mngmnt/internal/database"
    "log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	
	"mngmnt/internal/handlers"
	"mngmnt/internal/models"
)

func TestMain(m *testing.M) {
  
    if err := database.ConnectDB(); err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
   
    m.Run()
}
func TestProductCreationIntegration(t *testing.T) {

	product := models.Product{
		UserID:          1,
		ProductName:     "New Product",
		ProductDescription: "New description",
		ProductImages:   []string{"http://example.com/image.jpg"},
		ProductPrice:    19.99,
	}

	err := database.InsertProduct(&product)
	assert.NoError(t, err)
	req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(product.ID), nil)
	assert.NoError(t, err)


	recorder := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/products/:id", handlers.GetProductByID)
	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var responseProduct models.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &responseProduct)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, responseProduct.ID)
	assert.Equal(t, product.ProductName, responseProduct.ProductName)
}
