package model

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type Message struct {
	Id          string
	Body        string `json:"body"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
	Caption     string `json:"caption"`
	Document    string `json:"document"`
	FileName    string `json:"filename"`
	ContextInfo waProto.ContextInfo
}

type MessageResponse struct {
	Id        string      `json:"id"`
	Details   string      `json:"details"`
	Timestamp interface{} `json:"timestamp"`
}
