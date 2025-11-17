# syntax=docker/dockerfile:1.20

############################
# Etapa de build ARM64
############################
FROM golang:1.25.4 AS builder
WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct
ARG TARGETARCH

# 1) Deps (capa estable + cache)
COPY go.mod go.sum ./
RUN go mod download

# 2) Codi
COPY . .

# 3) Build (cache de compilació)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /out/pinger ./cmd/pinger

# --- runtime ---
FROM alpine:3.22
# instal·la ping + setcap a la imatge final
RUN apk add --no-cache iputils libcap
# copia el teu binari
COPY --from=builder /out/pinger /usr/local/bin/pinger
# aplica capabilities al ping de la imatge final
RUN setcap cap_net_raw+ep "$(command -v ping)" && getcap "$(command -v ping)"

ENTRYPOINT ["/usr/local/bin/pinger"]
