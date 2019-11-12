BUILD_FOLDER = dist
BUILD_OPTIONS =

OS := darwin
ARCH := 386
BINARY_NAME := "vimpack"
FULL_BINARY_NAME := $(BINARY_NAME)-$(OS)-$(ARCH)

VERSION := $$(cat VERSION)

PROJECT_USERNAME := kdisneur
PROJECT_REPOSITORY := vimpack
GITHUB_TOKEN :=

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git branch --no-color | awk '/^\* / { print $$2 }')
GIT_STATE := $(shell if [ -z "$(shell git status --short)" ]; then echo clean; else echo dirty; fi)
ALREADY_RELEASED := $(shell if [ $$(curl --silent --output /dev/null --write-out "%{http_code}" https://api.github.com/repos/$(PROJECT_USERNAME)/$(PROJECT_REPOSITORY)/releases/tags/$(VERSION)) -eq 200 ]; then echo "true"; else echo "false"; fi)
PRERELASE := $(shell if [ "$(GIT_BRANCH)" != "master" ]; then echo "true"; else echo "false"; fi)

GOFMT_PATH = gofmt
GOLINT_PATH = golint
GHR_PATH = ghr
STATICCHECK_PATH = staticcheck

PACKAGE = internal
MOCK_PACKAGE = $(PACKAGE)/mock_$(PACKAGE)

MOCK_FILES=$(wildcard $(MOCK_PACKAGE)/*)

TEST_OPTIONS =

$(MOCK_PACKAGE)/%.go: $(PACKAGE)/%.go
	mockgen -source $< -destination $@

setup:
	go get github.com/tcnksm/ghr
	go get golang.org/x/lint/golint
	go get github.com/golang/mock/mockgen
	go get honnef.co/go/tools/cmd/staticcheck

build: _dependencies
	@touch internal/version.go

	GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILD_OPTIONS) \
		-ldflags \
			"-X vimpack/internal.versionNumber=$(VERSION) \
			 -X vimpack/internal.gitBranch=$(GIT_BRANCH) \
			 -X vimpack/internal.gitCommit=$(GIT_COMMIT) \
			 -X vimpack/internal.gitState=$(GIT_STATE) \
			 -X vimpack/internal.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		-o $(BUILD_FOLDER)/$(BINARY_NAME)

	@tar czf $(BUILD_FOLDER)/$(FULL_BINARY_NAME).tgz -C $(BUILD_FOLDER) $(BINARY_NAME)
	@echo "archive generated at $(BUILD_FOLDER)/$(FULL_BINARY_NAME).tgz"

	@mv $(BUILD_FOLDER)/$(BINARY_NAME) $(BUILD_FOLDER)/$(FULL_BINARY_NAME)
	@echo "archive generated at $(BUILD_FOLDER)/$(FULL_BINARY_NAME)"

release:
	@if [ $(ALREADY_RELEASED) = "true" ]; then \
	  echo "Already released... skipping" \
	  exit 0; \
	fi; \
	tmpfolder=$$(mktemp -d /tmp/$(PROJECT_REPOSITORY)-artifacts-XXXXX);\
	cp $(BUILD_FOLDER)/*.tgz $${tmpfolder}; \
	$(GHR_PATH) -t $(GITHUB_TOKEN) \
		-u $(PROJECT_USERNAME) \
		-r $(PROJECT_REPOSITORY) \
		-c $(GIT_COMMIT) \
		$$(if [ "$(PRERELASE)" = "true" ]; then echo "-prerelease"; else echo ""; fi) \
		$(VERSION) $${tmpfolder};

refresh-mocks: $(MOCK_FILES)

refresh-templates:
	@version=$$(awk '$$1 == "golang" { print $$2 }' .tool-versions);\
	tempfile=$$(mktemp /tmp/circleci.config.XXXXX);\
	echo "# THIS FILE IS GENERATED, DO NOT EDIT MANUALLY.\n# SOURCE: config.yml.in" > $${tempfile};\
	m4 -D GO_VERSION=$${version} .circleci/config.yml.in >> $${tempfile};\
	mv $${tempfile} .circleci/config.yml

test-circleci-config:
	@tempfile=$$(mktemp /tmp/circleci.config.XXXXX);\
	cp .circleci/config.yml $${tempfile};\
	make refresh-templates;\
	if ! diff $${tempfile} .circleci/config.yml; then \
	  echo "---";\
	  echo "<: actual | >: expected";\
	  cp $${tempfile} .circleci/config.yml; \
	  exit 1; \
	fi

test-style: _gofmt _golint

test-unit:
	go test $(TEST_OPTIONS) ./...

test-staticcheck:
	$(STATICCHECK_PATH) ./...

_dependencies:
	go mod download

_gofmt:
	@data=$$($(GOFMT_PATH) -l .);\
	if [ -n "$$data" ]; then \
		echo $$data; \
		exit 1; \
	fi
_golint:
	@data=$$($(GOLINT_PATH) internal . | grep -v "should have comment");\
	if [ -n "$$data" ]; then \
		echo $$data; \
		exit 1; \
	fi
