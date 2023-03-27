package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB      DataBase
	SERVER  Server
	TLS     Certificate
	AUTH    Authentication
	LICENSE License
}

type DataBase struct {
	Host     string `name:"dbhost" description:"Set the http server host." default:"0.0.0.0"`
	Port     string `name:"dbport" description:"Set the http server port." default:"8080"`
	User     string `name:"dbuser" description:"Set the database user."`
	Password string `name:"dbpassword" description:"Set the database password."`
	Dialect  string `name:"dbdialect" description:"Set the database dialect." default:"sqlite"`
	DbURI    string `name:"dburi" description:"Set the database dialect." default:"file:dbdata/wppserver.db?_foreign_keys=on"`
}

type Server struct {
	Host string
}

type Authentication struct {
	UserEmail    string `name:"useremail" description:"Set initial admin email." default:"admin@wppserver.com"`
	UserPassword string `name:"userpassword" description:"Set initial admin password." default:"admin"`
	SecretKey    string `name:"username" description:"Set your secret key for encryption. Use a random phrase."`
}

type Certificate struct {
	CertFile string `name:"certfile" description:"The certificate file path."`
	KeyFile  string `name:"keyfile" description:"The certificate key file path."`
}

type License struct {
	Key string `name:"username" description:"Set your purchase key. Ex. XXXX-XXXX-XXXXX-XXXX"`
}

func GetConfig(input ...string) *Config {
	if len(input) != 0 && input[0] != "" {
		err := godotenv.Load(input[0])
		if err != nil {
			log.Print("The .env file in path not found : ", input[0], err)
		}
	}

	return &Config{
		DB: DataBase{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Dialect:  os.Getenv("DATABASE_TYPE"),
			DbURI:    os.Getenv("DATABASE_URI"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
		},
		SERVER: Server{
			Host: os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("SERVER_PORT"),
		},
		TLS: Certificate{
			CertFile: os.Getenv("CERTIFICATE_FILE"),
			KeyFile:  os.Getenv("CERTIFICATE_KEY_FILE"),
		},
		AUTH: Authentication{
			UserEmail:    os.Getenv("AUTH_INITIAL_EMAIL"),
			UserPassword: os.Getenv("AUTH_INITIAL_PASSWORD"),
			SecretKey:    os.Getenv("AUTH_JWT_SECRET"),
		},
		LICENSE: License{
			Key: os.Getenv("LICENSE_KEY"),
		},
	}
}
