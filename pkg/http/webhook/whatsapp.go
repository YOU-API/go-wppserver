package webhook

import (
	"database/sql"
	"log"
	"wppserver/pkg/http/model"
	"wppserver/pkg/whatsapp"

	"github.com/google/uuid"
)

// Events from Whatsapp websocket
func Handler(db *sql.DB, device *whatsapp.Device, rawEvt interface{}) {
	rows, err := db.Query("SELECT id, userid, deviceid, description, url, secret, status FROM wppserver_webhooks WHERE userid=$1", device.UserId)
	if err != nil {
		log.Panicf("query select error: %v\n", err)
		return
	}

	for rows.Next() {
		var webhooksData model.WebHook
		if err := rows.Scan(&webhooksData.Id, &webhooksData.UserId, &webhooksData.DeviceId, &webhooksData.Description, &webhooksData.URL, &webhooksData.Secrete, &webhooksData.Status); err != nil {
			log.Panic(err)
			return
		}
		if webhooksData.Status != "enabled" {
			continue
		}

		type PostData struct {
			Event    string
			UserID   uuid.UUID
			DeviceID uuid.UUID
			State    string
			Data     model.Message
		}

		// switch evt := rawEvt.(type) {
		// case *events.Message:
		// 	response := PostData{
		// 		Event:    "Message",
		// 		UserID:   device.UserId,
		// 		DeviceID: device.Id,
		// 	}

		// 	ImageMessage := evt.Message.GetImageMessage()
		// 	if ImageMessage != nil {

		// 	}
		// }
	}

}
