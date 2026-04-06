# variable go at top
.PHONY: help dev down prod backend-test backend-build backend-clean

BACKEND_DIR = ./backend/

DOCKER_COMPOSE = docker compose --env-file $(BACKEND_DIR).env 

BINARY_NAME=main
LDFLAGS=-w -s

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  dev             Start dev environment (detached)"
	@echo "  dev-down        Stop dev environment"
	@echo "  prod            Start prod environment"
	@echo "  prod-down       Stop prod environment"
	@echo "  backend-test    Run backend tests"
	@echo "  backend-build   Build backend binary"
	@echo "  backend-clean   Remove backend build artifacts"

dev:
	# $(DOCKER_COMPOSE) --profile dev up --watch --build
	$(DOCKER_COMPOSE) --profile dev up --build -d

dev-down:
	$(DOCKER_COMPOSE) --profile dev down

prod-down:
	$(DOCKER_COMPOSE) --profile prod down

prod:
	$(DOCKER_COMPOSE) --profile prod up --build

# multi-line will not work since each line runs in a separate shell
# to write multi-line either use && or "\"
backend-test:
	cd $(BACKEND_DIR) && go test -v ./...

backend-build:
	cd $(BACKEND_DIR) && CGO_ENABLED=0 GOOS=linux go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/main.go

# Clean build artifacts
backend-clean:
	cd $(BACKEND_DIR) && rm -f bin/$(BINARY_NAME)
