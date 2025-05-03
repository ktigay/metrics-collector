SERVER_PORT=12345
AGENT_PATH=./cmd/agent/agent
SERVER_PATH=./cmd/server/server

build:
	go build -o $(SERVER_PATH) ./cmd/server/*.go
	go build -o $(AGENT_PATH) ./cmd/agent/*.go

run-test: \
	build \
	run-test-a \
	run-test-u \
	run-test-s \
	run-lint \

run-test-u:
	go test ./...

run-test-s:
	go vet -vettool=$$(which statictest) ./...

run-test-a: \
	run-test-a1 \
	run-test-a2 \
	run-test-a3 \
	run-test-a4 \
	run-test-a5 \
	run-test-a6 \
	run-test-a7 \
	run-test-a8 \

run-test-a1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=$(SERVER_PATH)
run-test-a2:
	metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=$(AGENT_PATH)
run-test-a3:
	metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=. -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH)
run-test-a4:
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH) -server-port=$(SERVER_PORT) -source-path=.
run-test-a5:
	metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH) -server-port=$(SERVER_PORT) -source-path=.
run-test-a6:
	metricstest -test.v -test.run=^TestIteration6$$ -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH) -server-port=$(SERVER_PORT) -source-path=.
run-test-a7:
	metricstest -test.v -test.run=^TestIteration7$$ -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH) -server-port=$(SERVER_PORT) -source-path=.
run-test-a8:
	export LOG_LEVEL="error" && metricstest -test.v -test.run=^TestIteration8$$ -agent-binary-path=$(AGENT_PATH) -binary-path=$(SERVER_PATH) -server-port=$(SERVER_PORT) -source-path=.

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
	git fetch template && git checkout template/main .github

run-lint:
	golangci-lint run