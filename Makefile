.PHONY: build release upload realize cdeps

GO        = go
GOGET     = $(GO) get -u
GOTEST    ?= go test
GOFLAGS   ?= -mod=vendor

export GOFLAGS

BUILD_FLAVOUR         ?= crust
BUILD_APPS            ?= system compose messaging monolith
BUILD_TIME            ?= $(shell date +%FT%T%z)
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_ARCH            ?= $(shell go env GOARCH)
BUILD_OS              ?= $(shell go env GOOS)
BUILD_OS_is_windows    = $(filter windows,$(BUILD_OS))
BUILD_DEST_DIR        ?= build
BUILD_NAME             = $(BUILD_FLAVOUR)-server-$*-$(BUILD_VERSION)-$(BUILD_OS)-$(BUILD_ARCH)
BUILD_BIN_NAME         = $(BUILD_NAME)$(if $(BUILD_OS_is_windows),.exe,)

RELEASE_BASEDIR        = $(BUILD_DEST_DIR)/pkg/$(BUILD_FLAVOUR)-server-$*
RELEASE_NAME           = $(BUILD_NAME).tar.gz
RELEASE_EXTRA_FILES   ?= README.md LICENSE CONTRIBUTING.md DCO .env.example
RELEASE_PKEY          ?= .upload-rsa

LDFLAGS_BUILD_TIME     = -X github.com/cortezaproject/corteza-server/pkg/version.BuildTime=$(BUILD_TIME)
LDFLAGS_VERSION        = -X github.com/cortezaproject/corteza-server/pkg/version.Version=$(BUILD_VERSION)
LDFLAGS_EXTRA         ?=
LDFLAGS                = -ldflags "$(LDFLAGS_BUILD_TIME) $(LDFLAGS_GIT_TAG) $(LDFLAGS_EXTRA)"

########################################################################################################################

help:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - build             build all apps
	@echo - build.<app>       build a specific app
	@echo - release           release all apps
	@echo - release.<app>     release a specific app
	@echo

########################################################################################################################
# Building & packing

build: $(addprefix build., $(BUILD_APPS))

build.%: cmd/%
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build $(LDFLAGS) -o $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME) cmd/$*/main.go

release.%: $(addprefix build., %)
	@ mkdir -p $(RELEASE_BASEDIR) $(RELEASE_BASEDIR)/bin
	@ cp $(RELEASE_EXTRA_FILES) $(RELEASE_BASEDIR)/
	@ cp $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME) $(RELEASE_BASEDIR)/bin/$(BUILD_FLAVOUR)-server-$*
	@ tar -C $(dir $(RELEASE_BASEDIR)) -czf $(BUILD_DEST_DIR)/$(RELEASE_NAME) $(notdir $(RELEASE_BASEDIR))

release: $(addprefix release.,$(BUILD_APPS))

release-clean:
	@ rm -rf $(RELEASE_BASEDIR)

upload: $(RELEASE_PKEY)
	@ echo "put $(BUILD_DEST_DIR)/*.tar.gz" | sftp -q -i $(RELEASE_PKEY) $(RELEASE_SFTP_URI)
	@ rm -f $(RELEASE_PKEY)

$(RELEASE_PKEY):
	@ echo $(RELEASE_SFTP_KEY) | base64 -d > $(RELEASE_PKEY)
	@ chmod 0400 $@

########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

cdeps:
	$(GO) get github.com/cortezaproject/corteza-server
	$(GO) mod vendor

