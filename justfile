
build:
  go build -o im ./cmd/im

install: build
  mkdir -p ~/.local/bin
  cp im ~/.local/bin/im

lint:
  golangci-lint run ./...

check: lint test

test:
  go test ./...

prod id:
  glow $(spec-driver show spec PROD-{{id}} --path) --pager

delta id:
  glow $(spec-driver show delta {{id}} --path) --pager
