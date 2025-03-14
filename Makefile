# Variables
ENV:=dev
GROUP_NAME:=viebiz
PROJECT_NAME:=lit

# Shorten cmd
DOCKER_BUILD_BIN := docker
COMPOSE_BIN := ENV=$(ENV) GROUP_NAME=$(GROUP_NAME) PROJECT_NAME=$(PROJECT_NAME) docker compose
COMPOSE_TOOL_RUN := $(COMPOSE_BIN) run --rm --service-ports tool

init: pg redis collector
	echo "Start Postgres, Redis!"
pg:
	@$(COMPOSE_BIN) up postgres -d

redis:
	@$(COMPOSE_BIN) up redis -d

collector:
	@$(COMPOSE_BIN) up collector -d

test:
	@$(COMPOSE_TOOL_RUN) sh -c "go test -mod=vendor -vet=all -coverprofile=coverage.out -failfast -timeout 5m ./..."

benchmark:
	@$(COMPOSE_TOOL_RUN) sh -c "go test ./... -bench=. -run=^#"

gen-mocks:
	@$(COMPOSE_TOOL_RUN) sh -c "mockery"

gen-proto:
	@$(COMPOSE_TOOL_RUN) sh -c "buf generate"
