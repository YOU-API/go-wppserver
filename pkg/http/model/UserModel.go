package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `sql:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Type     string    `json:"type"`
	Password string    `json:"password"`
	Status   string    `json:"status"`
}

func (u *User) IsAdmin() bool {
	return u.Type == "admin"
}

func (u *User) CanAuthorization(tid uuid.UUID) bool {
	if u.IsAdmin() || u.Id == tid {
		return true
	}

	return false
}
