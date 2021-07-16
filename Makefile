# Project related variables
PKG = github.com/dp1140a/geoip
APP_NAME=geoip

# Directories
WD := $(subst $(BSLASH),$(FSLASH),$(shell pwd))
MD := $(subst $(BSLASH),$(FSLASH),$(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))
BUILD_DIR = $(WD)/build
PKG_DIR = $(MD)
CMD_DIR = $(PKG_DIR)/cmd
DIST_DIR = $(WD)/dist
LOG_DIR = $(WD)/log
REPORT_DIR = $(WD)/reports

M = $(shell printf "\033[34;1m▶\033[0m")
DONE="$(M) done ✨"
VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
ifndef VERSION
	VERSION := dev
endif
GIT_TAG := $(shell git describe --exact-match --tags 2>git_describe_error.tmp; rm -f git_describe_error.tmp)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-parse HEAD)

MAKEFLAGS += --no-print-directory
CMDS := $(shell find "$(CMD_DIR)/" -mindepth 1 -maxdepth 1 -type f | sed 's/ /\\ /g' | xargs -n1 basename)

GOBIN = $(shell go env GOPATH)/bin
ARCHES ?= amd64
OSES ?= linux darwin
OUTTPL = $(DIST_DIR)/$(APP_NAME)-$(VERSION)-{{.OS}}_{{.Arch}}
GZCMD = tar -czf
ZIPCMD = zip
SHACMD = sha256sum
VET_RPT=vet.out
COVERAGE_RPT=coverage.out

LDFLAGS = -X $(PKG)/version.APP_NAME=$(APP_NAME) \
	-X $(PKG)/version.commit=$(GIT_COMMIT) \
	-X $(PKG)/version.branch=$(GIT_BRANCH) \
	-X $(PKG)/version.version=$(VERSION) \
	-X $(PKG)/version.buildTime=$(shell date -Iseconds)

## deps: Download and Install any missing dependecies
.PHONY: deps
deps:
	go mod download
	@echo $(DONE) "-- Deps"

## build: Install missing dependencies. Builds binary in ./build
.PHONY: build
build: deps tidy fmt reports
	@mkdir -pv $(BUILD_DIR)
	@echo "$(LDFLAGS)"
	@echo "  $(M)  Checking if there is any missing dependencies...\n"
	@$(MAKE) deps
	@echo "  $(M)  Building...\n"
	@echo "GOBIN: $(GOBIN)"
	$(GOBIN)/gox -rebuild -gocmd="go" -arch="$(ARCHES)" -os="$(OSES)" -output="$(OUTTPL)/{{.Dir}}" \
		-tags "$(BUILD_TAGS)" -ldflags "$(LDFLAGS)"
	$(info "Built version:$(VERSION), build:$(GIT_COMMIT)")
	@echo $(DONE) "-- Build"

## Creates a distribution
.PHONY: dist
dist: clean build
	cd "$(DIST_DIR)"; for dir in ./**; do \
		cp $(PKG_DIR)/config.toml $$dir; \
		cp -r $(PKG_DIR)/etc $$dir; \
		$(GZCMD) "$(basename "$$dir").tar.gz" "$$dir"; \
	done
	cd "$(DIST_DIR)"; find . -maxdepth 1 -type f -printf "$(SHACMD) %P | tee \"./%P.sha\"\n" | sh
	$(info "Built v$(VERSION), build $(COMMIT_ID)")
	@echo $(DONE) "-- Dist"

## docker: Builds a docker image
.PHONY: docker
docker:
	docker build --build-arg="$(APP_NAME)" \
	--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
	--build-arg GIT_BRANCH="$(GIT_BRANCH)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg GOOS="$(GOOS)" \
	-t fofx/fofx:$(VERSION) . -f Dockerfile
	@echo $(DONE) "-- Docker"

## tidy: Verifies and downloads all required dependencies
.PHONY: tidy
tidy:
	@echo "$(M) 🏃 go mod tidy..."
	@mkdir -pv $(REPORT_DIR)
	go mod verify
	go mod tidy
	@if ! git diff --quiet; then \
		echo "'go mod tidy' resulted in changes or working tree is dirty:"; \
		git --no-pager diff > $(REPORT_DIR)/diff.out; \
	fi
	@echo $(DONE) "-- Tidy"

## fmt: Runs gofmt on all source files
.PHONY: fmt
fmt:
	@echo "$(M) 🏃 gofmt..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./...); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret
	@echo $(DONE) "-- Tidy"

## test: Tests code coverage
.PHONY: test
test:
	@echo "$(M)  👀 testing code...\n"
	@mkdir -pv $(REPORT_DIR)
	go test ./... >$(REPORT_DIR)/test.out 2>&1
	@echo $(DONE) "-- Test"

## testwithcoverge: Tests code coverage
.PHONY: testwithcoverage
testwithcoverage:
	@echo "$(M)  👀 testing code with coverage...\n"
	@mkdir -pv $(REPORT_DIR)
	go test ./... -coverprofile=$(REPORT_DIR)/$(COVERAGE_RPT)
	@echo $(DONE) "-- Test with Coverage"

## missing: Displays lines of code missing from coverage. Puts report in ./build/coverage.out
.PHONY: missing
missing: testwithcoverage
	@echo "$(M)  👀 missing coverage...\n"
	@mkdir -pv $(REPORT_DIR)
	go tool cover -func=$(REPORT_DIR)/$(COVERAGE_RPT) -o $(REPORT_DIR)/missing.out
	@echo $(DONE) "-- Missing"

## vet: Run go vet.  Puts report in ./build/vet.out
.PHONY: vet
vet:
	@echo "  $(M) 🏃 go vet..."
	@mkdir -pv $(REPORT_DIR)
	go vet -v ./... 2>&1 | tee $(REPORT_DIR)/vet.out
	@echo $(DONE) "-- Vet"

## reports: Runs vet, coverage, and missing reports
.PHONY: reports
reports: vet missing
	@echo $(DONE) "-- Reports"

## clean: Removes build, dist and report dirs
.PHONY: clean
clean:
	@echo "$(M)  🧹 Cleaning build ..."
	go clean $(PKG) || true
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -rf $(REPORT_DIR)
	@echo $(DONE) "-- Clean"

## cleanDocker: Clean build files. Runs `go clean` internally
.PHONY: cleanDocker
cleanDocker:
	@echo "$(M)  🧹 Cleaning docker build ..."
	#docker stop $(APP_NAME)
	$(shell docker rm $(APP_NAME) 2>/dev/null)
	$(shell docker rmi $$(docker images -f "dangling=true" -q) 2>/dev/null)
	@echo $(DONE) "-- Clean Docker"

## gencerts: Generates a sample self signed cert and key to enable TLS
.PHONY: gencerts
gencerts:
	@echo "$(M) Generating SSL certs"
	$(shell openssl req -newkey rsa:4096 \
		-x509 \
		-sha256 \
    	-days 3650 \
    	-nodes \
    	-out etc/example.crt \
    	-keyout etc/example.key \
    	-subj "/C=US/ST=WA/L=Seattle/O=Example Org/OU=Example/CN=www.example.com" \
    ) @echo $(DONE) "-- Gen Certs"

## debug: Print make env information
.PHONY: debug
debug:
	$(info PKG=$(PKG))
	$(info APP_NAME=$(APP_NAME))
	$(info MD=$(MD))
	$(info WD=$(WD))
	$(info PKG_DIR=$(PKG_DIR))
	$(info CMD_DIR=$(CMD_DIR))
	$(info BUILD_DIR=$(BUILD_DIR))
	$(info DIST_DIR=$(DIST_DIR))
	$(info LOG_DIR=$(LOG_DIR))
	$(info REPORT_DIR=$(REPORT_DIR))
	$(info VERSION=$(VERSION))
	$(info GIT_COMMIT=$(GIT_COMMIT))
	$(info GIT_TAG=$(GIT_TAG))
	$(info CMDS=$(CMDS))
	$(info ARCHES=$(ARCHES))
	$(info OSES=$(OSES))
	$(info LDFLAGS=$(LDFLAGS))
	@echo $(DONE) "-- Debug"

.PHONY: help
help: Makefile
	@echo "\n Choose a command run in "$(PROJECTNAME)":\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'