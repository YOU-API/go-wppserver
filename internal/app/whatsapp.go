package app

import (
	"net/http"
	"wppserver/pkg/http/handler"
)

func (a *App) loginDevice(w http.ResponseWriter, r *http.Request) {
	handler.LoginDevice(a.DB, a.Devices, w, r)
}

func (a *App) logoutDevice(w http.ResponseWriter, r *http.Request) {
	handler.LogoutDevice(a.DB, a.Devices, w, r)
}

func (a *App) status(w http.ResponseWriter, r *http.Request) {
	handler.Status(a.DB, a.Devices, w, r)
}

func (a *App) connect(w http.ResponseWriter, r *http.Request) {
	handler.Connect(a.DB, a.Devices, w, r)
}

func (a *App) disconnect(w http.ResponseWriter, r *http.Request) {
	handler.Disconnect(a.DB, a.Devices, w, r)
}

func (a *App) createWebhook(w http.ResponseWriter, r *http.Request) {
	handler.CreateWebhook(a.DB, a.Devices, w, r)
}

func (a *App) deleteWebhook(w http.ResponseWriter, r *http.Request) {
	handler.DeleteWebhook(a.DB, a.Devices, w, r)
}

func (a *App) getWebhooks(w http.ResponseWriter, r *http.Request) {
	handler.GetWebhooks(a.DB, a.Devices, w, r)
}

func (a *App) sendText(w http.ResponseWriter, r *http.Request) {
	handler.SendText(a.DB, a.Devices, w, r)
}

func (a *App) sendImage(w http.ResponseWriter, r *http.Request) {
	handler.SendImage(a.DB, a.Devices, w, r)
}

func (a *App) sendDocument(w http.ResponseWriter, r *http.Request) {
	handler.SendDocument(a.DB, a.Devices, w, r)
}

func (a *App) getContacts(w http.ResponseWriter, r *http.Request) {
	handler.GetContacts(a.DB, a.Devices, w, r)
}

func (a *App) scrapingPhones(w http.ResponseWriter, r *http.Request) {
	handler.ScrapingPhones(a.DB, a.Devices, w, r)
}
