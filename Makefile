BINARY_NAME=eth_balance_proxy
build: clean dep
	BINARY_NAME=${BINARY_NAME} ./scripts/build.sh

clean:
	BINARY_NAME=${BINARY_NAME} ./scripts/clean.sh

dep:
	go mod download

lint: fmt
	./scripts/lint.sh

fmt:
	./scripts/fmt.sh

mock:
	./scripts/mock.sh

start_server: build
	./${BINARY_NAME} server --endpoints="https://www.google.com,https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/,https://mainnet.infura.io/v3/e2edc69a0cef4ff28466331d6d972560,https://fittest-falling-smoke.discover.quiknode.pro/"

test:
	./scripts/test.sh

vet: fmt
	go vet
