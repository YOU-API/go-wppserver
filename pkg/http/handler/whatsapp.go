package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"wppserver/pkg/http/model"
	"wppserver/pkg/utils"
	"wppserver/pkg/whatsapp"

	uuid "github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/vincent-petithory/dataurl"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// Gets QR code encoded in Base64
func LoginDevice(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		deviceId, err := uuid.NewUUID()
		if err != nil {
			log.Panic(err)
		}
		device, ok = ds.RegisterNew(deviceId)
		if !ok {
			log.Panic("Device create error")
		}

		stmt, err := db.Prepare("INSERT INTO wppserver_devices(id, userid, jid, Connected) VALUES($1,$2,$3,$4);")
		if err != nil {
			log.Panicf("query error: %v\n", err)
			respondError(w, http.StatusInternalServerError, "")
			return
		}

		result, err := stmt.Exec(device.Id, device.UserId, "", "yes")
		if err != nil && result != nil {
			log.Panicf("query error: %v\n", err)
			respondError(w, http.StatusBadRequest, "No rows affected")
			return
		}

		device.StartClient(ds.Container)
	}

	if device.Client.Store.ID == nil {
		if device.Client.IsConnected() {
			device.Client.Disconnect()
		}
		qrChan, err := device.Client.GetQRChannel(context.Background())
		if err != nil {
			// This error means that we're already logged in, so ignore it.
			if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				log.Panicf("Failed to get QR channel: %v\n", err)
				respondError(w, http.StatusInternalServerError, "")
				return
			}
		} else {
			err = device.Client.Connect()
			if err != nil {
				log.Panic(err)
				respondError(w, http.StatusInternalServerError, "Failed to connect client")
				return
			}

			for evt := range qrChan {
				if evt.Event == "code" {
					//qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					image, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
					device.QrCode.Base64QrCode = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
					device.QrCode.Expiration = time.Now().UTC().Add(time.Second * time.Duration(30)).Format("2006-01-02T15:04:05-0700")
					break
				}
			}
			respondJSON(w, http.StatusOK, device.QrCode)
			return
		}
	} else {
		if !device.Client.IsConnected() {
			err := device.Client.Connect()
			if err != nil {
				log.Panic(err)
				respondError(w, http.StatusInternalServerError, "Failed to connect client")
				return
			}
		}

		respondJSON(w, http.StatusNoContent, nil)
		return
	}
}

// Logs out device from Whatsapp (requires to scan QR next time)
func LogoutDevice(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	} else if device.Client.IsLoggedIn() {
		if !device.Client.IsConnected() {
			device.Client.Connect()
		}

		err := device.Client.Logout()
		if err != nil {
			log.Panicf("Could not perform logout: %v\n", err)
			respondError(w, http.StatusInternalServerError, "Could not perform logout")
			return
		}

		if ds.Remove(device.UserId) {
			log.Printf("Unable to remove device from memory: %v\n", device)
		}

		log.Printf("Logged out: %v\n", device)
		respondJSON(w, http.StatusOK, nil)
		return
	}
}

// Status from Device websocket
func Status(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	type Status struct {
		IsConnected bool `json:"isconnected"`
		IsLogged    bool `json:"islogged"`
	}

	response := Status{
		IsConnected: device.Client.IsConnected(),
		IsLogged:    device.Client.IsLoggedIn(),
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Connects to WppServer
func Connect(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	device.Client = device.StartClient(ds.Container)

	if device.Client.Store.ID == nil {
		respondError(w, http.StatusInternalServerError, "No logged")
		return
	}

	err := device.Client.Connect()
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "Failed to Connect")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
	return
}

// Disconnects from WppServer websocket, does not log out device
func Disconnect(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	device.Client.Disconnect()

	respondJSON(w, http.StatusNoContent, nil)
	return
}

// Sends a regular text message
func SendText(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	message := model.Message{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&message); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if message.Phone == "" {
		respondError(w, http.StatusBadRequest, "Missing Phone in Payload")
		return
	}

	if message.Body == "" {
		respondError(w, http.StatusBadRequest, "Missing Body in Payload")
		return
	}

	message.Id = whatsmeow.GenerateMessageID()
	recipient, ok := device.ParseJID(message.Phone)

	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &message.Body,
		},
	}

	if message.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*message.ContextInfo.StanzaId),
			Participant:   proto.String(*message.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	result, err := device.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		log.Panicf("Send message: %v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.MessageResponse{
		Id:        message.Id,
		Details:   "Sent",
		Timestamp: result,
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Sends a image message
func SendImage(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	message := model.Message{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&message); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if message.Phone == "" {
		respondError(w, http.StatusBadRequest, "Missing Phone in Payload")
		return
	}

	if message.Image == "" {
		respondError(w, http.StatusBadRequest, "Missing Image in Payload")
		return
	}

	dataURL, err := dataurl.DecodeString(message.Image)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Could not decode base64 encoded data from payload")
		return
	}

	mediaType := strings.ToLower(dataURL.MediaType.ContentType())
	if mediaType != "image/png" && mediaType != "image/jpg" && mediaType != "image/jpeg" && mediaType != "image/webp" {
		log.Panicf("Send image: %v\n", err)
		respondError(w, http.StatusBadRequest, "Unsupported image format. Only use JPG, JPEG, PNG, WEBP Images")
		return
	}

	uploaded, err := device.Client.Upload(context.Background(), dataURL.Data, whatsmeow.MediaImage)
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	message.Id = whatsmeow.GenerateMessageID()
	recipient, ok := device.ParseJID(message.Phone)

	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(message.Caption),
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(dataURL.Data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(dataURL.Data))),
		},
	}

	if message.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*message.ContextInfo.StanzaId),
			Participant:   proto.String(*message.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	result, err := device.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	response := model.MessageResponse{
		Id:        message.Id,
		Details:   "Sent",
		Timestamp: result,
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Sends a document message
func SendDocument(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	message := model.Message{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&message); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if message.Phone == "" {
		respondError(w, http.StatusBadRequest, "Missing Phone in Payload")
		return
	}

	if message.Document == "" {
		respondError(w, http.StatusBadRequest, "Missing Document in Payload")
		return
	}

	if message.FileName == "" {
		respondError(w, http.StatusBadRequest, "Missing FileName in Payload")
		return
	}

	dataURL, err := dataurl.DecodeString(message.Document)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Could not decode base64 encoded data from payload")
		return
	}

	mediaType := strings.ToLower(dataURL.MediaType.ContentType())
	if mediaType != "application/pdf" && mediaType != "application/docx" && mediaType != "application/xlsx" && mediaType != "application/pptx" && mediaType != "application/doc" && mediaType != "application/xls" && mediaType != "application/ppt" && mediaType != "application/txt" && mediaType != "application/docm" && mediaType != "application/xlsm" && mediaType != "application/pptm" && mediaType != "application/rtf" && mediaType != "application/csv" && mediaType != "application/tsv" && mediaType != "application/zip" && mediaType != "application/x-rar-compressed" && mediaType != "application/octet-stream" {
		log.Panicf("Send document: %v\n", err)
		respondError(w, http.StatusBadRequest, "Unsupported image format. Only use JPG, JPEG, PNG, WEBP Images")
		return
	}

	uploaded, err := device.Client.Upload(context.Background(), dataURL.Data, whatsmeow.MediaDocument)
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	message.Id = whatsmeow.GenerateMessageID()
	recipient, ok := device.ParseJID(message.Phone)

	msg := &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			FileName:      &message.FileName,
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(dataURL.Data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(dataURL.Data))),
		},
	}

	if message.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*message.ContextInfo.StanzaId),
			Participant:   proto.String(*message.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	result, err := device.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.MessageResponse{
		Id:        message.Id,
		Details:   "Sent",
		Timestamp: result,
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Scraping BY phone number
func ScrapingPhones(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	type checkPhones struct {
		Phones []string `json:"phones"`
	}

	phones := checkPhones{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&phones); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	isOnWhatsapp, err := device.Client.IsOnWhatsApp(phones.Phones)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	var jids []types.JID
	response := []model.PhoneInfo{}

	// set jids to registered phones
	// set a not registered phones
	for _, item := range isOnWhatsapp {
		if item.IsIn {
			jids = append(jids, item.JID)
			log.Print(item.JID)
		} else {
			// set not registered phones
			response = append(response, model.PhoneInfo{Phone: item.Query})
		}
	}

	// get profile info
	infoPhones, err := device.Client.GetUserInfo(jids)
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	// check if the phone is a contact
	// get the profile picture
	for jid, userInfo := range infoPhones {
		phoneNumber := strings.Split(jid.String(), "@")[0]
		phoneNumber = strings.Split(phoneNumber, ".")[0]

		finfo := model.PhoneInfo{
			Jid:        jid,
			Phone:      phoneNumber,
			InWhatsapp: true,
		}

		cdata, err := device.Client.Store.Contacts.GetContact(jid)
		finfo.IsContact = cdata.Found
		if finfo.IsContact {
			finfo.PessoalName = cdata.PushName
			finfo.BusinessName = cdata.BusinessName
		} else if userInfo.VerifiedName != nil {
			finfo.PessoalName = userInfo.VerifiedName.Details.GetVerifiedName()
		}

		picture, err := device.Client.GetProfilePictureInfo(jid, nil)
		if err == nil {
			finfo.PictureURL = picture.URL
		}

		response = append(response, finfo)
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Get all contacts
func GetContacts(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	result, err := device.Client.Store.Contacts.GetAllContacts()
	if err != nil {
		log.Panic(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	response := []model.PhoneInfo{}
	for cjid, cdata := range result {
		phoneNumber := strings.Split(cjid.String(), "@")[0]
		phoneNumber = strings.Split(phoneNumber, ".")[0]

		cinfo := model.PhoneInfo{
			Phone:        phoneNumber,
			PessoalName:  cdata.PushName,
			BusinessName: cdata.BusinessName,
		}

		picture, err := device.Client.GetProfilePictureInfo(cjid, nil)
		if err == nil {
			cinfo.PictureURL = picture.URL
		}

		response = append(response, cinfo)
	}

	respondJSON(w, http.StatusOK, response)
	return
}

// Create WebHook
func CreateWebhook(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	webhook := model.WebHook{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&webhook); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if webhook.URL == "" {
		respondError(w, http.StatusBadRequest, "Missing URL in Payload")
		return
	}

	stmt, err := db.Prepare("INSERT INTO wppserver_webhooks(id, userid, deviceid, description, url, secret, status) VALUES($1,$2,$3,$4,$5,$6,$7);")
	if err != nil {
		log.Panicf("query prepare error: %v\n", err)
	}

	webhook.Id, err = uuid.NewUUID()
	if err != nil {
		log.Panicf("generate uuid error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	webhook.UserId = device.UserId
	webhook.DeviceId = device.Id

	res, err := stmt.Exec(webhook.Id, webhook.UserId, webhook.DeviceId, webhook.Description, webhook.URL, webhook.Secrete, "enabled")
	if err != nil && res != nil {
		log.Panicf("query exec error: %v\n", err)
	}

	log.Printf("Key registed %q\n", webhook.Id)
	respondJSON(w, http.StatusCreated, webhook)
}

// Delete WebHook
func DeleteWebhook(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	webhook := model.WebHook{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&webhook); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM wppserver_webhooks WHERE id=$1 AND userid=$2;", webhook.Id, device.UserId)
	if err != nil {
		log.Panicf("query exec error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusConflict, "Target not found")
		return
	}

	log.Printf("Webhook deleted %q\n", webhook.Id)
	respondJSON(w, http.StatusNoContent, nil)
}

// Gets WebHook
func GetWebhooks(db *sql.DB, ds *whatsapp.Devices, w http.ResponseWriter, r *http.Request) {
	var device *whatsapp.Device

	auth, ok := utils.GetRequestAuth(db, r)
	if !ok || !auth.Scope.Contain("whatsapp") {
		respondError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	device, ok = ds.Get(auth.User.Id)
	if !ok {
		respondError(w, http.StatusInternalServerError, "No session")
		return
	}

	webhook := model.WebHook{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&webhook); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, userid, deviceid, description, url, secret, status FROM wppserver_webhooks WHERE userid=$1", device.UserId)
	if err != nil {
		log.Panicf("query select error: %v\n", err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	webhooks := make([]model.WebHook, 0)
	for rows.Next() {
		var webhooksData model.WebHook
		if err := rows.Scan(&webhooksData.Id, &webhooksData.UserId, &webhooksData.DeviceId, &webhooksData.Description, &webhooksData.URL, &webhooksData.Secrete, &webhooksData.Status); err != nil {
			log.Panic(err)
			respondError(w, http.StatusInternalServerError, "")
			return
		}
		webhooks = append(webhooks, webhooksData)
	}
	respondJSON(w, http.StatusOK, webhooks)

}
