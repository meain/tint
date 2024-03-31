.PHONY: run build test clean

BINARY_NAME=tint

build:
	go build -o $(BINARY_NAME) -v

run:
	go run ./... .

test: build # TODO: add proper go test
	./tint lint --config .tint.toml.sample .

clean:
	go clean
	rm -f $(BINARY_NAME)