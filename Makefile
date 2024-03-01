all: build install
build: build-idcost build-returnevm build-deploy build-call build-eof_gen build-eof_mod build-eof_upd build-yulreturn build-gentruncpush build-createaddress build-createaddress2 build-mnem2evm build-evm2mnem build-opinfo build-oplist build-eof_dasm build-eof_desc build-eof_fuzz_gen
install: install-idcost install-returnevm install-deploy install-call install-eof_gen install-eof_mod install-eof_upd install-yulreturn install-gentruncpush install-createaddress install-createaddress2 install-mnem2evm install-evm2mnem install-opinfo install-oplist install-eof_dasm install-eof_desc install-eof_fuzz_gen
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
build-eof_desc:
	go build -o ./build/ ./cmd/eof_desc
build-eof_fuzz_gen:
	go build -o ./build/ ./cmd/eof_fuzz_gen
install-idcost:
	go install ./cmd/idcost
install-returnevm:
	go install ./cmd/returnevm
install-deploy:
	go install ./cmd/deploy
install-call:
	go install ./cmd/call
install-eof_gen:
	go install ./cmd/eof_gen
install-eof_mod:
	go install ./cmd/eof_mod
install-eof_upd:
	go install ./cmd/eof_upd
install-yulreturn:
	go install ./cmd/yulreturn
install-gentruncpush:
	go install ./cmd/gentruncpush
install-createaddress:
	go install ./cmd/create_address
install-createaddress2:
	go install ./cmd/create_address2
install-mnem2evm:
	go install ./cmd/mnem2evm
install-evm2mnem:
	go install ./cmd/evm2mnem
install-opinfo:
	go install ./cmd/opinfo
install-oplist:
	go install ./cmd/oplist
install-eof_dasm:
	go install ./cmd/eof_dasm
install-eof_desc:
	go install ./cmd/eof_desc
install-eof_fuzz_gen:
	go install ./cmd/eof_fuzz_gen

test:
	cd common && go test
