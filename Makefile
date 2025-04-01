install:
	brew update
	brew install protobuf golang-migrate golangci-lint hadolint sqlfluff mkcert grpcurl jq

proto:
	protoc -Iprotos/v1 --go_out=. --go-grpc_out=. ./protos/v1/api.proto

migrate:
	migrate create -dir ./sql -ext .sql $(MIGRATION_NAME)

lint:
	golangci-lint run
	hadolint ./images/*.Dockerfile
	sqlfluff lint --dialect postgres ./sql/*.sql

test:
	go test -v ./...

certs:
	mkcert -install
	mkcert -cert-file ./deploy/dev/server.crt.pem -key-file ./deploy/dev/server.key.pem server localhost

dev: certs
	docker compose -f ./deploy/dev/docker-compose.yaml up -d --build

stop:
	docker compose -f ./deploy/dev/docker-compose.yaml stop
	docker compose -f ./deploy/dev/docker-compose.yaml rm --force
