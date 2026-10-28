BINARY := pomo
VERSION ?= 0.1.0
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build install test vet fmt clean perf

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/pomo

install:
	go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/pomo

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

clean:
	rm -f $(BINARY)

# Startup time and binary size. RSS is checked manually while running.
perf: build
	@ls -lh $(BINARY)
	@/usr/bin/time -f 'startup: %e s' ./$(BINARY) --version >/dev/null
