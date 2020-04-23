GO ?= GO111MODULE=on CGO_ENABLED=0 go
PACKAGES = $(shell go list ./... | grep -v /vendor/)
BUILDX ?= docker buildx
PLATFORMS ?= linux/amd64,linux/i386,linux/arm64,linux/arm/v7
DOCKER_IMAGE ?= kennedyoliveira/transmission-exporter

ifeq ($(DOCKER_TAG),)
	DOCKER_TAG = latest
endif

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
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: dashboards
dashboards:
	jsonnet fmt -i dashboards/transmission.jsonnet
	jsonnet -J dashboards/vendor -m dashboards -e "(import 'dashboards/transmission.jsonnet').grafanaDashboards"


.PHONE: docker-init
docker-init:
	@$(BUILDX) create --name transmission-exporter-builder
	@$(BUILDX) use transmission-exporter-builder
	@$(BUILDX) inspect --bootstrap transmission-exporter-builder

.PHONE: docker-build
docker-build: clean
	@echo ">> building multi-arch docker images, tag=$(DOCKER_TAG)"
	@$(BUILDX) build -f Dockerfile-cross \
			  --platform $(PLATFORMS) \
			  --tag $(DOCKER_IMAGE):$(DOCKER_TAG) \
			  --push \
			  .

.PHONE: docker-clean
docker-clean:
	@$(BUILDX) rm transmission-exporter-builder