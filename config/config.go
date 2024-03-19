package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_FILE = ".env"

	ENV                   = "ENV"
	APP_PORT              = "APP_PORT"
	JWT_SECRET            = "JWT_SECRET"
	OAUTH_AUTHORIZE_URL   = "OAUTH_AUTHORIZE_URL"
	OAUTH_ACCESSTOKEN_URL = "OAUTH_ACCESSTOKEN_URL"
	OAUTH_CLIENTID        = "OAUTH_CLIENTID"
	OAUTH_SECRET          = "OAUTH_SECRET"
	OAUTH_URL             = "OAUTH_URL"
	RESUME_POST_ID        = "RESUME_POST_ID"
	DATABASE_URL          = "DATABASE_URL"
	STORE_VOL_PATH        = "STORE_VOL_PATH"

	LIBSQL_URL        = "LIBSQL_URL"
	LIBSQL_AUTH_TOKEN = "LIBSQL_AUTH_TOKEN"

	CF_ACCOUNT_ID           = "CF_ACCOUNT_ID"
	CF_R2_ACCESS_KEY_ID     = "CF_R2_ACCESS_KEY_ID"
	CF_R2_SECRET_ACCESS_KEY = "CF_R2_SECRET_ACCESS_KEY"
	CLERK_SECRET            = "CLERK_SECRET"
	CLERK_PUBLISHABLE       = "CLERK_PUBLISHABLE"
	CLERK_SRC_URL           = "CLERK_SRC_URL"
)

func init() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("Failed to load env file: %s\n", err)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
