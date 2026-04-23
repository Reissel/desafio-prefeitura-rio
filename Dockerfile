FROM golang:1.26.2-alpine3.23

RUN apk add --no-cache curl bash

RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY . .

CMD ["just", "run"]