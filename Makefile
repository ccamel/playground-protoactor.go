.EXPORT_ALL_VARIABLES:

.PHONY: tools deps gen-static check build

default: build

tools: ./bin/golangci-lint $(GOPATH)/bin/esc $(GOPATH)/bin/gothanks
protos=$(addsuffix .pb.go,$(basename $(shell find . -maxdepth 5 -type f -name *.proto)))

%.pb.go: %.proto
	@echo "🖋 Generating proto $(notdir $<)"
	protoc -I=$(dir $@) -I=$(GOPATH)/src -I=$(GOPATH)/pkg/mod --go_out=$(dir $@) $(notdir $<)

deps:
	go get .

protobuf: $(protos)

check: tools
	./bin/golangci-lint run ./...

thanks: tools
	$(GOPATH)/bin/gothanks -y | grep -v "is already"

build:
	go build -o playground-protoactor .

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
