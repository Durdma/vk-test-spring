package models

type Film struct {
	Name        string
	Description string
	Date        string
	Rating      float64
	Actors      []Actor
}
