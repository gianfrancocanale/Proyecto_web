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

wait-api:
	@echo "Esperando a que la API responda"
	@for i in $$(seq 1 60); do \
		curl -s -o /dev/null $(API_URL)/usuarios && exit 0; \
		sleep 1; \
	done; \
	echo "La API no respondió a tiempo"; docker compose logs api; exit 1

compile:
	@echo "Compilando"
	@go build ./...
	@go test -run '^$$' ./...

test:
	$(MAKE) clean
	$(MAKE) generate
	$(MAKE) start
	$(MAKE) wait
	$(MAKE) compile
	@echo "Ejecutando Tests"
	go test ./...
	$(MAKE) clean

run:
	$(MAKE) clean
	$(MAKE) generate
	$(MAKE) start
	$(MAKE) wait
	@templ generate
	@echo "Ejecutando API en Docker"
	docker compose up --build api
	$(MAKE) clean

# Todo el flujo de la entrega: preparar, levantar, probar endpoints y limpiar.
# El trap hace que se limpie siempre, incluso si algo falla, y conserva el código de salida.
test-api:
	@trap '$(MAKE) clean' EXIT; \
	$(MAKE) clean && \
	$(MAKE) generate && \
	$(MAKE) compile && \
	$(MAKE) start && \
	$(MAKE) wait && \
	docker compose up -d --build api && \
	$(MAKE) wait-api && \
	$(MAKE) hurl

hurl:
	@echo "Ejecutando todas las pruebas Hurl"
	hurl --test requests/usuarios.hurl
	hurl --test requests/apuestas.hurl
	hurl --test requests/escuderias.hurl
	hurl --test requests/pilotos.hurl
	hurl --test requests/pilotohistoricos.hurl
	hurl --test requests/granpremio.hurl
	hurl --test requests/granpremiohistorico.hurl

hurlUsuarios:
	@echo "Ejecutando pruebas Hurl para usuarios"
	hurl --test requests/usuarios.hurl

hurlApuestas:
	@echo "Ejecutando pruebas Hurl para apuestas"
	hurl --test requests/apuestas.hurl

hurlEscuderias:
	@echo "Ejecutando pruebas Hurl para escuderías"
	hurl --test requests/escuderias.hurl

hurlPilotos:
	@echo "Ejecutando pruebas Hurl para pilotos"
	hurl --test requests/pilotos.hurl

hurlPilotosHistoricos:
	@echo "Ejecutando pruebas Hurl para pilotos históricos"
	hurl --test requests/pilotohistoricos.hurl

hurlGrandesPremios:
	@echo "Ejecutando pruebas Hurl para grandes premios"
	hurl --test requests/granpremio.hurl

hurlGrandesPremiosHistoricos:
	@echo "Ejecutando pruebas Hurl para grandes premios históricos"
	hurl --test requests/granpremiohistorico.hurl