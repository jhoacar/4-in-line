app:
	go build -o ${PWD}/bin/ ${PWD}/cmd/console
run_app:
	make app
	${PWD}/bin/app
server:
	go build -o ${PWD}/bin/ ${PWD}/cmd/server
run_server:
	make server
	${PWD}/bin/server

