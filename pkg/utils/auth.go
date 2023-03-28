package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"wppserver/pkg/http/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = "your-256-bit-secret" // wppserver-rest-api

func MakeRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func MakeRandomString() string {
	uuid, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	return strings.Replace(uuid.String(), "-", "", -1)
}

func MakeHMAC256(message string) string {
	hmacSecret := []byte(jwtKey)

	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := []byte(jwtKey)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func MakeTokenFromUUID(uuid uuid.UUID) string {
	hmacSecret := []byte(jwtKey)

	expirationTime := time.Now().Add(60 * 15 * time.Second)

	claims := &model.Claims{
		ID: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSecret)

	if err != nil {
		panic(err)
	}

	return tokenString
}

func MakeToken(userId uuid.UUID, escope string, args ...uuid.UUID) string {
	hmacSecret := []byte(jwtKey)

	apiKeyId := uuid.Nil
	if len(args) != 0 {
		apiKeyId = args[0]
	}

	expirationTime := time.Now().Add(60 * 15 * time.Second)
	claims := model.Claims{
		ApiKeyId: apiKeyId,
		UserId:   userId,
		Scope:    escope,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSecret)

	if err != nil {
		log.Panic(err)
	}

	return tokenString
}

func GetRequestAuth(db *sql.DB, r *http.Request) (model.Auth, bool) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	auth := model.Auth{}
	reqToken := r.Header.Get("Authorization")

	userAuthorization, okGetUserFromClaims := GetUserFromClaims(reqToken, db)
	scopeAuthorization, okGetScopesFromClaims := GetScopeFromClaims(reqToken, db)

	if !okGetUserFromClaims || !okGetScopesFromClaims {
		return auth, false
	}

	auth.User = userAuthorization
	auth.Scope = scopeAuthorization

	return auth, true
}

func GetUserFromClaims(reqToken string, db *sql.DB) (model.User, bool) {
	user := model.User{}

	claims, okExtractUUID := ExtractClaims(reqToken)
	if !okExtractUUID {
		return user, false
	}

	userUUID, errUUID := uuid.Parse(claims["userid"].(string))
	if errUUID != nil {
		return user, false
	}

	err := db.QueryRow("SELECT id, name, email, password, type, status FROM wppserver_users WHERE id=$1", userUUID).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Type, &user.Status)

	if err != nil {
		log.Panicf("query error: %v\n", err)
		return user, false
	}

	return user, true
}

func GetScopeFromClaims(reqToken string, db *sql.DB) (model.Scope, bool) {
	scope := model.Scope{}

	claims, okExtractUUID := ExtractClaims(reqToken)
	if !okExtractUUID {
		return scope, false
	}

	scope.List = strings.Split(claims["Scope"].(string), " ")

	return scope, true
}

func GetRequestToken(db *sql.DB, r *http.Request) (string, bool) {
	var userId uuid.UUID
	var scope string
	var token string

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	granType := r.Form.Get("grant_type")

	if granType == "client" {

		clientId := r.Form.Get("client_id")
		clientSecret := r.Form.Get("client_secret")

		err = db.QueryRowContext(ctx, "SELECT id, userid, scope FROM wppserver_apikeys WHERE clientid=$1 AND clientsecret=$2 AND status='enabled'", clientId, MakeHMAC256(clientSecret)).Scan(
			&userId, &scope)

		if err != nil {
			return "", false
		}

		token = MakeToken(userId, scope)
		return token, true

	} else if granType == "password" || granType == "" {

		user := model.User{}
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		err = db.QueryRowContext(ctx, "SELECT id, password, type FROM wppserver_users WHERE email=$1 AND status='enabled'", username).Scan(
			&user.Id, &user.Password, &user.Type)

		if err != nil {
			return "", false
		}

		match := CheckPasswordHash(password, user.Password)
		if !match {
			return "", false
		}

		if user.IsAdmin() {
			token = MakeToken(user.Id, "admin:*")
		} else {
			token = MakeToken(user.Id, "user:*")
		}

		return token, true
	}

	return "", false
}

func BasicAuthRequest(db *sql.DB, r *http.Request) (model.Key, bool) {
	key := model.Key{}

	ck, cs, ok := r.BasicAuth()
	if !ok {
		return key, false
	}

	err := db.QueryRow("SELECT userid FROM wppserver_apikeys WHERE clientid=$1 AND clientsecret=$2 AND status='enabled'", ck, MakeHMAC256(cs)).Scan(
		&key.UserId)

	if err != nil {
		return key, false
	}

	return key, true
}
