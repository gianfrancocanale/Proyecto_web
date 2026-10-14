.PHONY: test generate clean start wait compile run hurl hurlUsuarios hurlApuestas hurlEscuderias hurlPilotos hurlPilotosHistoricos hurlGrandesPremios hurlGrandesPremiosHistoricos

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
	@echo "Ejecutando API en Docker"
	docker compose up --build api
	$(MAKE) clean

hurl:
	@echo "Ejecutando todas las pruebas Hurl"
	hurl --test requests/usuarios.hurl
	hurl --test requests/apuestas.hurl
	hurl --test requests/escuderias.hurl
	hurl --test requests/pilotos.hurl
	hurl --test requests/pilotos_historicos.hurl
	hurl --test requests/gran_premios.hurl
	hurl --test requests/gran_premios_historicos.hurl

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
	hurl --test requests/pilotos_historicos.hurl

hurlGrandesPremios:
	@echo "Ejecutando pruebas Hurl para grandes premios"
	hurl --test requests/gran_premios.hurl

hurlGrandesPremiosHistoricos:
	@echo "Ejecutando pruebas Hurl para grandes premios históricos"
	hurl --test requests/gran_premios_historicos.hurl
