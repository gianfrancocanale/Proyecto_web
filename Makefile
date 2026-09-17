.PHONY: test generate clean start wait compile

generate:
	@sqlc generate

clean:
	@echo "Limpiando contenedores y volúmenes"
	@docker compose down -v --remove-orphans

start:
	@echo "Iniciando contenedor de PostgreSQL"
	@docker compose up -d database

wait:
	@echo "Esperando a que PostgreSQL esté listo"
	@until docker exec postgres-db pg_isready -U ChiaraGian -d DB_PredictOne > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "PostgreSQL está listo."

compile:
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
	$(MAKE) compile
	@echo "Ejecutando API en Docker"
	docker compose up --build api
	$(MAKE) clean
