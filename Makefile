# ---------------------------------------------------------------------
# -- Load environment variables from .env file if it exists
# ---------------------------------------------------------------------
-include .env

# ---------------------------------------------------------------------
# -- Image URL to use all building/pushing image targets
# ---------------------------------------------------------------------
IMAGE_TAG ?= $(shell git describe --tags --abbrev=0)
FULL_IMAGE_TAG ?= ${REGISTRY}/${REPO}:${IMAGE_TAG}

# ---------------------------------------------------------------------
# -- ENVTEST_K8S_VERSION refers to the version of kubebuilder assets 
# --  to be downloaded by envtest binary.
# ---------------------------------------------------------------------
ENVTEST_K8S_VERSION = 1.28.0
# ---------------------------------------------------------------------
# -- CONTROLLER_GET_VERSION a version of the controller-get tool
# ---------------------------------------------------------------------
CONTROLLER_GET_VERSION = v0.14.0
# ---------------------------------------------------------------------
# -- K8s version to start a local kubernetes
# ---------------------------------------------------------------------
K8S_VERSION ?= v1.22.3
# ---------------------------------------------------------------------
# -- Helper tools version
# ---------------------------------------------------------------------
GOLANGCI_LINT_VERSION ?= v1.55.2
# ---------------------------------------------------------------------
# -- Get the currently used golang install path 
# --  (in GOPATH/bin, unless GOBIN is set)
# ---------------------------------------------------------------------
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
# ---------------------------------------------------------------------
# -- Which container tool to use
# ---------------------------------------------------------------------
CONTAINER_TOOL ?= docker
# ---------------------------------------------------------------------
# -- Which compose tool to use
# -- For example
# -- $ export COMPOSE_TOOL="nerdctl compose"
# ---------------------------------------------------------------------
COMPOSE_TOOL ?= docker-compose
# ---------------------------------------------------------------------
# -- It's required when you want to use k3s and nerdctl
# -- $ export CONTAINER_TOOL_NAMESPACE_ARG="--namespace k8s.io"
# ---------------------------------------------------------------------
CONTAINER_TOOL_NAMESPACE_ARG ?=
# -- Set additional arguments to container tool
# -- To use, set an environment variable:
# -- $ export CONTAINER_TOOL_ARGS="--build-arg GOARCH=arm64 --platform=linux/arm64"
# ---------------------------------------------------------------------
CONTAINER_TOOL_ARGS ?=
# -- Use buildx or not
# ---------------------------------------------------------------------
USE_BUILDX ?= false
# -- Platforms for buildx
# ---------------------------------------------------------------------
PLATFORMS ?= linux/amd64,linux/arm64
# -- Push image or not
# ---------------------------------------------------------------------
PUSH_IMAGE ?= false
# ---------------------------------------------------------------------
# -- A path to store binaries that are used in the Makefile
# ---------------------------------------------------------------------
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# ---------------------------------------------------------------------
# -- Rules
# ---------------------------------------------------------------------
help: ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@grep -E '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sort | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: all
all: build

.PHONY: build
build: lint fmt vet manifests ## Build the project without creating an image
	@echo "Building the project..."

.PHONY: image
image: ## Build and optionally push the container image
	docker buildx create --use || true
	docker buildx build --platform linux/amd64,linux/arm64 -t ${FULL_IMAGE_TAG} . \
		--build-arg="OPERATOR_VERSION=${IMAGE_TAG}" \
		--output type=oci,dest=my-image.tar ${CONTAINER_TOOL_ARGS}
	
ifeq ($(PUSH_IMAGE),true)
	docker buildx build --platform linux/amd64,linux/arm64 -t ${FULL_IMAGE_TAG} . \
		--build-arg="OPERATOR_VERSION=${IMAGE_TAG}" \
		--push ${CONTAINER_TOOL_ARGS}
endif

.PHONY: push
push: ## Push the image to the registry
ifneq ($(PUSH_IMAGE), false)
	$(CONTAINER_TOOL) push $(FULL_IMAGE_TAG)
endif

.PHONY: clean
clean: ## Clean up old binaries and images
	rm -rf $(LOCALBIN)/*
	rm -f my-image.tar
	rm -rf vendor
	rm -f go.sum

# ---------------------------------------------------------------------
# -- Go related rules
# ---------------------------------------------------------------------
lint: ## lint go code
	@go mod tidy
	test -s $(LOCALBIN)/golangci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
	$(LOCALBIN)/golangci-lint run ./...

fmt: ## Format go code
	@test -s $(LOCALBIN)/gofumpt || GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@latest
	@$(LOCALBIN)/gofumpt -l -w .

vet: ## go vet to find issues
	@go vet ./...

# ---------------------------------------------------------------------
# -- Tests are not passing when GOARCH is not amd64, so it's hardcoded
# ---------------------------------------------------------------------
unit: ## run go unit tests
	GOARCH=amd64 go test -tags tests -run "TestUnit" ./... -v

test: ## run go all test
	$(COMPOSE_TOOL) down
	$(COMPOSE_TOOL) up -d
	$(COMPOSE_TOOL) restart sqladmin
	sleep 10
	GOARCH=amd64 go test -count=1 -tags tests ./... -v -cover
	$(COMPOSE_TOOL) down

# ---------------------------------------------------------------------
# -- Kubebuilder related rule
# ---------------------------------------------------------------------
addexamples: ## add examples via kubectl create -f examples/
	cd ./examples/; ls | while read line; do kubectl apply -f $$line; done

manifests: controller-gen ## generate custom resource definitions
	$(LOCALBIN)/controller-gen crd rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
generate: controller-gen ## generate supporting code for custom resource types
	$(LOCALBIN)/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GET_VERSION}

.PHONY: envtest
envtest: ## Download envtest-setup locally if necessary.
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
	${LOCALBIN}/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path

# ---------------------------------------------------------------------
# -- Additional helpers
# ---------------------------------------------------------------------
k3s_mac_lima_create: ## create local k8s using lima
	limactl start --tty=false ./resources/lima/k3s.yaml

k3s_mac_lima_start: ## start local lima k8s
	limactl start k3s

k3s_mac_deploy: build k3s_mac_image ## build image and import image to local lima k8s

k3s_mac_image: ## import built image to local lima k8s
	limactl copy my-image.tar k3s:/tmp/db.tar
	limactl shell k3s sudo k3s ctr images import --all-platforms /tmp/db.tar

