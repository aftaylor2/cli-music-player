BINARY = player
CMD = ./cmd/player
PREFIX ?= /usr/local

.PHONY: build run test vet lint fmt clean install uninstall

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

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 $(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)
