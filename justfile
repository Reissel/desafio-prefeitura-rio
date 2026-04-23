default: run

run:
    @go run main.go

build:
    @go build -o target/main main.go

clean:
    rm -rf target/