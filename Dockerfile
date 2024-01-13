# ===== build stage ====
FROM golang:1.21.5-bullseye as builder

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -trimpath -ldflags="-w -s" -o cmd/bin/main cmd/main.go

# ===== deploy stage ====
FROM golang:1.21.5-bullseye as deploy

WORKDIR /app

RUN apt update -y

COPY --from=builder /app/cmd/bin/main .

CMD ["/app/main"]
