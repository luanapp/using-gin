package species

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	Handlerer interface {
		GetAll(c *gin.Context)
		GetById(c *gin.Context)
		Save(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}

	Handler struct {
		repo *repository
	}

	Species struct {
		Id             string `json:"id" form:"id" example:"996ff476-09bc-45f8-b79d-83b268de2485"`
		ScientificName string `json:"scientific_name" form:"scientific_name" binding:"required" example:"Phyllobates terribilis"`
		Genus          string `json:"genus" form:"genus" binding:"required" example:"Phyllobates"`
		Family         string `json:"family" form:"family" binding:"required" example:"Dendrobatidae"`
		Order          string `json:"order" form:"order" binding:"required" example:"Anura"`
		Class          string `json:"class" form:"class" binding:"required" example:"Amphibia"`
		Phylum         string `json:"phylum" form:"phylum" binding:"required" example:"Chordata"`
		Kingdom        string `json:"kingdom" form:"kingdom" binding:"required" example:"Animalia"`
	}
)

// NewHandler returns a new repository object
func NewHandler(repo *repository) Handlerer {
	return &Handler{
		repo: repo,
	}
}

// GetAll retrieve all species
// @Summary Get all species
// @Description Retrieves all species from database
// @Tags species
// @Accept  json
// @Produce  json
// @Success 200 {array} Species
// @Failure 500 {object} object
// @Router /species [get]
func (h *Handler) GetAll(c *gin.Context) {
	sps, err := h.repo.getAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sps)
}

// GetById retrieve species with the given id
// @Summary Get species with the given id
// @Description Retrieves the species from database with the given id
// @Tags species
// @Accept  json
// @Produce  json
// @Param id path string true "Species ID"
// @Success 200 {object} Species
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /species/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	id := c.Param("id")
	sp, err := h.repo.getById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("species with id %s not found", id)})
		return
	}

	c.JSON(http.StatusOK, sp)
}

// Save saves a species into database
// @Summary Saves the given species
// @Description Saves the species information into database
// @Tags species
// @Accept  json
// @Produce  json
// @Param species body Species true "Species Payload"
// @Success 201 {object} Species
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /species [post]
func (h *Handler) Save(c *gin.Context) {
	sp := new(Species)
	err := c.Bind(sp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedSP, err := h.repo.save(sp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, savedSP)
}

// Update updates the species with the given id in the database
// @Summary Updates the species with the given id
// @Description Updates the species information with the given id in the database
// @Tags species
// @Accept  json
// @Produce  json
// @Param species body Species true "Species Payload"
// @Param id path string true "Species ID"
// @Success 202
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /species/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	sp := new(Species)
	err := c.Bind(sp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sp.Id = id
	err = h.repo.update(sp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}

// Delete deletes the species with the given id from database
// @Summary Deletes the species with the given id
// @Description Deletes the species information with the given id from database
// @Tags species
// @Accept  json
// @Produce  json
// @Param id path string true "Species ID"
// @Success 202
// @Failure 500 {object} object
// @Router /species/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
