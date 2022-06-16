all: build-idcost build-returnevm build-deploy build-call build-eof_gen build-yulreturn
build-idcost:
	go build -o ./build/ ./cmd/idcost
build-returnevm:
	go build -o ./build/ ./cmd/returnevm
build-deploy:
	go build -o ./build/ ./cmd/deploy
build-call:
	go build -o ./build/ ./cmd/call
build-eof_gen:
	go build -o ./build/ ./cmd/eof_gen
build-yulreturn:
	go build -o ./build/ ./cmd/yulreturn
