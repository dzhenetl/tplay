PREFIX ?= /usr/local
DESTDIR ?=
VERSION ?= $(patsubst v%,%,$(shell git describe --tags --always --dirty 2>/dev/null || echo dev))
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build test lint install uninstall clean

build:
	go build -ldflags "$(LDFLAGS)" -o tplay ./cmd/tplay

test:
	go test -race ./...

lint:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi
	go vet ./...

install: build
	install -d "$(DESTDIR)$(PREFIX)/bin"
	install -m 755 tplay "$(DESTDIR)$(PREFIX)/bin/tplay"

uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/tplay"

clean:
	rm -f tplay
