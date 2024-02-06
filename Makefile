BINARY_NAME=ethBalanceProxy
build: clean
	BINARY_NAME=${BINARY_NAME} ./scripts/build.sh

clean:
	BINARY_NAME=${BINARY_NAME} ./scripts/clean.sh

fmt:
	./scripts/fmt.sh

install: build
	./scripts/install.sh

kind\:init:
	./scripts/kind-init.sh

kind\:install:
	./scripts/kind-install.sh

kind\:test:
	./scripts/kind-test.sh

kind\:uninstall:
	./scripts/kind-uninstall.sh

kind\:teardown:
	./scripts/kind-teardown.sh

lint: fmt
	./scripts/lint.sh

local\:start: build
	./${BINARY_NAME} server --upstreams="https://google.com,https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/,https://mainnet.infura.io/v3/e2edc69a0cef4ff28466331d6d972560,https://fittest-falling-smoke.discover.quiknode.pro/"

mock:
	./scripts/mock.sh

test:
	./scripts/test.sh

vet: fmt
	go vet
