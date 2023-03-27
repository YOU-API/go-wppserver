package app

import (
	"net/http"
	"wppserver/pkg/http/handler"
)

func (a *App) accessToken(w http.ResponseWriter, r *http.Request) {
	handler.AccessToken(a.DB, w, r)
}

func (a *App) refreshToken(w http.ResponseWriter, r *http.Request) {
	handler.RefreshToken(a.DB, w, r)
}
