package model

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
)

type Claims struct {
	ID       uuid.UUID `json:"id"`
	ApiKeyId uuid.UUID `json:"apikeyid"`
	UserId   uuid.UUID `json:"userid"`
	Scope    string
	jwt.StandardClaims
}
