package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/google/uuid"

	"wppserver/pkg/config"
	"wppserver/pkg/http/model"
	"wppserver/pkg/utils"
)

func GetAllUsers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.IsAdmin() {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, name, email, type, status FROM wppserver_users")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Type, &user.Status); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	respondJSON(w, http.StatusOK, users)
}

func FindUsers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	searchKeyword := r.Form.Get("keyword")

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.IsAdmin() {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, name, email, type, status FROM wppserver_users WHERE name LIKE '%' || $1 || '%' OR email LIKE '%' || $1 || '%'", searchKeyword)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Type, &user.Status); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	respondJSON(w, http.StatusOK, users)
}

func RegisterUser(db *sql.DB, w http.ResponseWriter, r *http.Request, config *config.Config) {
	user := model.User{Status: "disabled"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	stmt, err := db.Prepare("INSERT INTO wppserver_users(id, name, email, password, type, status) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Panic(err)
	}

	uuid, err := uuid.NewUUID()
	userPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Panic(err)
	}

	if config.SETTINGS.NewUsersStatus == "ENABLED" {
		user.Status = "enabled"
	} else {
		user.Status = "disabled"
	}

	result, err := stmt.Exec(uuid, user.Name, user.Email, userPassword, "agent", user.Status)
	if err != nil && result != nil {
		log.Panic(err)
	}

	err = db.QueryRow("SELECT id, name, email, type, status FROM wppserver_users WHERE id=$1", uuid).Scan(
		&user.Id, &user.Name, &user.Email, &user.Type, &user.Status)

	if err != nil || err == sql.ErrNoRows {
		log.Printf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("User registed %q\n", user.Email)
	respondJSON(w, http.StatusOK, user)
}

func GetUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userTarget := model.User{}

	/*decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userTarget); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()*/

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	err := db.QueryRowContext(ctx, "SELECT id, name, email, type, status FROM wppserver_users WHERE id=$1", auth.User.Id).Scan(
		&userTarget.Id, &userTarget.Name, &userTarget.Email, &userTarget.Type, &userTarget.Status)

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, userTarget)
}

func UpdateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var err error
	var result sql.Result

	ctx := context.Background()

	userTarget := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userTarget); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(userTarget.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	if auth.User.IsAdmin() {
		result, err = db.ExecContext(ctx, "UPDATE wppserver_users SET name=$1,email=$2,type=$3,status=$4 WHERE id=$5;",
			userTarget.Name, userTarget.Email, userTarget.Type, userTarget.Status, userTarget.Id)
	} else {
		result, err = db.ExecContext(ctx, "UPDATE wppserver_users SET name=$1,email=$2 WHERE id=$3;",
			userTarget.Name, userTarget.Email, userTarget.Id)
	}

	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "No rows affected")
		return
	}

	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
		if err != nil {
			log.Panicf("query error: %v\n", err)
			respondError(w, http.StatusBadRequest, "Expected to affect no more than 1 a row")
			return
		}
	}

	log.Printf("User updated %q\n", userTarget.Email)
	respondJSON(w, http.StatusNoContent, nil)
}

func DeleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userTarget := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userTarget); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(userTarget.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM wppserver_users WHERE id=$1;", userTarget.Id)
	if err != nil {
		log.Fatalf("query error: %v\n", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "No rows affected")
		return
	}

	if rows != 1 {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "Expected to affect no more than 1 a row")
		return
	}

	log.Printf("User deleted %q\n", userTarget.Id)
	respondJSON(w, http.StatusNoContent, nil)
}

func StatusUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var err error
	var result sql.Result

	userTarget := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&userTarget); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.IsAdmin() {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	result, err = db.ExecContext(ctx, "UPDATE wppserver_users SET status=$1 WHERE id=$2;", userTarget.Status, userTarget.Id)
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "No rows affected")
		return
	}
	if rows != 1 {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "Expected to affect no more than 1 a row")
		return
	}

	log.Print("User status updated")
	respondJSON(w, http.StatusNoContent, nil)
}

func PasswordUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var err error
	var result sql.Result

	type PasswordUpdate struct {
		Id          uuid.UUID `json:"id"`
		OldPassword string    `json:"oldpassword"`
		NewPassword string    `json:"newpassword"`
	}

	userPassword := PasswordUpdate{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&userPassword); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(userPassword.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	user := model.User{}
	err = db.QueryRowContext(ctx, "SELECT id, name, email, password FROM wppserver_users WHERE id=$1", userPassword.Id).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusBadRequest, "No rows affected")
		return
	}

	match := utils.CheckPasswordHash(userPassword.OldPassword, user.Password)
	if !match {
		respondError(w, http.StatusBadRequest, "Wrong password")
		return
	}

	newPassword, err := utils.HashPassword(userPassword.NewPassword)
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	result, err = db.ExecContext(ctx, "UPDATE wppserver_users SET password=$1 WHERE id=$2;", newPassword, userPassword.Id)
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		respondError(w, http.StatusBadRequest, "No rows affected")
		return
	}
	if rows == 0 {
		respondError(w, http.StatusBadRequest, "Passwords do not match.")
		return
	}

	log.Print("User password updated")
	respondJSON(w, http.StatusNoContent, nil)
}

func CreateUserKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type UserToken struct {
		Token string `json:"token"`
	}

	keyData := model.Key{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&keyData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(keyData.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	stmt, err := db.Prepare("INSERT INTO wppserver_apikeys(id, userid, description, clientid, clientsecret, truncatedsecret) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	keyData.Id = uuid
	keyData.ClientID = "ck_" + utils.MakeRandomString()

	keyData.ClientSecret = "cs_" + utils.MakeRandomString() // Return the unhashed key on creation to be displayed once and never stored.
	keyData.TruncatedSecret = keyData.ClientSecret[0:14]    // first clientsecret characters

	res, err := stmt.Exec(uuid, keyData.UserId, keyData.Description, keyData.ClientID, utils.MakeHMAC256(keyData.ClientSecret), keyData.TruncatedSecret)
	if err != nil && res != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	log.Printf("Key registed %q\n", uuid)
	respondJSON(w, http.StatusCreated, keyData)
}

func DeleteUserKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	keyData := model.Key{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&keyData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(keyData.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM wppserver_apikeys WHERE id=$1 AND userid=$2;", keyData.Id, keyData.UserId)
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusConflict, "Target not found")
		return
	}

	log.Printf("Key deleted %q\n", keyData.Id)
	respondJSON(w, http.StatusNoContent, nil)
}

func GetUserKeys(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userTarget := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userTarget); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.User.CanAuthorization(userTarget.Id) {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, userid, description, clientid, truncatedsecret FROM wppserver_apikeys WHERE userid=$1", userTarget.Id)
	if err != nil {
		log.Panicf("query error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	keys := make([]model.Key, 0)
	for rows.Next() {
		var keyData model.Key
		if err := rows.Scan(&keyData.Id, &keyData.UserId, &keyData.Description, &keyData.ClientID, &keyData.TruncatedSecret); err != nil {
			log.Fatal(err)
		}
		keys = append(keys, keyData)
	}
	respondJSON(w, http.StatusOK, keys)
}
