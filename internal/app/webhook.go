package app

import (
	"wppserver/pkg/http/webhook"
	"wppserver/pkg/whatsapp"
)

func (a *App) webhook(device *whatsapp.Device, rawEvt interface{}) {
	webhook.Handler(a.DB, device, rawEvt)
}
