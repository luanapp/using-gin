package crud

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	id, _       = uuid.NewUUID()
	mockSpecies = []Species{
		{Id: id.String(), ScientificName: "Phyllobates terribilis", Genus: "Phyllobates", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Allobates alessandroi", Genus: "Allobates", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega bilinguis", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega andina", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
		{Id: id.String(), ScientificName: "Ameerega boehmei", Genus: "Ameerega", Family: "Dendrobatidae", Order: "Anura", Class: "Amphibia", Phylum: "Chordata", Kingdom: "Animalia"},
	}
)

func TestNewHandler(t *testing.T) {
	repo := &repository{conn: &pgxpool.Pool{}}
	got := NewHandler(repo)
	expected := &Handler{repo}
	assert.True(t, reflect.DeepEqual(got, expected))
}

func TestHandler_GetAll_StatusOK(t *testing.T) {
	// Arrange
	repo := &repositoryMock{}
	repo.On("getAll").Return(mockSpecies, nil)

	e := gin.Default()
	handler := NewHandler(repo)
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
	repo := &repositoryMock{}
	repo.On("getAll").Return([]Species{}, errors.New("some unexpected error"))

	e := gin.Default()
	handler := NewHandler(repo)
	e.GET("/species", handler.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/species", nil)
	req.Header.Set("Content-Type", "application/json")

	// Act
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}
