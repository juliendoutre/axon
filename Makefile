include ./deploy/dev/.env

install:
	brew update
	brew install protobuf golang-migrate golangci-lint hadolint sqlfluff mkcert grpcurl tfenv jq

proto:
	protoc -Iprotos/v1 --go_out=. --go-grpc_out=. ./protos/v1/api.proto

migrate:
	migrate create -dir ./sql -ext .sql $(MIGRATION_NAME)

lint:
	golangci-lint run
	hadolint ./images/*.Dockerfile
	sqlfluff lint --dialect postgres ./sql/*.sql
	terraform fmt -recursive

test:
	go test -v ./...

certs:
	mkcert -install
	mkdir -p ./certs
	mkcert -cert-file ./deploy/dev/server.crt.pem -key-file ./deploy/dev/server.key.pem server localhost

dev: certs
	docker compose -f ./deploy/dev/docker-compose.yaml up -d --build
	sleep 10
	@cd ./deploy/dev && terraform init && TF_VAR_keycloak_http_port=${KEYCLOAK_HTTP_PORT} TF_VAR_keycloak_admin_password=${KEYCLOAK_ADMIN_PASSWORD} terraform apply -auto-approve

jwt:
	@cd ./deploy/dev && curl -s -d "client_id=axon" -d "client_secret=$$(terraform output -raw openid_client_secret)" -d "username=admin" -d "password=${KEYCLOAK_ADMIN_PASSWORD}" -d "grant_type=password" "http://localhost:${KEYCLOAK_HTTP_PORT}/realms/master/protocol/openid-connect/token" | jq -c -r '.access_token'

stop:
	docker compose -f ./deploy/dev/docker-compose.yaml stop
	docker compose -f ./deploy/dev/docker-compose.yaml rm --force
