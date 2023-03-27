package model

import "github.com/google/uuid"

type WebHook struct {
	Id          uuid.UUID `sql:"id;default:uuid_generate_v4()"`
	UserId      uuid.UUID `sql:"userid;default:uuid_generate_v4()"`
	DeviceId    uuid.UUID `sql:"userid;default:uuid_generate_v4()"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	URL         string    `json:"url"`
	Events      string    `json:"events"`
	Secrete     string    `json:"secrete"`
}
