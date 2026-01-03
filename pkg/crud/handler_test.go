package crud

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	crud "github.com/luanapp/gin-example/pkg/crud/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/luanapp/gin-example/pkg/model"
)

var (
	id, _       = uuid.NewUUID()
	mockSpecies = []model.Species{
		{Id: id.String(), ScientificName: "Phyllobates terribilis", Genus: "Phyllobates", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Allobates alessandroi", Genus: "Allobates", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega bilinguis", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega andina", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega boehmei", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
	}
)

func TestNewHandler(t *testing.T) {
	repo := &repository[model.Species]{conn: &pgxpool.Pool{}}
	got := NewHandler[model.Species](repo)
	expected := &Handler[model.Species]{repo}
	assert.True(t, reflect.DeepEqual(got, expected))
}

func TestHandler_GetAll_StatusOK(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	repo.EXPECT().GetAll().Return(mockSpecies, nil)

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.GET("/species", handler.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/species", nil)
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_GetAll_Status500(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	repo.EXPECT().GetAll().Return([]model.Species{}, errors.New("some unexpected error"))

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.GET("/species", handler.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/species", nil)
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_GetById_StatusOK(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	repo.EXPECT().GetById(specie1.Id).Return(&specie1, nil)

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.GET("/species/:id", handler.GetById)

	req := httptest.NewRequest(http.MethodGet, "/species/"+specie1.Id, nil)
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	expected, _ := json.Marshal(specie1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, rec.Body.Bytes())

}

func TestHandler_GetById_Status500(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	repo.EXPECT().GetById(specie1.Id).Return(nil, errors.New("some unexpected error"))

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.GET("/species/:id", handler.GetById)

	req := httptest.NewRequest(http.MethodGet, "/species/"+specie1.Id, nil)
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_Create_StatusOK(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	expected, _ := json.Marshal(specie1)
	repo.EXPECT().Save(&specie1).Return(&specie1, nil)

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.POST("/species", handler.Save)

	req := httptest.NewRequest(http.MethodPost, "/species", bytes.NewReader(expected))
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, string(expected), rec.Body.String())
}

func TestHandler_Create_StatusBadRequest(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.POST("/species", handler.Save)

	req := httptest.NewRequest(http.MethodPost, "/species", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_Create_Status500(t *testing.T) {
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	expected, _ := json.Marshal(specie1)
	repo.EXPECT().Save(&specie1).Return(&specie1, errors.New("some unexpected error"))

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.POST("/species", handler.Save)

	req := httptest.NewRequest(http.MethodPost, "/species", bytes.NewReader(expected))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_Update_StatusOK(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	expected, _ := json.Marshal(specie1)
	repo.EXPECT().Update(specie1.Id, &specie1).Return(nil)

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.PUT("/species/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/species/"+specie1.Id, bytes.NewReader(expected))
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestHandler_Update_StatusBadRequest(t *testing.T) {
	// Arrange
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.PUT("/species/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/species/"+specie1.Id, bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_Update_Status500(t *testing.T) {
	repo := crud.NewMockRepository[model.Species](t)
	specie1 := mockSpecies[0]
	expected, _ := json.Marshal(specie1)
	repo.EXPECT().Update(specie1.Id, &specie1).Return(errors.New("some unexpected error"))

	e := gin.Default()
	handler := NewHandler[model.Species](repo)
	e.PUT("/species/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/species/"+specie1.Id, bytes.NewReader(expected))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
