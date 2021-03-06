# ----------------------------------------------------------------------------
# global

.DEFAULT_GOAL = static
APP = kt
CMD_PREFIX = $(PKG)/cmd/
CMD = $(CMD_PREFIX)$(APP)

# ----------------------------------------------------------------------------
# target

.PHONY: all
all: mod pkg/install static

# ----------------------------------------------------------------------------
# include

include hack/make/go.mk

# ----------------------------------------------------------------------------
# overlays

override GO_PKGS := $(shell go list ${GO_MOD_FLAGS} -f '{{if and (or .GoFiles .CgoFiles) (ne .Name "main")}}{{.ImportPath}}{{end}}' ./...)
