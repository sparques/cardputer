GO_TOOLCHAIN_VERSION := 1.24.6
GO_TOOLCHAIN_CACHE := $(HOME)/go/pkg/mod/golang.org/toolchain@v0.0.1-go$(GO_TOOLCHAIN_VERSION).linux-amd64/bin
BOARD ?= cardputer

ifeq ($(BOARD),cardputer-adv)
BOARD_TAGS := cardputer_adv
else
BOARD_TAGS :=
endif

ifeq ($(wildcard $(GO_TOOLCHAIN_CACHE)/go),)
TINYGO_ENV = GOTOOLCHAIN=go$(GO_TOOLCHAIN_VERSION)
else
TINYGO_ENV = PATH=$(GO_TOOLCHAIN_CACHE):$$PATH
endif

CACHE_ENV = GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod XDG_CACHE_HOME=$(CURDIR)/.cache

.PHONY: build tinygo-version clean

build:
	env $(TINYGO_ENV) $(CACHE_ENV) tinygo test -target=esp32s3-generic $(if $(BOARD_TAGS),-tags=$(BOARD_TAGS)) ./...

tinygo-version:
	env $(TINYGO_ENV) tinygo version

clean:
	rm -rf .cache
