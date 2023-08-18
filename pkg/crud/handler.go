package crud

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

	Handler[T any] struct {
		repo Repository[T]
	}
)

// NewHandler returns a new repository object
func NewHandler[T any](repo Repository[T]) Handlerer {
	return &Handler[T]{
		repo: repo,
	}
}

func DefaultHandler[T any]() Handlerer {
	return NewHandler(defaultRepository[T]())
}

func (h *Handler[T]) GetAll(c *gin.Context) {
	entities, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities)
}

func (h *Handler[T]) GetById(c *gin.Context) {
	id := c.Param("id")
	sp, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("species with id %s not found", id)})
		return
	}

	c.JSON(http.StatusOK, sp)
}

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

func (h *Handler[T]) Update(c *gin.Context) {
	id := c.Param("id")

	e := new(T)
	err := c.Bind(e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.Update(id, e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *Handler[T]) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
