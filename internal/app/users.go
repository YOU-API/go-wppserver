package app

import (
	"net/http"
	"wppserver/pkg/http/handler"
)

func (a *App) getAllUsers(w http.ResponseWriter, r *http.Request) {
	handler.GetAllUsers(a.DB, w, r)
}

func (a *App) registerUser(w http.ResponseWriter, r *http.Request) {
	handler.RegisterUser(a.DB, w, r, a.Config)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}

func (a *App) findUsers(w http.ResponseWriter, r *http.Request) {
	handler.FindUsers(a.DB, w, r)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	handler.DeleteUser(a.DB, w, r)
}

func (a *App) statusUser(w http.ResponseWriter, r *http.Request) {
	handler.StatusUser(a.DB, w, r)
}

func (a *App) passwordUser(w http.ResponseWriter, r *http.Request) {
	handler.PasswordUser(a.DB, w, r)
}

func (a *App) createUserKey(w http.ResponseWriter, r *http.Request) {
	handler.CreateUserKey(a.DB, w, r)
}

func (a *App) deleteUserKey(w http.ResponseWriter, r *http.Request) {
	handler.DeleteUserKey(a.DB, w, r)
}

func (a *App) getUserKeys(w http.ResponseWriter, r *http.Request) {
	handler.GetUserKeys(a.DB, w, r)
}
