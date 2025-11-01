# syntax=docker/dockerfile:1.7

############################
# Etapa de build ARM64
############################
FROM golang:1.25.3 AS build
WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct

# 1) Deps (capa estable + cache)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# 2) Codi
COPY . .

# 3) Build (cache de compilació)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o /out/pinger ./cmd/pinger

############################
# Etapa final (distroless)
############################
FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/pinger /pinger
USER nonroot
ENTRYPOINT ["/pinger"]
