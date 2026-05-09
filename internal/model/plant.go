package model

type Plant struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required,min=2"`
	Species string `json:"species" validate:"required"`
	Age     int    `json:"age" validate:"required,gt=0"`
}
