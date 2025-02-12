# Variables
ENV:=dev
GROUP_NAME:=bizgroup
PROJECT_NAME:=lit

# Shorten cmd
DOCKER_BUILD_BIN := docker
COMPOSE_BIN := ENV=$(ENV) GROUP_NAME=$(GROUP_NAME) PROJECT_NAME=$(PROJECT_NAME) docker compose
COMPOSE_TOOL_RUN := $(COMPOSE_BIN) run --rm --service-ports tool

pull:
	@$(COMPOSE_BIN) pull || true

test:
	@$(COMPOSE_TOOL_RUN) sh -c "go test -mod=vendor -coverprofile=coverage.out -failfast -timeout 5m ./..."

benchmark:
	@$(COMPOSE_TOOL_RUN) sh -c "go test ./... -bench=. -run=^#"

gen-mocks:
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir env --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir httpclient --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir httpserv --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir iam --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir jwt --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir postgres --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir monitoring --all --recursive --inpackage"
	@$(COMPOSE_TOOL_RUN) sh -c "mockery --dir vault --all --recursive --inpackage"

gen-proto:
	@$(COMPOSE_TOOL_RUN) sh -c "buf generate"
