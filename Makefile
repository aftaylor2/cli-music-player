BINARY = player
CMD = ./cmd/player

.PHONY: build run test vet lint fmt clean

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
