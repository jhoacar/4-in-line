app:
	go build -o ${PWD}/bin/ ${PWD}/cmd/console
run_app:
	make app
	${PWD}/bin/app

