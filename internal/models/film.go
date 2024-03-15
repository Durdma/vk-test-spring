package models

import "github.com/google/uuid"

type Film struct {
	ID          uuid.UUID    `json:"id,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Date        string       `json:"date"`
	Rating      float64      `json:"rating"`
	Actors      []FilmActors `json:"actors"`
}

type FilmActors struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	SecondName string    `json:"second_name"`
	Patronymic string    `json:"patronymic"`
}
