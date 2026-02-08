# variable go at top 
.PHONY: dev down prod backend-test backend-build backend-clean

BACKEND_DIR = ./backend/

DOCKER_COMPOSE = docker compose --env-file $(BACKEND_DIR).env 

BINARY_NAME=main
LDFLAGS=-w -s

dev:
	$(DOCKER_COMPOSE) --profile dev up --watch --build

down:
	$(DOCKER_COMPOSE) down

prod:
	$(DOCKER_COMPOSE) --profile prod up

# multi-line will not work since each line runs in a separate shell
# to write multi-line either use && or "\"
backend-test:
	cd $(BACKEND_DIR) && go test -v ./...

backend-build:
	cd $(BACKEND_DIR) && CGO_ENABLED=0 GOOS=linux go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/main.go

# Clean build artifacts
backend-clean:
	cd $(BACKEND_DIR) && rm -f bin/$(BINARY_NAME)
