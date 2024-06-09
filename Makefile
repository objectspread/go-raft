.DEFAULT_GOAL := fmt

SHELL := /bin/bash
GORAFT_IMPORT_PATH = github.com/objectspread/go-raft

# These DOCKER_xxx vars are used when building Docker images.
DOCKER_NAMESPACE?=objectspread
DOCKER_TAG?=latest

# All .go files that are not auto-generated and should be auto-formatted and linted.
ALL_SRC = $(shell find . -name '*.go' \
				   -not -name '_*' \
				   -not -name '.*' \
				   -not -name 'mocks*' \
				   -not -name '*.pb.go' \
				   -not -path './vendor/*' \
				   -not -path '*/mocks/*' \
				   -not -path '*/*-gen/*' \
				   -type f | \
				sort)

SED=sed
GO=go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
GOBUILD=CGO_ENABLED=0 installsuffix=cgo $(GO) build -trimpath
GOFMT=gofmt
GOFUMPT=gofumpt

FMT_LOG=.fmt.log
IMPORT_LOG=.import.log

GOVERSIONINFO=goversioninfo
SYSOFILE=resource.syso

UNAME := $(shell uname -m)
ifeq ($(UNAME), s390x)
# go test does not support -race flag on s390x architecture
	RACE=
else
	RACE=-race
endif

# GOBUILD=CGO_ENABLED=0 installsuffix=cgo $(GO) build -trimpath
GOTEST_QUIET=$(GO) test $(RACE)
GOTEST=$(GOTEST_QUIET) -v

COLORIZE ?= | $(SED) 's/PASS/✅ PASS/g' | $(SED) 's/FAIL/❌ FAIL/g' | $(SED) 's/SKIP/🔕 SKIP/g'

GIT_SHA=$(shell git rev-parse HEAD)
GIT_CLOSEST_TAG=$(shell git describe --abbrev=0 --tags)

ifneq ($(GIT_CLOSEST_TAG),$(shell echo ${GIT_CLOSEST_TAG} | grep -E "$(semver_regex)"))
	$(warning GIT_CLOSEST_TAG=$(GIT_CLOSEST_TAG) is not in the semver format $(semver_regex))
endif
GIT_CLOSEST_TAG_MAJOR := $(shell echo $(GIT_CLOSEST_TAG) | $(SED) -n 's/v\([0-9]*\)\.[0-9]*\.[0-9]/\1/p')

GIT_CLOSEST_TAG_MINOR := $(shell echo $(GIT_CLOSEST_TAG) | $(SED) -n 's/v[0-9]*\.\([0-9]*\)\.[0-9]/\1/p')
GIT_CLOSEST_TAG_PATCH := $(shell echo $(GIT_CLOSEST_TAG) | $(SED) -n 's/v[0-9]*\.[0-9]*\.\([0-9]\)/\1/p')
DATE=$(shell TZ=UTC0 git show --quiet --date='format-local:%Y-%m-%dT%H:%M:%SZ' --format="%cd")
BUILD_INFO_IMPORT_PATH=$(GORAFT_IMPORT_PATH)/pkg/version
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).commitSHA=$(GIT_SHA) -X $(BUILD_INFO_IMPORT_PATH).latestVersion=$(GIT_CLOSEST_TAG) -X $(BUILD_INFO_IMPORT_PATH).date=$(DATE)"

.PHONY: fmt
fmt:
	@echo Running gofmt on ALL_SRC ...
	@$(GOFMT) -e -s -l -w $(ALL_SRC)
	

.PHONY: lint
lint:
	golangci-lint -v run
	@[ ! -s "$(FMT_LOG)" -a ! -s "$(IMPORT_LOG)" ] || (echo "License check or import ordering failures, run 'make fmt'" | cat - $(FMT_LOG) $(IMPORT_LOG) && false)
	
.PHONY: clean
clean:
	rm -rf cover*.out .cover/ cover.html $(FMT_LOG) $(IMPORT_LOG) \
	GOCACHE=$(GOCACHE) go clean -cache -testcache
	find cmd -type f -executable | xargs -I{} sh -c '(git ls-files --error-unmatch {} 2>/dev/null || rm -v {})'

.PHONY: test
test:
	bash -c "set -e; set -o pipefail; $(GOTEST) -tags=memory_storage_integration ./... $(COLORIZE)"

# Requires variables: $(BIN_NAME) $(BIN_PATH) $(GO_TAGS) $(DISABLE_OPTIMIZATIONS) $(SUFFIX) $(GOOS) $(GOARCH) $(BUILD_INFO)
# Other targets can depend on this one but with a unique suffix to ensure it is always executed.
BIN_PATH = ./cmd/$(BIN_NAME)
.PHONY: _build-a-binary
_build-a-binary-%:
	$(GOBUILD) $(DISABLE_OPTIMIZATIONS) $(GO_TAGS) -o $(BIN_PATH)/$(BIN_NAME)$(SUFFIX)-$(GOOS)-$(GOARCH) $(BUILD_INFO) $(BIN_PATH)

.PHONY: build-client
build-client: BIN_NAME = client
build-client: _build-a-binary-client$(SUFFIX)-$(GOOS)-$(GOARCH)

.PHONY: build-server
build-server: BIN_NAME = server
build-server: _build-a-binary-server$(SUFFIX)-$(GOOS)-$(GOARCH)


# Magic values:
# - LangID "0409" is "US-English".
# - CharsetID "04B0" translates to decimal 1200 for "Unicode".
# - FileOS "040004" defines the Windows kernel "Windows NT".
# - FileType "01" is "Application".
define VERSIONINFO
{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": $(GIT_CLOSEST_TAG_MAJOR),
            "Minor": $(GIT_CLOSEST_TAG_MINOR),
            "Patch": $(GIT_CLOSEST_TAG_PATCH),
            "Build": 0
        },
        "ProductVersion": {
            "Major": $(GIT_CLOSEST_TAG_MAJOR),
            "Minor": $(GIT_CLOSEST_TAG_MINOR),
            "Patch": $(GIT_CLOSEST_TAG_PATCH),
            "Build": 0
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "FileDescription": "$(NAME)",
        "FileVersion": "$(GIT_CLOSEST_TAG_MAJOR).$(GIT_CLOSEST_TAG_MINOR).$(GIT_CLOSEST_TAG_PATCH).0",
        "LegalCopyright": "2015-2024 The Go-Raft Project Authors",
		"ProductName": "$(NAME)",
        "ProductVersion": "$(GIT_CLOSEST_TAG_MAJOR).$(GIT_CLOSEST_TAG_MINOR).$(GIT_CLOSEST_TAG_PATCH).0"
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0409",
            "CharsetID": "04B0"
        }
    }
}
endef

export VERSIONINFO

.PHONY: _prepare-winres
_prepare-winres:
	$(MAKE) _prepare-winres-helper NAME="Go-Raft Client"            PKGPATH="cmd/client"
	$(MAKE) _prepare-winres-helper NAME="Go-Raft Server"        PKGPATH="cmd/server"

.PHONY: _prepare-winres-helper
_prepare-winres-helper:
	echo $$VERSIONINFO | $(GOVERSIONINFO) -o="$(PKGPATH)/$(SYSOFILE)" -

.PHONY: build-binaries-linux
build-binaries-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) _build-platform-binaries

.PHONY: build-binaries-windows
build-binaries-windows: _prepare-winres
	GOOS=windows GOARCH=amd64 $(MAKE) _build-platform-binaries
	rm ./cmd/*/$(SYSOFILE)

.PHONY: build-binaries-darwin
build-binaries-darwin:
	GOOS=darwin GOARCH=amd64 $(MAKE) _build-platform-binaries

.PHONY: build-binaries-darwin-arm64
build-binaries-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) _build-platform-binaries

.PHONY: build-binaries-s390x
build-binaries-s390x:
	GOOS=linux GOARCH=s390x $(MAKE) _build-platform-binaries

.PHONY: build-binaries-arm64
build-binaries-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) _build-platform-binaries

.PHONY: build-binaries-ppc64le
build-binaries-ppc64le:
	GOOS=linux GOARCH=ppc64le $(MAKE) _build-platform-binaries

# build all binaries for one specific platform GOOS/GOARCH
.PHONY: _build-platform-binaries
_build-platform-binaries: \
	build-client \
	build-server

.PHONY: build-all-platforms
build-all-platforms: \
	build-binaries-linux \
	build-binaries-darwin \
	build-binaries-darwin-arm64 \
	build-binaries-s390x \
	build-binaries-arm64 \
	build-binaries-ppc64le \
	build-binaries-windows



.PHONY: install-test-tools
install-test-tools:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	$(GO) install mvdan.cc/gofumpt@latest

.PHONY: install-build-tools
install-build-tools:
	$(GO) install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.4.0

.PHONY: install-tools
install-tools: install-test-tools install-build-tools
	$(GO) install github.com/vektra/mockery/v2@v2.14.0

.PHONY: install-ci
install-ci: install-test-tools install-build-tools
