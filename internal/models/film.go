package models

import "github.com/google/uuid"

type Film struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	Rating      float64   `json:"rating"`
	Actors      []Actor   `json:"actors"`
}
