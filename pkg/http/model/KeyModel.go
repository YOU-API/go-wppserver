package model

import "github.com/google/uuid"

type Key struct {
	Id              uuid.UUID `sql:"id;default:uuid_generate_v4()"`
	UserId          uuid.UUID `sql:"userid;default:uuid_generate_v4()"`
	Description     string    `json:"description"`
	ClientID        string    `json:"clientid"`
	ClientSecret    string    `json:"clientsecret"`
	TruncatedSecret string    `json:"truncatedsecret"`
}
