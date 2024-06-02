# TODO(adamc): this is taken from Crossplane's make system. Let's just switch to using that.

# ====================================================================================
# Logger

TIME_LONG	= `date +%Y-%m-%d' '%H:%M:%S`
TIME_SHORT	= `date +%H:%M:%S`
TIME		= $(TIME_SHORT)

INFO	= echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN	= echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR		= echo ${TIME} ${RED}[FAIL]${CNone}
OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)

.PHONY: all
all: test

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: install-go-tools
install-go-tools: ## Install dev tools
	go generate -tags tools

.PHONY: generate
generate: install-go-tools ## Generates any Kelpie mocks
	go generate ./...

# prepare for code review
reviewable:
	@$(MAKE) generate
	@$(MAKE) lint
	@$(MAKE) test

.PHONY: check-diff
# ensure generate target doesn't create a diff
check-diff: generate
	@$(INFO) checking that branch is clean
	git diff
	@if git status --porcelain | grep . ; then $(ERR) There are uncommitted changes after running make generate. Please ensure you commit all generated files in this branch after running make generate. && false; else $(OK) branch is clean; fi

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test -coverprofile cover.out ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run
