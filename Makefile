# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

# Populate version variables
# Add to compile time flags
PKG := github.com/kleijnweb/git-mirror-manager
VERSION := 0.1
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
OS := $(shell uname)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

print-%: ; @echo $*=$($*)

.PHONY: help build test benchmark run after_test

default: test build

build: vendor ## Builds git-mirror-manager
	CGO_ENABLED=0 GOOS=linux go build -a ${GO_LDFLAGS_STATIC} ${PKG}

mocks: vendor ## Generate mocks
	mockery -all -dir internal

vendor: ## Runs dep ensure
	dep ensure
	# TODO remove once https://github.com/golang/dep/issues/433 is resolved
	# Workaround for OSX sed behaving differently
	case "${OS}" in \
		Darwin) find vendor -type f -name "*.go" -print0 | xargs -0 sed -i '' 's/Sirupsen\/logrus/sirupsen\/logrus/g' ;;\
		*)      find vendor -type f -name "*.go" -print0 | xargs -0 sed -i    's/Sirupsen\/logrus/sirupsen\/logrus/g' ;;\
	esac ;\
	touch $@

checkstyle: ## Run lint and fmt
	gofmt -s -w manager main.go
	golint manager
	golint main.go

test: mocks vendor ## Run tests
	GIT_TERMINAL_PROMPT=0 go test -v ./...

cover: mocks vendor ## Run tests with code coverage
	GIT_TERMINAL_PROMPT=0 go test -covermode=count -coverprofile=cover.out ./... >/dev/null
	sed -ir '/\/mock_.*\.go:/d' ./cover.out
	go tool cover -func=cover.out

run: ## Runs git-mirror-manager without building
	go run ./main.go

# Magic as explained here: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

help: ## Shows help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
