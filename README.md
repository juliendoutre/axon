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

```shell
docker compose -f ./deploy/dev/docker-compose.yaml up -d --build
grpcurl localhost:8000 list
grpcurl localhost:8000 describe axon.api.v1.axon
export TEST_JWT=$(curl -d "client_id=axon" -d "client_secret=ZU4ZKQ6H1sWT6SeB7uh3exUACThJ2Ma3" -d "username=admin" -d "password=${KEYCLOAK_ADMIN_PASSWORD}" -d "grant_type=password" "http://localhost:7080/realms/master/protocol/openid-connect/token" | jq -c -r '.access_token')
grpcurl -H authorization:"Bearer ${TEST_JWT}" localhost:8000 axon.api.v1.axon/GetVersion
grpcurl -H authorization:"Bearer ${TEST_JWT}" -d '{"asset_type": 2, "asset_id": "google.com", "attributes": {"test": "a"}}' localhost:8000 axon.api.v1.axon/Observe
docker compose -f ./deploy/dev/docker-compose.yaml stop
docker compose -f ./deploy/dev/docker-compose.yaml rm
```
