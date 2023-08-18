package model

type (
	Species struct {
		Id             string `json:"id" form:"id" example:"996ff476-09bc-45f8-b79d-83b268de2485" db:"id,pk"`
		ScientificName string `json:"scientific_name" form:"scientific_name" binding:"required" example:"Phyllobates terribilis" db:"scientific_name"`
		Genus          string `json:"genus" form:"genus" binding:"required" example:"Phyllobates" db:"genus"`
		Family         string `json:"family" form:"family" binding:"required" example:"Dendrobatidae" db:"family"`
		Order          string `json:"order" form:"order" binding:"required" example:"Anura" db:"order"`
		Class          string `json:"class" form:"class" binding:"required" example:"Amphibia" db:"class"`
		Phylum         string `json:"phylum" form:"phylum" binding:"required" example:"Chordata" db:"phylum"`
		Kingdom        string `json:"kingdom" form:"kingdom" binding:"required" example:"Animalia" db:"kingdom"`
	}
)
