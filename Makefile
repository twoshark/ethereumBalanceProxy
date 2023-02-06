BINARY_NAME=eth_balance_proxy
build:
	BINARY_NAME=${BINARY_NAME} ./scripts/build.sh

clean:
	BINARY_NAME=${BINARY_NAME} ./scripts/clean.sh

dep:
	go mod download

lint:
	./scripts/lint.sh

mock:
	./scripts/mock.sh

start_server: build
	./${BINARY_NAME} server

test:
	./scripts/test.sh

vet:
	go vet
