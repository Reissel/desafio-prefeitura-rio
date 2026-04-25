### BUILD
FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o prefeitura_app .

### EXECUTION
FROM golang:1.26.2-alpine3.23

RUN apk add --no-cache curl bash

RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY --from=builder app/prefeitura_app .
COPY justfile .

CMD ["./prefeitura_app"]