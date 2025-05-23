PROJECT_NAME=metrics-collector

AGENT_PATH=./cmd/agent/agent
SERVER_PATH=./cmd/server/server
TEMP_FILE=/tmp/metric_storage.txt
TEST_SERVER_PORT=\$$(random unused-port)

SHELL := /bin/bash
CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)

# локальный, внешний порт.
LOCAL_PORT=4001
# внутренний порт сервера.
SRV_PORT=8080
# адрес, который слушает сервер.
SRV_LISTEN=:$(SRV_PORT)
# адрес, по которому стучится агент.
SRV_ADDR=server:$(SRV_PORT)
# dsn к postgres
DATABASE_DSN=postgres://postgres:postgres@192.168.1.46:5430/praktikum?sslmode=disable

COMPOSE := export PROJECT_NAME=$(PROJECT_NAME) CURRENT_UID=$(CURRENT_UID) \
 		   CURRENT_GID=$(CURRENT_GID) SRV_LISTEN=$(SRV_LISTEN) SRV_PORT=$(SRV_PORT) \
 		   LOCAL_PORT=$(LOCAL_PORT) SRV_ADDR=$(SRV_ADDR) DATABASE_DSN=$(DATABASE_DSN) && cd docker &&

DOCKER_RUN := cd docker && docker run --rm -v ${PWD}:/app -it $(PROJECT_NAME)-app

build-local:
	go build -o $(SERVER_PATH) ./cmd/server/*.go
	go build -o $(AGENT_PATH) ./cmd/agent/*.go

build:
	$(COMPOSE) docker compose -f docker-compose.build.yml build app

go-build-server:
	cd docker && docker run --rm -v ${PWD}:/app -it $(PROJECT_NAME)-app \
	go build -gcflags "all=-N -l" -o /app/cmd/server/server -tags dynamic /app/cmd/server/

go-build-agent:
	cd docker && docker run --rm -v ${PWD}:/app -it $(PROJECT_NAME)-app \
	go build -gcflags "all=-N -l" -o /app/cmd/agent/agent -tags dynamic /app/cmd/agent/

run-test: \
	go-build-server \
	go-build-agent \
	run-test-a \
	run-test-u \
	run-test-s \
	run-lint \

run-test-u:
	$(DOCKER_RUN) sh -c "go test ./..."

run-test-s:
	$(DOCKER_RUN) sh -c "go vet -vettool=\$$(which statictest) ./..."

run-lint:
	$(DOCKER_RUN) golangci-lint run

cs-fix:
	$(DOCKER_RUN) gofumpt -w -extra internal/ cmd/

run-test-a: \
	run-test-a1 \
	run-test-a2 \
	run-test-a3 \
	run-test-a4 \
	run-test-a5 \
	run-test-a6 \
	run-test-a7 \
	run-test-a8 \
	run-test-a9 \
	run-test-a10 \
	run-test-a11 \
	run-test-a12 \

define test_cmd
	$(DOCKER_RUN)  sh -c "metricstest -test.v -test.run=^TestIteration$(1)$$ -source-path=. -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH)$(2)"
endef

run-test-a1:
	$(call test_cmd,1)
run-test-a2:
	$(call test_cmd,2[AB])
run-test-a3:
	$(call test_cmd,3[AB]*)
run-test-a4:
	$(call test_cmd,4, -server-port=$(TEST_SERVER_PORT))
run-test-a5:
	$(call test_cmd,5, -server-port=$(TEST_SERVER_PORT))
run-test-a6:
	$(call test_cmd,6, -server-port=$(TEST_SERVER_PORT))
run-test-a7:
	$(call test_cmd,7, -server-port=$(TEST_SERVER_PORT))
run-test-a8:
	$(call test_cmd,8, -server-port=$(TEST_SERVER_PORT))
run-test-a9:
	$(call test_cmd,9, -server-port=$(TEST_SERVER_PORT) -file-storage-path=$(TEMP_FILE))
run-test-a10:
	$(call test_cmd,10[AB], -server-port=$(TEST_SERVER_PORT) -file-storage-path=$(TEMP_FILE) -database-dsn='$(DATABASE_DSN)')
run-test-a11:
	$(call test_cmd,11, -server-port=$(TEST_SERVER_PORT) -file-storage-path=$(TEMP_FILE) -database-dsn='$(DATABASE_DSN)')
run-test-a12:
	$(call test_cmd,12, -server-port=$(TEST_SERVER_PORT) -file-storage-path=$(TEMP_FILE) -database-dsn='$(DATABASE_DSN)')

up: \
	up-db \
	up-server \
	up-agent \

up-server: \
	go-build-server
	$(COMPOSE) docker compose up -d server

up-agent: \
	go-build-agent
	$(COMPOSE) docker compose up -d agent

up-db:
	$(COMPOSE) docker compose up -d postgres

down:
	$(COMPOSE) docker compose down server agent postgres

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
	git fetch template && git checkout template/main .github