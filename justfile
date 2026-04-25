default: run

run:
    docker compose up -d

stop:
    docker compose down

build:
    @go build -o target/main main.go

clean:
    rm -rf target/

test:
    go test -v ./logic/... && docker-compose -f docker-compose.test.yml up   --build   --abort-on-container-exit   --exit-code-from app-tester; docker-compose -f docker-compose.test.yml down
    