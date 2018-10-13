GO ?= GO111MODULE=on CGO_ENABLED=0 go
PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: all
all: install

.PHONY: clean
clean:
	$(GO) clean -i ./...

.PHONY: install
install:
	$(GO) install -v ./cmd/transmission-exporter

.PHONY: build
build:
	$(GO) build -v ./cmd/transmission-exporter

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;
