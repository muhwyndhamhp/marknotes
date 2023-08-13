package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhwyndhamhp/gotes-mx/utils/errs"
)

const (
	ENV_FILE = ".env"

	APP_PORT    = "APP_PORT"
	DB_HOST     = "DB_HOST"
	DB_PORT     = "DB_PORT"
	DB_USER     = "DB_USER"
	DB_NAME     = "DB_NAME"
	DB_PASSWORD = "DB_PASSWORD"
)

func init() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		log.Fatal(errs.Wrap(err))
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
