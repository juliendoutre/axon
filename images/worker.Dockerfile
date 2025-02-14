# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
ARG ALPINE_VERSION=3.19
ARG DEBIAN_VERSION=12

FROM --platform=$BUILDPLATFORM index.docker.io/golang:$GO_VERSION-alpine$ALPINE_VERSION AS builder

ARG TARGETOS
ARG TARGETARCH
ARG GO_VERSION

WORKDIR /axon/worker

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/worker ./cmd/worker
COPY ./pkg ./pkg
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-s -w -X main.GoVersion=$GO_VERSION -X main.Os=$TARGETOS -X main.Arch=$TARGETARCH" -o /worker ./cmd/worker

FROM --platform=$TARGETPLATFORM gcr.io/distroless/base-debian$DEBIAN_VERSION:latest AS runner

LABEL org.opencontainers.image.authors Julien Doutre <jul.doutre@gmail.com>
LABEL org.opencontainers.image.title axon.worker
LABEL org.opencontainers.image.url https://github.com/juliendoutre/axon
LABEL org.opencontainers.image.documentation https://github.com/juliendoutre/axon
LABEL org.opencontainers.image.source https://github.com/juliendoutre/axon/tree/${GIT_COMMIT_SHA}/images/worker.Dockerfile
LABEL org.opencontainers.image.licenses MIT
LABEL org.opencontainers.revision ${GIT_COMMIT_SHA}

WORKDIR /

COPY --from=builder /worker /worker

USER nonroot:nonroot

ENTRYPOINT ["/worker"]
