.PHONY: help docker docker-push realize dep dep.update test test.messaging test.compose qa critic vet codegen integration

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

BASEPKGS = system compose messaging
IMAGES   = crust-server-system crust-server-compose crust-server-messaging crust-server

########################################################################################################################
# Tool bins
DEP         = $(GOPATH)/bin/dep
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

dep.update: $(DEP)
	$(DEP) ensure -update -v

cdeps: $(DEP)
	$(DEP) ensure -update github.com/cortezaproject/corteza-server
	$(DEP) ensure -v

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

########################################################################################################################
# QA

test:
	# Run basic unit tests
	$(GO) test ./opt/... ./internal/... ./compose/... ./messaging/... ./system/...

integration:
	# Run drone's integration pipeline
	rm -f build/gen*
	drone exec --pipeline integration

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

$(DEP):
	$(GOGET) github.com/tools/godep

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

$(STATICCHECK):
	$(GOGET) honnef.co/go/tools/cmd/staticcheck

clean:
	rm -f $(REALIZE) $(GOCRITIC) $(GOTEST)
