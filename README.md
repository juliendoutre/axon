# Axon

## Development

### Setup

```shell
make install
```

### Regenerate the protobuf code

```shell
make proto
```

### Add a SQL migration

```shell
make migrate MIGRATION_NAME=...
```

### Lint the code

```shell
make lint
```

### Run unit tests

```shell
make test
```

### Run locally

```shell
make dev
export TEST_JWT=$(CLIENT_ID=... CLIENT_SECRET=... goidc)
grpcurl localhost:8000 list
grpcurl localhost:8000 describe axon.api.v1.axon
grpcurl -H authorization:"Bearer ${TEST_JWT}" localhost:8000 axon.api.v1.axon/GetVersion
grpcurl -H authorization:"Bearer ${TEST_JWT}" -d '{"asset_type": 2, "asset_id": "google.com", "attributes": {"test": "a"}}' localhost:8000 axon.api.v1.axon/Observe
grpcurl -H authorization:"Bearer ${TEST_JWT}" -d '{"from": "2025-02-12T20:16:06Z", "to": "2025-02-12T23:20:06Z"}' localhost:8000 axon.api.v1.axon/CountObservations
grpcurl -H authorization:"Bearer ${TEST_JWT}" -d '{"from": "2025-02-12T20:16:06Z", "to": "2025-02-12T23:20:06Z", "page_size": 100, "page_token": 0}' localhost:8000 axon.api.v1.axon/ListObservations
make stop
```

## TODO

- search filters: https://github.com/grindlemire/go-lucene
- graph DB
- geoloc data
- jobs RPC
- documentation
- release
- myeline
