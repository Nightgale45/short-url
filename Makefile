# variable go at top 

BACKEND_DIR = ./backend/


dev:
	docker compose --env-file $(BACKEND_DIR).env --profile dev up --watch

down:
	docker compose down

prod:
	docker compose --env-file $(BACKEND_DIR).env --profile prod

# multi-line will not work since each line runs in a separate shell
# to write multi-line either use && or \
backend-test:
	cd $(BACKEND_DIR) && go test -v
