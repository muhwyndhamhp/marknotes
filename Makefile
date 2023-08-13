ifeq ("$(wildcard .env)","")
    $(shell cp env.sample .env)
endif

include .env

$(eval export $(grep -v '^#' .env | xargs -0))

GO_MODULE := github.com/muhwyndhamhp/todo-mx
VERSION  ?= $(shell git describe --tags --abbrev=0)
LDFLAGS   := -X "$(GO_MODULE)/config.Version=$(VERSION)"
DB_DSN    := "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

tools: $(MIGRATE) $(AIR) $(MOCKERY) $(GOLANGCI) $(CHGLOG)
tools:
	@echo "Required tools are installed"

setup-local: tools
	@docker-compose up -d
	@sleep 5
	@docker exec -it material-postgres psql -h localhost -p 5432 -U $(DB_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || (docker exec -it material-postgres psql -h localhost -p 5432 -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME)" && echo "Database $(DB_NAME) created")

run:
	@air -c .air.toml --build.cmd "go build -ldflags \"$(LDFLAGS)\" -o ./tmp/main ."
