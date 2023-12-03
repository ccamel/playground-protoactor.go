.EXPORT_ALL_VARIABLES:

# Docker images
DOCKER_IMAGE_BUF = bufbuild/buf:1.28.1

.PHONY: tools deps gen-static check build lint lint-proto

default: build

tools: ./bin/golangci-lint $(GOPATH)/bin/esc $(GOPATH)/bin/gothanks

deps:
	go get .

protobuf:
	@echo "🖋 Generating proto..."
	@docker run --rm \
		-v `pwd`:/proto \
		-w /proto \
		${DOCKER_IMAGE_BUF} \
		generate --verbose

check: tools
	./bin/golangci-lint run ./...

thanks: tools
	$(GOPATH)/bin/gothanks -y | grep -v "is already"

build:
	go build -o playground-protoactor .

lint: lint-proto

lint-proto:
	@echo "🖋 Linting proto..."
	@docker run --rm \
		-v `pwd`:/proto \
		-w /proto \
		${DOCKER_IMAGE_BUF} \
		generate --verbose

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
