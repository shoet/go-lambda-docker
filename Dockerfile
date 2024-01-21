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
    # go build -trimpath -ldflags="-w -s" -o cmd/bin/main cmd/main.go
    go build -trimpath -ldflags="-w -s" -o cmd/bin/main function/main.go

# ===== deploy stage ====
# FROM golang:1.20.12-bullseye as deploy
FROM mcr.microsoft.com/playwright:v1.40.0-jammy as deploy

RUN apt update
RUN apt install -y golang-1.20
ENV GOPATH=/go
ENV GOROOT=/usr/lib/go-1.20
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

WORKDIR /app

# ENV PLAYWRIGHT_BROWSERS_PATH=/tmp/playwright/browser
RUN mkdir -p /var/playwright/browser
ENV PLAYWRIGHT_BROWSERS_PATH=/var/playwright/browser
RUN go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps

COPY --from=builder /app/cmd/bin/main .

CMD ["/app/main"]
