all: build install
build: build-idcost build-returnevm build-returneofdata build-deploy build-call build-eof_gen build-eof_mod build-eof_upd build-yulreturn build-gentruncpush build-createaddress build-createaddress2 build-mnem2evm build-evm2mnem build-opinfo build-oplist build-eof_dasm
build-idcost:
	go build -o ./build/ ./cmd/idcost
build-returnevm:
	go build -o ./build/ ./cmd/returnevm
build-returneofdata:
	go build -o ./build/ ./cmd/returneofdata
build-deploy:
	go build -o ./build/ ./cmd/deploy
build-call:
	go build -o ./build/ ./cmd/call
build-eof_gen:
	go build -o ./build/ ./cmd/eof_gen
build-eof_mod:
	go build -o ./build/ ./cmd/eof_mod
build-eof_upd:
	go build -o ./build/ ./cmd/eof_upd
build-yulreturn:
	go build -o ./build/ ./cmd/yulreturn
build-gentruncpush:
	go build -o ./build/ ./cmd/gentruncpush
build-createaddress:
	go build -o ./build/ ./cmd/create_address
build-createaddress2:
	go build -o ./build/ ./cmd/create_address2
build-mnem2evm:
	go build -o ./build/ ./cmd/mnem2evm
build-evm2mnem:
	go build -o ./build/ ./cmd/evm2mnem
build-opinfo:
	go build -o ./build/ ./cmd/opinfo
build-oplist:
	go build -o ./build/ ./cmd/oplist
build-eof_dasm:
	go build -o ./build/ ./cmd/eof_dasm
install:
	mv ./build/* ~/.bin/
test:
	cd common && go test
