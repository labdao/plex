# This makefile is for bringing up local stack using minikube


.PHONY: all
all: checks unit-test

.PHONY: checks
checks: license lint

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh

.PHONY: unit-test
unit-test:
	@scripts/check_unit.sh

.PHONY: sandbox-cli-docker
sandbox-cli-docker:
	@echo "Building sandbox-cli docker image"
	@docker build -f ./images/sandbox-cli/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(SANDBOX_CLI_IMAGE_NAME):latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) .
