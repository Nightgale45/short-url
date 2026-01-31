# variable go at top 

BACKEND_DIR = ./backend/

DOCKER_COMPOSE = docker compose --env-file $(BACKEND_DIR).env 


dev:
	$(DOCKER_COMPOSE) --profile dev up --watch

down:
	$(DOCKER_COMPOSE) down

prod:
	$(DOCKER_COMPOSE) --profile prod

# multi-line will not work since each line runs in a separate shell
# to write multi-line either use && or \
backend-test:
	cd $(BACKEND_DIR) && go test -v
