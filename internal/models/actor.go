package models

import (
	"github.com/google/uuid"
)

type Actor struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	SecondName  string    `json:"second_name"`
	Patronymic  string    `json:"patronymic"`
	Sex         string    `json:"sex"`
	DateOfBirth string    `json:"date_of_birth"`
	Films       []Film    `json:"films"`
}
