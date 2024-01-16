# ===== build stage ====
FROM golang:1.20.12-bullseye as builder

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
    # go build -trimpath -ldflags="-w -s" -o cmd/bin/main lambda/main.go

# ===== deploy stage ====
FROM golang:1.20.12-alpine3.18 as deploy

WORKDIR /app

RUN apk update

RUN apk add chromium
RUN apk add libc6-compat
RUN apk add gcompat
RUN ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2

COPY --from=builder /app/cmd/bin/main .

CMD ["/app/main"]
