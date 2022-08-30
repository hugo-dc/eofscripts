all: build-idcost build-returnevm build-deploy build-call build-eof_gen build-yulreturn build-gentruncpush build-createaddress build-createaddress2
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
build-gentruncpush:
	go build -o ./build/ ./cmd/gentruncpush
build-createaddress:
	go build -o ./build/ ./cmd/create_address
build-createaddress2:
	go build -o ./build/ ./cmd/create_address2
