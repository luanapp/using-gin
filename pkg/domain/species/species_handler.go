package species

import (
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
		Family         string `json:"family" form:"family" binding:"required"`
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
	}
	c.JSON(http.StatusOK, sps)
}

func (h *Handler) GetById(c *gin.Context) {
	//id := c.Param("id")
	c.JSON(http.StatusOK, Species{})
}

func (h *Handler) Save(c *gin.Context) {
	sp := new(Species)
	err := c.Bind(sp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Species{})
}
