package model

import (
	uuid "github.com/google/uuid"
	"go.mau.fi/whatsmeow/types"
)

type UserDevice struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"userid"`
	Phone     string    `json:"phone"`
	Jid       types.JID `json:"jid"`
	Connected string    `json:"connected"`
}
