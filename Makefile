.PHONY: verify go-test frontend-check frontend-build studio-dev

GOCACHE ?= $(CURDIR)/.cache/go-build
GOPATH ?= $(CURDIR)/.cache/go
WAILS_BIN ?= $(CURDIR)/bin/wails

verify: go-test frontend-check frontend-build

go-test:
	GOCACHE=$(GOCACHE) GOPATH=$(GOPATH) go test ./...

frontend-check:
	cd apps/studio/frontend && npm run check

frontend-build:
	cd apps/studio/frontend && npm run build

studio-dev:
	cd apps/studio && if [ -x "$(WAILS_BIN)" ]; then "$(WAILS_BIN)" dev; else wails dev; fi
