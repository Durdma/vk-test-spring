package models

type User struct {
	ID   int
	Name string
	Role Role
}

type Role struct {
	ID   int
	Role string
}

type Admin struct {
	User
}
