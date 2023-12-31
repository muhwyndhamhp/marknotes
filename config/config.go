package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_FILE = ".env"

	APP_PORT              = "APP_PORT"
	JWT_SECRET            = "JWT_SECRET"
	OAUTH_AUTHORIZE_URL   = "OAUTH_AUTHORIZE_URL"
	OAUTH_ACCESSTOKEN_URL = "OAUTH_ACCESSTOKEN_URL"
	OAUTH_CLIENTID        = "OAUTH_CLIENTID"
	OAUTH_SECRET          = "OAUTH_SECRET"
	OAUTH_URL             = "OAUTH_URL"
	RESUME_POST_ID        = "RESUME_POST_ID"
	DATABASE_URL          = "DATABASE_URL"
)

func init() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("Failed to load env file: %s\n", err)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
