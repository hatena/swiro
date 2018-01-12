NAME     := swiro
VERSION := $(shell git describe --tags --exact-match 2> /dev/null || git rev-parse --short HEAD || echo "unknown")

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -s -w -extldflags "-static"
LDFLAGS += \
	-X "github.com/taku-k/swiro/build.version=$(VERSION)" \
	-X "github.com/taku-k/swiro/build.name=$(NAME)"

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o bin/$(NAME)

.PHONY: cross-build
cross-build: deps
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_darwin_amd64
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_darwin_386
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_linux_amd64
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_linux_386
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_windows_amd64.exe
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '$(LDFLAGS)' -o dist/$(NAME)_windows_386.exe

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
    go get github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep
	dep ensure

.PHONY: fmt
fmt:
	gofmt -s -w $$(find . -type f -name '*.go' | grep -v -e vendor)

.PHONY: imports
imports:
	goimports -w $$(find . -type f -name '*.go' | grep -v -e vendor)

.PHONY: test
test:
	go test -v -race $$(go list ./...)

