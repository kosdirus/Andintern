export GO111MODULE=on
export GOFLAGS=-mod=vendor
export GOBIN=$(PWD)/gobin

GO_PACKAGE := github.com/kosdirus/andintern
OUT := bin/main
GIT_BRANCH := $$(git rev-parse --abbrev-ref HEAD)
GIT_REV := $$(git rev-parse HEAD)
GIT_VER := $$(git describe --tags --abbrev=0 --always)

clean:
	rm -rf ./bin/*

vendor:
	go mod download
	go mod tidy
	go mod vendor
.PHONY: vendor

build: clean
	go build -o $(OUT) ./cmd
.PHONY: build

run: build
	$(OUT) -config ./config.json
.PHONY: run