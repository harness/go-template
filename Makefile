MAIN_PACKAGE_PATH	:= $(CURDIR)/cmd/go-template
BINARY_NAME			:= go-template
BINDIR				:= $(CURDIR)/bin
DISTDIR				:= $(CURDIR)/dist
PARALLEL_BUILDS_CNT	:= 4
TARGETS				:= darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
LDFLAGS				:= -w -s
GOFLAGS				:=

# version
TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
GIT_REF := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
ifeq ($(VERSION),)
    VERSION := 0.0.1-dev
endif
ifneq ($(GIT_REF), $(TAG_COMMIT))
    VERSION := $(VERSION)-next-$(GIT_REF)-$(DATE)
endif
ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif

ifdef VERSION
	VERSION_MAJOR = $(word 1, $(subst ., ,$(VERSION)))
	VERSION_MINOR = $(word 2, $(subst ., ,$(VERSION)))
	VERSION_PATCH = $(word 3, $(subst ., ,$(VERSION)))
endif

LDFLAGS += -X main.versionMajor=$(VERSION_MAJOR)
LDFLAGS += -X main.versionMinor=$(VERSION_MINOR)
LDFLAGS += -X main.versionPatch=$(VERSION_PATCH)
LDFLAGS += -X main.gitref=$(GIT_REF)

# dependencies tools
GOBIN	= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN 	= $(shell go env GOPATH)/bin
endif
GOX 	= $(GOBIN)/gox
GOLANGCI_LINT = $(GOBIN)/golangci-lint

$(GOLANGCI_LINT):
	@echo "Instal golangci-lint"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.55.1

$(GOX):
	@echo "Install gox"
	(cd /; go install github.com/mitchellh/gox@latest)

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## lint: code analysis
.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run \
				-E gofmt \
				-E govet \
				-E stylecheck \
				-E gosec \
				-E asciicheck \
				-E goimports \
				-E exportloopref \
				-E unparam 
	@echo "Everything good"

.PHONY: test
test:
	go test -v -race -buildvcs ./...

## build: build binary for current platform
.PHONY: build
build:
	CGO_ENABLED=0 go build -o=${BINDIR}/${BINARY_NAME} -ldflags '$(LDFLAGS)' ${MAIN_PACKAGE_PATH}

.PHONY: build-release
build-release: $(GOX)
build-release: LDFLAGS += -extldflags "-static"
build-release: 
	$(if $(VERSION),, $(error Version is required for build-release target))
	$(GOX) -parallel=$(PARALLEL_BUILDS_CNT) -output="${DISTDIR}/$(BINARY_NAME)/release/v$(VERSION)/bin/{{.OS}}/{{.Arch}}/$(BINARY_NAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ${MAIN_PACKAGE_PATH}