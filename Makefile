console:
	go build -o ${PWD}/bin/ ${PWD}/cmd/console
dev:
	go run ${PWD}/cmd/console
build:
	go build -o ${PWD}/bin/ ${PWD}/cmd/server
start:
	go run ${PWD}/cmd/server --port 6060 --client ${PWD}/client

