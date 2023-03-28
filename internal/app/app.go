package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"wppserver/assets"
	"wppserver/pkg/config"
	"wppserver/pkg/http/model"
	"wppserver/pkg/utils"
	"wppserver/pkg/whatsapp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

//App has router and db instances
type App struct {
	DB      *sql.DB
	Router  *mux.Router
	Devices *whatsapp.Devices
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) error {
	var err error

	if a.DB, err = a.DatabaseSetup(config); err != nil {
		return err
	}

	a.Devices = whatsapp.NewDevices()
	a.connectDevices()
	a.Router = mux.NewRouter()
	a.setRouters()

	if err := a.Run(config); err != nil {
		return err
	}

	return nil
}

// Initial database settings
func (a *App) DatabaseSetup(config *config.Config) (*sql.DB, error) {
	// Create sqlite default db
	if config.DB.DbURI == "file:dbdata/wppserver.db?_foreign_keys=true" {
		path, err := utils.GetWorkDir()
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(path + "/dbdata/wppserver.db"); err != nil {
			zipFile, errDir := assets.F.ReadFile("dbdata.zip")
			if errDir != nil {
				return nil, errDir
			}
			if errUnzip := utils.UnzipReader(zipFile, path); errUnzip != nil {
				return nil, errUnzip
			}
		}
	}

	// Open database connection
	db, err := sql.Open(config.DB.Dialect, config.DB.DbURI)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create initial tables in db
	for _, v := range []string{
		"CREATE TABLE IF NOT EXISTS wppserver_users( id VARCHAR(36) NOT NULL, name	TEXT, email	TEXT, type	TEXT, password	TEXT, status	TEXT DEFAULT 'disabled', PRIMARY KEY(id) );",
		"CREATE TABLE IF NOT EXISTS wppserver_apikeys( id VARCHAR(36) NOT NULL, userid VARCHAR(36) NOT NULL, description TEXT, permissions TEXT, consumerkey TEXT, consumersecret TEXT, truncatedsecret	TEXT, scope	TEXT, status TEXT DEFAULT 'enabled', jid TEXT, PRIMARY KEY(id) );",
		"CREATE TABLE IF NOT EXISTS wppserver_devices( id VARCHAR(36) NOT NULL, userid VARCHAR(36), jid	TEXT, qrcode	TEXT, connected	TEXT, PRIMARY KEY(id) );",
		"CREATE TABLE IF NOT EXISTS wppserver_webhooks( id VARCHAR(36) NOT NULL, deviceid VARCHAR(36), description	TEXT, status	TEXT, url	TEXT, events	TEXT, secret	TEXT, PRIMARY KEY(id) );",
	} {
		if _, err := db.Exec(v); err != nil {
			return nil, err
		}
	}

	// Create initial admin user
	if config.AUTH.UserEmail != "" && config.AUTH.UserPassword != "" {
		user := model.User{}
		err := db.QueryRow("SELECT id, name, email, type, status FROM wppserver_users WHERE type=admin").Scan(
			&user.Id, &user.Name, &user.Email, &user.Type, &user.Status)

		if err != nil || err == sql.ErrNoRows {
			stmt, err := db.Prepare("INSERT INTO wppserver_users(id, name, email, password, type, status) VALUES($1,$2,$3,$4,$5,$6);")
			if err != nil {
				return nil, err
			}

			user.Id, err = uuid.NewUUID()
			user.Password, err = utils.HashPassword(config.AUTH.UserPassword)
			if err != nil {
				return nil, err
			}

			result, err := stmt.Exec(user.Id, "root", config.AUTH.UserEmail, user.Password, "admin", "enabled")
			if err != nil && result != nil {
				return nil, err
			} else {
				log.Printf("Initial admin user registed %q\n", user.Id)
			}
		}
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

// Connect devices saved in DB
func (a *App) connectDevices() {
	log.Printf("Init connection to Whatsapp")

	rows, err := a.DB.Query("SELECT id, userid, jid, connected FROM wppserver_devices WHERE connected='yes'")
	if err != nil {
		log.Fatalf("query error: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var device whatsapp.Device

		if err := rows.Scan(&device.Id, &device.UserId, &device.Jid, &device.Connected); err != nil {
			log.Fatalf("rows error: %v\n", err)
		}

		if _, ok := a.Devices.Get(device.Id); !ok {
			log.Printf("Connect to Whatsapp on startup:  %v\n", device.Jid)

			device.FnEvent = a.webhook
			device.StartClient(a.Devices.Container)

			err := device.Client.Connect()
			if err != nil {
				log.Panic(err)
			}

			if device.Client.Store.ID != nil {
				a.Devices.Add(&device)
			} else {
				_, err := a.DB.Exec("DELETE FROM wppserver_devices WHERE id=$1 AND userid=$2;", device.Id, device.UserId)
				if err != nil {
					log.Panicf("query exec error: %v\n", err)
					return
				}
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("query error: %v\n", err)
	}
}

// Set all required routers
func (a *App) setRouters() {
	a.Post("/v1/auth", a.accessToken)
	a.Get("/v1/auth/refresh", a.refreshToken)

	a.Post("/v1/user", a.registerUser)
	a.Get("/v1/user", a.getUser)
	a.Put("/v1/user", a.updateUser)
	a.Delete("/v1/user", a.deleteUser)

	a.Put("/v1/user/status", a.statusUser)
	a.Put("/v1/user/password", a.passwordUser)

	a.Post("/v1/user/key", a.createUserKey)
	a.Delete("/v1/user/key", a.deleteUserKey)
	a.Get("/v1/user/key", a.getUserKeys)

	a.Get("/v1/users", a.getAllUsers)
	a.Get("/v1/users/findusers", a.findUsers)

	a.Get("/v1/device/login", a.loginDevice)
	a.Get("/v1/device/logout", a.logoutDevice)

	a.Get("/v1/session/status", a.status)
	a.Post("/v1/session/connect", a.connect)
	a.Post("/v1/session/disconnect", a.disconnect)

	//a.Post("/v1/webhook", a.createWebhook)
	//a.Get("/v1/webhook", a.deleteWebhook)
	//a.Get("/v1/webhook", a.getWebhooks)

	a.Post("/v1/chat/send/text", a.sendText)
	a.Post("/v1/chat/send/image", a.sendImage)
	a.Post("/v1/chat/send/document", a.sendDocument)

	a.Post("/v1/phone/scraping", a.scrapingPhones)
	a.Get("/v1/phone/contacts", a.getContacts)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(config *config.Config) error {
	if config.TLS.CertFile != "" && config.TLS.KeyFile != "" {
		log.Printf("Listening for TCP addr network: %v\n. Accepting connections over TLS.", config.SERVER.Host)
		if err := http.ListenAndServeTLS(
			config.SERVER.Host, config.TLS.CertFile, config.TLS.KeyFile, a.Router); err != nil {
			return err
		}
	} else {
		log.Printf("Listening for TCP addr network: %v\n", config.SERVER.Host)
		if err := http.ListenAndServe(config.SERVER.Host, a.Router); err != nil {
			return err
		}
	}

	return nil
}
