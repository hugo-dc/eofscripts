all: build install
build: build-eof_gen build-eof_mod build-eof_upd build-mnem2evm build-evm2mnem build-opinfo build-oplist build-eof_dasm build-eof_desc build-eof_fuzz_gen
install: install-eof_gen install-eof_mod install-eof_upd install-mnem2evm install-evm2mnem install-opinfo install-oplist install-eof_dasm install-eof_desc install-eof_fuzz_gen
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
install-eof_gen:
	go install ./cmd/eof_gen
install-eof_mod:
	go install ./cmd/eof_mod
install-eof_upd:
	go install ./cmd/eof_upd
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
