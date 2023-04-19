package models

type Card struct {
	Id    int   `json:"id"`
	Count uint8 `json:"count" validate:"required"`
	Fill  uint8 `json:"fill" validate:"required"`
	Shape uint8 `json:"shape" validate:"required"`
	Color uint8 `json:"color" validate:"required"`
}
