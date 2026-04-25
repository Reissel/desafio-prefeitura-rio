default: run

run:
    ./prefeitura_app

build:
    @go build -o target/main main.go

clean:
    rm -rf target/

test:
    docker-compose -f docker-compose.test.yml up   --build   --abort-on-container-exit   --exit-code-from app-tester && docker-compose -f docker-compose.test.yml down

integration-test:
    go test -v ./tests/...
    