app:
	go build -o ${PWD}/bin/ ${PWD}/cmd/app
run_app:
	go build -o ${PWD}/bin/ ${PWD}/cmd/app
	${PWD}/bin/app
