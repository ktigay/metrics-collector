SERVER_PORT=12345

build:
	go build -o ./cmd/server/server ./cmd/server/*.go
	go build -o ./cmd/agent/agent ./cmd/agent/*.go

run-test: \
	build \
	run-test-a \
	run-test-u \
	run-test-s

run-test-u:
	go test ./...

run-test-s:
	go vet -vettool=$$(which statictest) ./...

run-test-a: \
	run-test-a1 \
	run-test-a2 \
	run-test-a3 \
	run-test-a4 \
	run-test-a5

run-test-a1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=./cmd/server/server
run-test-a2:
	metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=./cmd/agent/agent
run-test-a3:
	metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=. -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server
run-test-a4:
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -server-port=$(SERVER_PORT) -source-path=.
run-test-a5:
	metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=$(SERVER_PORT) -source-path=.

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
	git fetch template && git checkout template/main .github