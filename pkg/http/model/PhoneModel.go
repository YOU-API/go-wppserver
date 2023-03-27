package model

import "go.mau.fi/whatsmeow/types"

type PhoneInfo struct {
	Jid          types.JID `json:"jid"`
	InWhatsapp   bool      `json:"onwhatsapp"`
	Phone        string    `json:"phone"`
	PessoalName  string    `json:"pessoalname"`
	BusinessName string    `json:"businessname"`
	PictureURL   string    `json:"pictureurl"`
	IsContact    bool      `json:"iscontact"`
}
