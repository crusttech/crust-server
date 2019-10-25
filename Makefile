.PHONY: help docker docker-push realize dep dep.update test test.messaging test.compose qa critic vet codegen integration

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u
GOTEST    ?= go test

BASEPKGS = system compose messaging
IMAGES   = corteza-server-system corteza-server-compose corteza-server-messaging corteza-server
TESTABLE = messaging system compose pkg internal

# Run watcher with a different event-trigger delay, eg:
# $> WATCH_DELAY=5s make watch.test.integration
WATCH_DELAY ?= 1s

# Run go test cmd with flags, eg:
# $> TEST_FLAGS="-v" make test.integration
# $> TEST_FLAGS="-v -run SpecialTest" make test.integration
TEST_FLAGS ?=

# Cover package maps for tests tasks
COVER_PKGS_messaging   = ./messaging/...
COVER_PKGS_system      = ./system/...
COVER_PKGS_compose     = ./compose/...
COVER_PKGS_pkg         = ./pkg/...
COVER_PKGS_all         = $(COVER_PKGS_pkg),$(COVER_PKGS_messaging),$(COVER_PKGS_system),$(COVER_PKGS_compose)
COVER_PKGS_integration = $(COVER_PKGS_all)

TEST_SUITE_pkg         = ./pkg/...
TEST_SUITE_services    = ./compose/... ./messaging/... ./system/...
TEST_SUITE_unit        = $(TEST_SUITE_pkg) $(TEST_SUITE_services)
TEST_SUITE_integration = ./tests/...
TEST_SUITE_all         = $(TEST_SUITE_unit) $(TEST_SUITE_integration)


########################################################################################################################
# Tool bins
REALIZE     = ${GOPATH}/bin/realize
GOTEST      = ${GOPATH}/bin/gotest
GOCRITIC    = ${GOPATH}/bin/gocritic
MOCKGEN     = ${GOPATH}/bin/mockgen
STATICCHECK = ${GOPATH}/bin/staticcheck

help:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - docker-images: builds docker images locally
	@echo - docker-push:   push built images
	@echo
	@echo - vet - run go vet on all code
	@echo - critic - run go critic on all code
	@echo - test - run all available unit tests
	@echo - qa - run vet, critic and test on code
	@echo


docker-images: $(IMAGES:%=docker-image.%)

docker-image.%: Dockerfile.%
	@ docker build --no-cache --rm -f Dockerfile.$* -t crusttech/$*:latest .

docker-push: $(IMAGES:%=docker-push.%)

docker-push.%: Dockerfile.%
	@ docker push crusttech/$*:latest


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

cdeps:
	$(GO) get github.com/cortezaproject/corteza-server
	$(GO) mod vendor

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

watch.test.%: $(NODEMON)
	# Development helper - watches for file
	# changes & reruns  tests
	$(WATCHER) "make test.$* || exit 0"

########################################################################################################################
# Quality Assurance

# Adds -coverprofile flag to test flags
# and executes test.cover... task
test.coverprofile.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverprofile=$(COVER_PROFILE)" make test.cover.$*

# Adds -coverpkg flag
test.cover.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverpkg=$(COVER_PKGS_$*)" make test.$*

# Runs integration tests
test.integration:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_integration)

# Runs one suite from integration tests
test.integration.%:
	$(GOTEST) $(TEST_FLAGS) ./tests/$*/...

# Runs ALL tests
test.all:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_all)

# Unit testing testing messaging, system or compose
test.unit.%:
	$(GOTEST) $(TEST_FLAGS) ./$*/...

# Runs ALL tests
test.unit:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_unit)

# Testing pkg
test.pkg:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_pkg)

test: test.unit

vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

qa: vet critic test

########################################################################################################################
# Toolset

$(GOTEST):
	$(GOGET) github.com/rakyll/gotest

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

$(STATICCHECK):
	$(GOGET) honnef.co/go/tools/cmd/staticcheck

clean:
	rm -f $(REALIZE) $(GOCRITIC) $(GOTEST)
