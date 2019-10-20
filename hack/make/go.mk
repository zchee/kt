# ----------------------------------------------------------------------------
# global

SHELL = /usr/bin/env bash

# hack for replace all whitespace to comma
comma := ,
empty :=
space := $(empty) $(empty)

# ----------------------------------------------------------------------------
# Go

ifneq ($(shell command -v go),)
GO_PATH ?= $(shell go env GOPATH)
GO_OS ?= $(shell go env GOOS)
GO_ARCH ?= $(shell go env GOARCH)

export GO111MODULE=on
GO_MOD_FLAGS =
ifneq ($(wildcard go.mod),)  # exist go.mod
ifneq ($(wildcard ./vendor),)  # exist vender directory
	GO_MOD_FLAGS=-mod=vendor
endif
endif

PKG := $(subst $(GO_PATH)/src/,,$(CURDIR))
GO_PKGS := $(shell go list ${GO_MOD_FLAGS} ./... | grep -v -e '.pb.go')
GO_APP_PKGS := $(shell go list ${GO_MOD_FLAGS} -f '{{if and (or .GoFiles .CgoFiles) (ne .Name "main")}}{{.ImportPath}}{{end}}' ${PKG}/...)
GO_TEST_PKGS := $(shell go list ${GO_MOD_FLAGS} -f='{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)
GO_VENDOR_PKGS=
ifneq ($(wildcard ./vendor),)  # exist vender directory
GO_VENDOR_PKGS = $(shell go list ${GO_MOD_FLAGS} -f '{{if and (or .GoFiles .CgoFiles) (ne .Name "main")}}./vendor/{{.ImportPath}}{{end}}' ./vendor/...)
endif
endif

CGO_ENABLED ?= 0
GO_GCFLAGS=
GO_CHECKPTR_FLAGS=all=-d=checkptr=1 -d=checkptr=2
# https://tip.golang.org/doc/diagnostics.html#debugging
GO_GCFLAGS_DEBUG=all=-N -l -dwarflocationlists=true
GO_LDFLAGS=-s -w
GO_LDFLAGS_STATIC="-extldflags=-fno-PIC -static"
GO_LDFLAGS_DEBUG=-compressdwarf=false

GO_BUILDTAGS=
ifeq (${CGO_ENABLED},0)
	GO_BUILDTAGS=osusergo netgo
endif
GO_BUILDTAGS_STATIC=static static_build
GO_INSTALLSUFFIX_STATIC=-installsuffix 'netgo'
GO_FLAGS = -tags='$(subst $(space),$(comma),${GO_BUILDTAGS})'

GO_FLAGS+=-gcflags='${GO_GCFLAGS}'
GO_FLAGS+=-ldflags='${GO_LDFLAGS}'

GO_TEST ?= go test
ifneq ($(shell command -v gotest),)
	GO_TEST=gotest
endif
GO_TEST_FUNC ?= .
GO_TEST_FLAGS ?=
GO_BENCH_FUNC ?= .
GO_BENCH_FLAGS ?= -benchmem
ifneq ($(wildcard go.mod),)  # exist go.mod
ifneq ($(GO111MODULE),off)
	GO_TEST_FLAGS+=${GO_MOD_FLAGS}
	GO_BENCH_FLAGS+=${GO_MOD_FLAGS}
endif
endif
GO_TEST_COVERAGE_OUT := coverage.out
ifneq ($(CIRCLECI),)
	GO_TEST_COVERAGE_OUT=/tmp/ci/artifacts/coverage.out
endif

CONTAINER_REGISTRY := gcr.io/containerz
CONTAINER_BUILD_TAG ?= $(VERSION:v%=%)
CONTAINER_BUILD_ARGS_BASE = --rm --pull --build-arg GOLANG_VERSION=${GOLANG_VERSION} --build-arg ALPINE_VERSION=${ALPINE_VERSION}
ifneq (${SHORT_SHA},)  ## for Cloud Build
	CONTAINER_BUILD_ARGS_BASE+=--build-arg SHORT_SHA=${SHORT_SHA}
endif
CONTINUOUS_INTEGRATION ?=  ## for buildkit
ifneq (${CONTINUOUS_INTEGRATION},)
	CONTAINER_BUILD_ARGS_BASE+=--progress=plain
endif
CONTAINER_BUILD_ARGS ?= ${CONTAINER_BUILD_ARGS_BASE}
CONTAINER_BUILD_TARGET ?= ${APP}
ifneq (${CONTAINER_BUILD_TARGET},${APP})
	CONTAINER_BUILD_TAG=${VERSION}-${CONTAINER_BUILD_TARGET}
endif

# ----------------------------------------------------------------------------
# defines

GOPHER = ""
define target
@printf "$(GOPHER)  \\x1b[1;32m$(patsubst ,$@,$(1))\\x1b[0m\\n"
endef

# ----------------------------------------------------------------------------
# targets

.PHONY: all
all: mod pkg/install static

##@ build and install

.PHONY: $(APP)
$(APP):
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -v $(strip $(GO_FLAGS)) -o $(APP) $(CMD)

.PHONY: build
build: GO_FLAGS+=${GO_MOD_FLAGS}
build: $(APP)  ## Builds a dynamic executable.

.PHONY: static
static: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
static: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
static: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
static: GO_FLAGS+=${GO_MOD_FLAGS}
static: $(APP)  ## Builds a static executable.

.PHONY: install
install: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
install: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
install: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
install: GO_FLAGS+=${GO_MOD_FLAGS}
install:  ## Installs the executable.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go install -v $(strip $(GO_FLAGS)) $(CMD)

.PHONY: pkg/install
pkg/install: GO_FLAGS+=${GO_MOD_FLAGS}
pkg/install: GO_LDFLAGS=
pkg/install: GO_BUILDTAGS=
pkg/install:
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go install -v ${GO_MOD_FLAGS} ${GO_PKGS}

##@ test, bench and coverage

.PHONY: test
test: CGO_ENABLED=1  # needs race test
test: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
test: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
test: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
test:  ## Runs package test including race condition.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v -race $(strip $(GO_TEST_FLAGS)) $(strip $(GO_FLAGS)) -run=$(GO_TEST_FUNC) $(GO_TEST_PKGS)

.PHONY: bench
bench: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
bench: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
bench: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
bench:  ## Take a package benchmark.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_BENCH_FLAGS)) $(strip $(GO_FLAGS)) -run='^$$' -bench=$(GO_BENCH_FUNC) $(GO_TEST_PKGS)

.PHONY: coverage
coverage: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
coverage: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
coverage: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
coverage:  ## Takes packages test coverage.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_TEST_FLAGS)) $(strip $(GO_FLAGS)) -covermode=atomic -coverpkg=$(PKG)/... -coverprofile=${GO_TEST_COVERAGE_OUT} $(GO_PKGS)

.PHONY: tools/go-junit-report
tools/go-junit-report:  # go get 'go-junit-report' binary
ifeq (, $(shell command -v go-junit-report))
	@cd $(mktemp -d); \
		go mod init tmp > /dev/null 2>&1; \
		go get -u github.com/jstemmer/go-junit-report@master
endif

.PHONY: coverage/ci
coverage/ci: tools/go-junit-report
coverage/ci: GO_LDFLAGS=${GO_LDFLAGS_STATIC}
coverage/ci: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
coverage/ci: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
coverage/ci:  ## Takes packages test coverage, and output coverage results to CI artifacts.
	$(call target)
	@mkdir -p /tmp/ci/artifacts /tmp/ci/test-results
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_TEST_FLAGS)) $(strip $(GO_FLAGS)) -covermode=atomic -coverpkg=$(PKG)/... -coverprofile=${GO_TEST_COVERAGE_OUT} $(GO_PKGS) 2>&1 | tee /dev/stderr | go-junit-report -set-exit-code > /tmp/ci/test-results/junit.xml
	@if [[ -f "${GO_TEST_COVERAGE_OUT}" ]]; then go tool cover -html=${GO_TEST_COVERAGE_OUT} -o $(dir GO_TEST_COVERAGE_OUT)/coverage.html; fi


##@ lint

.PHONY: lint
lint: lint/golangci-lint  ## Run all linters.

.PHONY: tools/golangci-lint
tools/golangci-lint:  # go get 'golangci-lint' binary
ifeq (, $(shell command -v golangci-lint))
	@cd $(mktemp -d); \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

.PHONY: lint/golangci-lint
lint/golangci-lint: tools/golangci-lint .golangci.yml  ## Run golangci-lint.
	$(call target)
	@golangci-lint run ./...


##@ mod

.PHONY: mod/init
mod/init:  ## Initializes and writes a new `go.mod` to the current directory.
	$(call target)
	@go mod init > /dev/null 2>&1 || true

.PHONY: mod/get
mod/get:  ## Updates all module packages and go.mod.
	$(call target)
	@go get -u -m -v -x

.PHONY: mod/tidy
mod/tidy:  ## Makes sure go.mod matches the source code in the module.
	$(call target)
	@go mod tidy -v

.PHONY: mod/vendor
mod/vendor: mod/tidy  ## Resets the module's vendor directory and fetch all modules packages.
	$(call target)
	@go mod vendor -v

.PHONY: mod/graph
mod/graph:  ## Prints the module requirement graph with replacements applied.
	$(call target)
	@go mod graph

.PHONY: mod/install
mod/install: mod/tidy mod/vendor
mod/install:  ## Install the module vendor package as an object file.
	$(call target)
	@GO111MODULE=off go install -v $(strip $(GO_FLAGS)) $(GO_VENDOR_PKGS) || go install -v ${GO_MOD_FLAGS} $(strip $(GO_FLAGS)) $(GO_VENDOR_PKGS)

.PHONY: mod/update
mod/update: mod/get mod/tidy mod/vendor mod/install  ## Updates all of vendor packages.
	@go mod edit -go 1.13

.PHONY: mod
mod: mod/init mod/tidy mod/vendor mod/install
mod:  ## Updates the vendoring directory using go mod.
	@go mod edit -go 1.13


##@ clean

.PHONY: clean
clean:  ## Cleanups binaries and extra files in the package.
	$(call target)
	@$(RM) $(APP) *.out *.test *.prof trace.log


##@ container

.PHONY: container/build
container/build:  ## Creates the container image.
	docker image build ${CONTAINER_BUILD_ARGS} --target ${CONTAINER_BUILD_TARGET} -t $(CONTAINER_REGISTRY)/$(APP):${CONTAINER_BUILD_TAG} .

.PHONY: container/push
container/push:  ## Pushes the container image to $CONTAINER_REGISTRY.
	docker image push $(CONTAINER_REGISTRY)/$(APP):$(VERSION)


##@ boilerplate

.PHONY: boilerplate/go/%
boilerplate/go/%: BOILERPLATE_PKG_DIR=$(shell printf $@ | cut -d'/' -f3- | rev | cut -d'/' -f2- | rev | awk -F. '{print $$1}')
boilerplate/go/%: BOILERPLATE_PKG_NAME=$(if $(findstring .go,$(suffix $(BOILERPLATE_PKG_DIR))),$(basename ${@F}),$(shell printf $@ | rev | cut -d/ -f2 | rev))
boilerplate/go/%: hack/boilerplate/boilerplate.go.txt
boilerplate/go/%:  ## Creates a go file based on boilerplate.go.txt in % location.
	@if [ -n ${BOILERPLATE_PKG_DIR} ] && [ ! -d ${BOILERPLATE_PKG_DIR} ]; then mkdir -p ${BOILERPLATE_PKG_DIR}; fi
	@if [[ ${@F} == *'.go'* ]] || [[ ${BOILERPLATE_PKG_DIR} == *'cmd'* ]] || [ -z ${BOILERPLATE_PKG_DIR} ]; then \
		cat hack/boilerplate/boilerplate.go.txt <(printf "\npackage $(basename ${@F})\\n") > $*; \
		else \
		cat hack/boilerplate/boilerplate.go.txt <(printf "\npackage ${BOILERPLATE_PKG_NAME}\\n") > $*; \
		fi
	@sed -i "s|YEAR|$(shell date '+%Y')|g" $*

## miscellaneous

.PHONY: AUTHORS
AUTHORS:  ## Creates AUTHORS file.
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@printf "$(shell git log --format="\n%aN <%aE>" | LC_ALL=C.UTF-8 sort -uf)" >> $@

.PHONY: TODO
TODO:  ## Print the all of (TODO|BUG|XXX|FIXME|NOTE) in packages.
	@rg -e '(TODO|BUG|XXX|FIXME|NOTE)(\(.+\):|:)' --follow --hidden --glob='!.git' --glob='!vendor' --glob='!internal' --glob='!Makefile' --glob='!snippets' --glob='!indent'


# ----------------------------------------------------------------------------
##@ help

.PHONY: help
help:  ## Show make target help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
