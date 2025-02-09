# Axon

## Development

### Generate the protobuf code

```shell
brew install protobuf
protoc -Iprotos/v1 --go_out=. --go-grpc_out=. ./protos/v1/api.proto
```

### Add a SQL migration

```shell
brew install golang-migrate
migrate create -dir ./sql -ext .sql <MIGRATION_NAME>
```

### Lint the code

```shell
brew install golangci-lint hadolint sqlfluff
golangci-lint run
hadolint ./images/*.Dockerfile
sqlfluff lint --dialect postgres ./sql/*.sql
```

### Run unit tests

```shell
go test -v ./...
```

### Generate certs

```shell
brew install mkcert
mkcert -install
mkdir -p ./certs
mkcert -cert-file ./deploy/dev/server.crt.pem -key-file ./deploy/dev/server.key.pem server localhost
```

### Run locally

Start dependencies and the server:
```shell
export $(cat ./deploy/dev/.env | xargs)
docker compose -f ./deploy/dev/docker-compose.yaml up -d --build
POSTGRES_HOST=localhost go run ./cmd/migrator
POSTGRES_HOST=localhost go run ./cmd/server
```

In another shell, interact with the server:
```shell
export $(cat ./deploy/dev/.env | xargs)
grpcurl localhost:${SERVER_PORT} list
grpcurl localhost:${SERVER_PORT} describe axon.api.v1.axon
export TEST_JWT=$(curl -d "client_id=axon" -d "client_secret=ZU4ZKQ6H1sWT6SeB7uh3exUACThJ2Ma3" -d "username=admin" -d "password=${KEYCLOAK_ADMIN_PASSWORD}" -d "grant_type=password" "http://localhost:${KEYCLOAK_HTTP_PORT}/realms/master/protocol/openid-connect/token" | jq -c -r '.access_token')
grpcurl -H authorization:"Bearer ${TEST_JWT}" localhost:${SERVER_PORT} axon.api.v1.axon/GetVersion
grpcurl -H authorization:"Bearer ${TEST_JWT}" -d '{"asset_type": 2, "asset_id": "google.com", "attributes": {"test": "a"}}' localhost:${SERVER_PORT} axon.api.v1.axon/Observe
```

Stop dependencies with:
```shell
docker compose -f ./deploy/dev/docker-compose.yaml stop
docker compose -f ./deploy/dev/docker-compose.yaml rm
```
