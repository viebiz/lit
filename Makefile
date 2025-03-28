# Variables
ENV:=dev
GROUP_NAME:=viebiz
PROJECT_NAME:=lit

# Shorten cmd
DOCKER_BUILD_BIN := docker
COMPOSE_BIN := ENV=$(ENV) GROUP_NAME=$(GROUP_NAME) PROJECT_NAME=$(PROJECT_NAME) docker compose
COMPOSE_TOOL_RUN := $(COMPOSE_BIN) run --rm --service-ports tool

mod:
	go mod tidy
	go mod vendor

init: kafka pg redis
	echo "Start Kafka, Postgres, and Redis!"

pg:
	@$(COMPOSE_BIN) up postgres -d

redis:
	@$(COMPOSE_BIN) up redis -d

kafka:
	@$(COMPOSE_BIN) up kafka -d

jaeger:
	@$(COMPOSE_BIN) up jaeger -d

test:
	@$(COMPOSE_TOOL_RUN) sh -c "go test -mod=vendor -vet=all -coverprofile=coverage.out -failfast -timeout 5m ./..."

kafka-topics: kafka
	@$(COMPOSE_BIN) run --rm kafka-topics sh -c "kafka-topics --create --topic $$TOPIC_TEST_1 --partitions 2 --replication-factor 1 --bootstrap-server kafka:9092"
	@$(COMPOSE_BIN) run --rm kafka-topics sh -c "kafka-topics --create --topic $$TOPIC_TEST_2 --partitions 2 --replication-factor 1 --bootstrap-server kafka:9092"

benchmark:
	@$(COMPOSE_TOOL_RUN) sh -c "go test ./... -bench=. -run=^#"

gen-mocks:
	@$(COMPOSE_TOOL_RUN) sh -c "mockery"

gen-proto:
	@$(COMPOSE_TOOL_RUN) sh -c "buf generate"

gen-pem:
	@$(COMPOSE_TOOL_RUN) sh -c "openssl genrsa -out ./crypto/rsa/testdata/pkcs1_rsa_private_key.pem 2048"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl rsa -pubout -in ./crypto/rsa/testdata/pkcs1_rsa_private_key.pem -outform PEM -RSAPublicKey_out -out ./crypto/rsa/testdata/pkcs1_rsa_public_key.pem"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out ./crypto/rsa/testdata/ecdsa_private_key.pem"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl ec -pubout -in ./crypto/rsa/testdata/ecdsa_private_key.pem -out ./crypto/rsa/testdata/ecdsa_public_key.pem"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl genpkey -algorithm RSA -out ./crypto/rsa/testdata/test_rsa_private_key.pem -pkeyopt rsa_keygen_bits:2048"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl rsa -pubout -in ./crypto/rsa/testdata/test_rsa_private_key.pem -out ./crypto/rsa/testdata/test_rsa_public_key.pem"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl genpkey -algorithm RSA -out ./crypto/rsa/testdata/rsa_private_key.pem -pkeyopt rsa_keygen_bits:2048"
	@$(COMPOSE_TOOL_RUN) sh -c "openssl rsa -pubout -in ./crypto/rsa/testdata/rsa_private_key.pem -out ./crypto/rsa/testdata/rsa_public_key.pem"

tear-down:
	@$(COMPOSE_BIN) down -v
