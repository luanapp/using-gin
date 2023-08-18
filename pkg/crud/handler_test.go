package crud

import (
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
	assert.Equal(t, rec.Code, http.StatusOK)
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
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}
