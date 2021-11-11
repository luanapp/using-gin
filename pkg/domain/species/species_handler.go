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
	}

	Handler struct {
		repo *repository
	}

	Species struct {
		Id             string `json:"id" form:"id"`
		ScientificName string `json:"scientific_name" form:"scientific_name" binding:"required"`
		Genus          string `json:"genus" form:"genus" binding:"required"`
		Family         string `json:"family" form:"family" binding:"required"`
		Order          string `json:"order" form:"order" binding:"required"`
		Class          string `json:"class" form:"class" binding:"required"`
		Phylum         string `json:"phylum" form:"phylum" binding:"required"`
		Kingdom        string `json:"kingdom" form:"kingdom" binding:"required"`
	}
)

func NewHandler(repo *repository) Handlerer {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	sps, err := h.repo.getAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sps)
}

func (h *Handler) GetById(c *gin.Context) {
	id := c.Param("id")
	sp, err := h.repo.getById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("species with id %s not found", id)})
		return
	}

	c.JSON(http.StatusOK, sp)
}

func (h *Handler) Save(c *gin.Context) {
	sp := new(Species)
	err := c.Bind(sp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedSP, err := h.repo.Save(sp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, savedSP)
}
