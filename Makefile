SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.EXPORT_ALL_VARIABLES:

GOPATH ?= $(shell go env GOPATH)
CURDIR := $(shell pwd)

# Versions
GOLANG_VERSION          ?= 1.23
ALPINE_VERSION          ?= 3.20
BUF_VERSION             ?= 1.45.0
GOLANGCI_LINT_VERSION   ?= v2.0.2
GOFUMPT_VERSION         ?= v0.7.0
GOTHANKS_VERSION        ?= latest
PROTOC_GEN_GO_VERSION   ?= latest
DOCKER_IMAGE_TAG        ?= latest

# Docker images
DOCKER_IMAGE_BUF    = bufbuild/buf:$(BUF_VERSION)

# Binary names
BINARY_NAME   = playground-protoactor
BINARY_AMD64  = $(BINARY_NAME).amd64

# Directories
TOOLS_DIR = ./bin

GOLANGCI_LINT_BIN = $(TOOLS_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)/golangci-lint
GOTHANKS_BIN      = $(TOOLS_DIR)/gothanks/$(GOTHANKS_VERSION)/gothanks
GOFUMPT_BIN       = $(TOOLS_DIR)/gofumpt/$(GOFUMPT_VERSION)/gofumpt
PROTOC_GEN_GO_BIN = $(TOOLS_DIR)/protoc-gen-go/$(PROTOC_GEN_GO_VERSION)/protoc-gen-go

# Some colors (if supported)
define get_color
$(shell tput -Txterm $(1) $(2) 2>/dev/null || echo "")
endef

COLOR_GREEN  = $(call get_color,setaf,2)
COLOR_YELLOW = $(call get_color,setaf,3)
COLOR_WHITE  = $(call get_color,setaf,7)
COLOR_CYAN   = $(call get_color,setaf,6)
COLOR_RED    = $(call get_color,setaf,1)
COLOR_RESET  = $(call get_color,sgr0,)

default: help

.PHONY: check-deps
check-deps: ## Check for required external tools (docker, curl)
	$(call echo_msg, 🛠, Checking, dependencies, ...)
	@command -v docker > /dev/null || { echo "Error: docker is not installed. Aborting." >&2; exit 1; }
	@command -v curl > /dev/null || { echo "Error: curl is not installed. Aborting." >&2; exit 1; }

.PHONY: deps
deps: ## Download Go module dependencies
	$(call echo_msg, 📥, Downloading, dependencies, ...)
	@go mod download

.PHONY: protobuf
protobuf: check-deps ## Generate protobuf files
	$(call echo_msg, 🖋, Generating, proto files, using $(COLOR_YELLOW)buf $(BUF_VERSION)$(COLOR_RESET)...)
	@docker run --rm \
		-v $(CURDIR):/proto \
		-w /proto \
		$(DOCKER_IMAGE_BUF) \
		generate --verbose

.PHONY: tools
tools: $(GOLANGCI_LINT_BIN) $(GOTHANKS_BIN) $(PROTOC_GEN_GO_BIN) $(GOFUMPT_BIN) ## Install necessary development tools

.PHONY: thanks
thanks: tools ## Thanks to the contributors
	$(call echo_msg, 🙏, Running, gothanks, ...)
	@$(GOTHANKS_BIN) -y | grep -v "is already"

.PHONY: build
build: deps ## Build the project
	$(call echo_msg, 🛠, Building, project, ...)
	@go build -o $(BINARY_NAME) .

.PHONY: test
test: deps ## Run tests
	$(call echo_msg, 🛠, Building, project, ...)
	@go test -v ./...

.PHONY: lint
lint: lint-proto lint-go ## Lint files

.PHONY: lint-proto
lint-proto: check-deps ## Lint proto files
	$(call echo_msg, 🖋, Linting, proto files, using $(COLOR_YELLOW)buf $(BUF_VERSION)$(COLOR_RESET)...)
	@docker run --rm \
		-v $(CURDIR):/work \
		-w /work \
		$(DOCKER_IMAGE_BUF) \
		lint

.PHONY: lint-go
lint-go: tools ## Lint Go source code
	$(call echo_msg, 🔍, Linting, Go code, ...)
	@$(GOLANGCI_LINT_BIN) run ./...

.PHONY: format
format: format-go ## Format files

.PHONY: format-go
format-go: tools ## Format Go source code
	$(call echo_msg, 📐, Formatting, Go source code, ...)
	@$(GOFUMPT_BIN) -w -l .

.PHONY: docker
docker: build ## Build Docker container
	$(call echo_msg, 📦, Building, Docker container, ...)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_AMD64) .
	@docker build -t $(BINARY_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: clean
clean: clean-artifacts clean-tools ## Clean up

.PHONY: clean-artifacts
clean-artifacts: ## Clean up build artifacts
	$(call echo_msg, 🧹, Cleaning, build artifacts, ...)
	@rm -f $(BINARY_NAME) $(BINARY_AMD64)

.PHONY: clean-tools
clean-tools: ## Clean up tools
	$(call echo_msg, 🧹, Cleaning, tools, ...)
	@rm -rf $(TOOLS_DIR)

.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${COLOR_YELLOW}make${COLOR_RESET} ${COLOR_GREEN}<target>${COLOR_RESET}'
	@echo ''
	@echo 'Targets:'
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic,$(origin $V)), \
			$(if $(filter TOOL_%,$V), \
				export $V="$($V)";))) \
	awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${COLOR_YELLOW}%-20s${COLOR_GREEN}%s${COLOR_RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${COLOR_CYAN}%s${COLOR_RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST) | envsubst

$(TOOLS_DIR):
	@mkdir -p $(TOOLS_DIR)

$(GOLANGCI_LINT_BIN): | $(TOOLS_DIR)
	$(call echo_msg, 📦, Installing, golangci-lint, $(COLOR_YELLOW)$(GOLANGCI_LINT_VERSION)$(COLOR_RESET)...)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@mkdir -p $(dir $(GOLANGCI_LINT_BIN))
	@cp $(shell go env GOPATH)/bin/golangci-lint $(dir $(GOLANGCI_LINT_BIN))

$(GOTHANKS_BIN):
	$(call echo_msg, 📦, Installing, gothanks, $(COLOR_YELLOW)$(GOTHANKS_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(GOTHANKS_BIN))
	@GOBIN="$$(cd $(dir $(GOTHANKS_BIN)) && pwd)" go install github.com/psampaz/gothanks@$(GOTHANKS_VERSION)

$(PROTOC_GEN_GO_BIN):
	$(call echo_msg, 📦, Installing, protoc, $(COLOR_YELLOW)$(PROTOC_GEN_GO_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(PROTOC_GEN_GO_BIN))
	@GOBIN="$$(cd $(dir $(PROTOC_GEN_GO_BIN)) && pwd)" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

$(GOFUMPT_BIN):
	$(call echo_msg, 📦, Installing, gofumpt, $(COLOR_YELLOW)$(GOFUMPT_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(GOFUMPT_BIN))
	@GOBIN="$$(cd $(dir $(GOFUMPT_BIN)) && pwd)" go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

# $(call echo_msg, <emoji>, <action>, <object>, <context>)
define echo_msg
	@echo "$(strip $(1)) ${COLOR_GREEN}$(strip $(2))${COLOR_RESET} ${COLOR_CYAN}$(strip $(3))${COLOR_RESET} $(strip $(4))"
endef
