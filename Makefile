.EXPORT_ALL_VARIABLES:

# Docker images
DOCKER_IMAGE_GOLANG	 = golang:1.21-alpine3.17
DOCKER_IMAGE_BUF     = bufbuild/buf:1.28.1

default: build

.PHONY: tools
tools: ./bin/golangci-lint $(GOPATH)/bin/esc $(GOPATH)/bin/gothanks

.PHONY: deps
deps:
	go get .

.PHONY: protobuf
protobuf:
	@echo "🖋 Generating proto..."
	@docker run --rm \
		-v `pwd`:/proto \
		-w /proto \
		${DOCKER_IMAGE_BUF} \
		generate --verbose

.PHONY: check
check: tools
	./bin/golangci-lint run ./...

.PHONY: thanks
thanks: tools
	$(GOPATH)/bin/gothanks -y | grep -v "is already"

.PHONY: build
build:
	go build -o playground-protoactor .

.PHONY: lint
lint: lint-proto

.PHONY: lint-proto
lint-proto:
	@echo "🖋 Linting proto..."
	@docker run --rm \
		-v `pwd`:/work \
		-w /work \
		${DOCKER_IMAGE_BUF} \
		generate --verbose

.PHONY: format
format: format-go

format-go:
	@echo "📐 Formatting go source code..."
	@docker run --rm \
  		-v `pwd`:/work:rw \
  		-w /work \
  		${DOCKER_IMAGE_GOLANG} \
  		sh -c \
		"go install mvdan.cc/gofumpt@v0.5.0; gofumpt -w -l ."

.PHONY: docker
docker:
	@echo "📦 building container"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o playground-protoactor.amd64 .
	docker build .

$(GOPATH)/bin/gothanks:
	@echo "📦 installing $(notdir $@)"
	go get -u github.com/psampaz/gothanks

$(GOPATH)/bin/esc:
	@echo "📦 installing $(notdir $@)"
	go get -u github.com/mjibson/esc

$(GOPATH)/bin/protoc:
	@echo "📦 installing $(notdir $@)"
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go

./bin/golangci-lint:
	@echo "📦 installing $(notdir $@)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.32.2
