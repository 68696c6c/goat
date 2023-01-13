package models

type Model interface {
	*Organization | *User
}
