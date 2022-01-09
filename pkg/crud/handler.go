package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/luanapp/gin-example/pkg/model"
)

type (
	Handlerer interface {
		GetAll(c *gin.Context)
		GetById(c *gin.Context)
		Save(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}

	Handler[T model.Model] struct {
		repo Repository[T]
	}
)

// NewHandler returns a new repository object
func NewHandler[T model.Model](repo Repository[T]) Handlerer {
	return &Handler[T]{
		repo: repo,
	}
}

func DefaultHandler[T model.Model]() Handlerer {
	return NewHandler(defaultRepository[T]())
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
func (h *Handler[T]) GetAll(c *gin.Context) {
	entities, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities)
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
func (h *Handler[T]) GetById(c *gin.Context) {
	id := c.Param("id")
	sp, err := h.repo.GetById(id)
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
func (h *Handler[T]) Save(c *gin.Context) {
	e := new(T)
	err := c.Bind(e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedEntity, err := h.repo.Save(e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, savedEntity)
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
func (h *Handler[T]) Update(c *gin.Context) {
	id := c.Param("id")

	e := new(T)
	err := c.Bind(e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	(*e).SetId(id)
	err = h.repo.Update(e)
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
func (h *Handler[T]) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
