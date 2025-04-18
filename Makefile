build:
	go build -o ./cmd/server/server ./cmd/server/*.go
	go build -o ./cmd/agent/agent ./cmd/agent/*.go

run-autotest: \
	run-autotest-1

run-autotest-1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=./cmd/server/server

run-test:
	go test ./...