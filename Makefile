APP_NAME := predictone
DB_URL := postgres://ChiaraGian:ChiaraGian@database:5432/DB_PredictOne?sslmode=disable

.PHONY: all run generate migrate apply status build test clean docker-up docker-down

all: build

run:
	@air

generate:
	@sqlc generate

migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	@atlas migrate diff "$(name)" \
		--dir "file://db/migrations" \
		--to "file://db/schema/schema.sql" \
		--dev-url "docker://postgres/18/dev?search_path=public"

apply:
	@atlas migrate apply \
		--dir "file://db/migrations" \
		--url "$(DB_URL)"


status:
	@atlas migrate status \
		--dir "file://db/migrations" \
		--url "$(DB_URL)"

build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

test:
	@go test ./...

clean:
	@rm -rf tmp

docker-up:
	@docker compose up --build

docker-down:
	@docker compose down

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
	$(MAKE) docker-up
	$(MAKE) wait
	$(MAKE) compile
	$(MAKE) test
	$(MAKE) docker-down
	$(MAKE) clean