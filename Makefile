ifeq (,$(wildcard .env))
$(shell cp .env.example .env)
endif
include .env
export

API_URL ?= http://localhost:8080

.PHONY: generate clean start wait wait-api compile test run test-api hurl hurlUsuarios hurlApuestas hurlEscuderias hurlPilotos hurlPilotosHistoricos hurlGrandesPremios hurlGrandesPremiosHistoricos

generate:
	@sqlc generate

clean:
	@echo "Limpiando contenedores y volúmenes"
	@docker compose down -v --remove-orphans

start:
	@echo "Iniciando contenedor de PostgreSQL"
	@docker compose up -d database

# -h localhost fuerza conexión TCP: durante el init de la imagen postgres el
# servidor temporal solo escucha por socket, así esperamos al servidor real.
wait:
	@echo "Esperando a que PostgreSQL esté listo"
	@until docker exec postgres-db pg_isready -h localhost -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "PostgreSQL está listo."


compile:
	@echo "Compilando"
	@go build ./...
	@go test -run '^$$' ./...


run:
	$(MAKE) clean
	$(MAKE) generate
	$(MAKE) start
	$(MAKE) wait
	@templ generate
	@echo "Ejecutando API en Docker"
	docker compose up --build api
	$(MAKE) clean

