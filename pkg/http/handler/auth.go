package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"wppserver/pkg/utils"
)

func AccessToken(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type Token struct {
		AccessToken string `json:"access_token"`
	}

	token, ok := utils.GetRequestToken(db, r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	respondJSON(w, http.StatusOK, Token{AccessToken: token})
}

func RefreshToken(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type Token struct {
		AccessToken string `json:"access_token"`
	}

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	token := utils.MakeToken(auth.User.Id, strings.Join(auth.Scope.List, " "))
	respondJSON(w, http.StatusOK, Token{AccessToken: token})

}
