package model

type Sample struct {
	Id       string `json:"id" form:"id" example:"d3b1d940-6b39-4836-a4a8-00445af31fdf" db:"id,pk"`
	SpecieId string `json:"speciesId" form:"species" example:"996ff476-09bc-45f8-b79d-83b268de2485" db:"specieId,fk,species:id"`
	Type     string `json:"type" form:"type" example:"tissue" db:"type"`
}

func (s Sample) GetId() string {
	return s.Id
}

func (s Sample) SetId(id string) {
	s.Id = id
}
