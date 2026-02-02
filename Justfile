test:
    go vet ./...
    go test ./...

lint:
    golangci-lint run

build:
    go build

repl: build
    ./kanren

debug-repl: build
    ./kanren -debug -pretty

fmt:
    go fmt ./...

clean:
    rm -rf *.out *.html ./kanren
