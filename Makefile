PROJECT := genda-api
.DEFAULT_GOAL := docker-up

docker-build:
	@docker build \
		-t genda-api \
		--build-arg PACKAGE_NAME=genda-api \
		--build-arg BUILD_REF=`git rev-parse --short HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

docker-push:
	@docker tag genda-api genda/genda-api:latest
	@docker push genda/genda-api:latest

docker-push-dev:
	@docker tag genda-api genda/genda-api:dev
	@docker push genda/genda-api:dev

genda-api: 
	@go run ./cmd/api/main.go

docker-up:
	@docker-compose up --build --remove-orphans

docker-down:
	@docker-compose down

build:
	@go build -o ./tmp/cmd/api cmd/api/main.go

test:
	@go test ./... -count=1

clean:
	@docker system prune -f

docker-stop-all:
	@docker stop $(docker ps -aq)

docker-remove-all:
	@docker rm $(docker ps -aq)

deps-reset:
	@git checkout -- go.mod
	@go mod tidy
	@go mod vendor

tidy:
	@go mod tidy
	@go mod vendor

deps-upgrade:
	@go get -u -t -d -v ./...
	@go mod vendor

deps-cleancache:
	@go clean -modcache

migrate-up:
	@migrate -path db/migration -database "postgresql://postgres:postgres@postgres:5432/genda_api?sslmode=disable" -verbose up

migrate-down:
	@migrate -path db/migration -database "postgresql://postgres:postgres@postgres:5432/genda_api?sslmode=disable" -verbose down

migrate-create:
	@migrate create -ext sql -dir db/migration -seq 