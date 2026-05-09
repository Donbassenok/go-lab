package repository

import "github.com/Donbassenok/go-lab/internal/model"

type PlantRepo interface {
	Create(plant model.Plant) (int, error)
	GetAll() ([]model.Plant, error)
	GetByID(id int) (model.Plant, error)
	Update(id int, plant model.Plant) error
	Patch(id int, updates map[string]interface{}) error
	Delete(id int) error
}
