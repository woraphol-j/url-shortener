.PHONY: all deps test build benchmark coveralls mockgen

DEP_VERSION=0.4.1
OS := $(shell uname | tr '[:upper:]' '[:lower:]')

prepare:
	@echo "Installing dep..."
	@curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-${OS}-amd64 -o ${GOPATH}/bin/dep
	@chmod a+x ${GOPATH}/bin/dep

	@echo "Installing goconvey"
	go get github.com/smartystreets/goconvey

	@echo "Installing richgo"
	go get -u github.com/kyoh86/richgo

	@echo "Installing realize"
	go get github.com/oxequa/realize

deps:
	@echo "Setting up the vendors folder..."
	@dep ensure -v
	@echo ""
	@echo "Resolved dependencies:"
	@dep status
	@echo ""

serve:
	@echo "Running bp-message-server"
	${GOPATH}/bin/realize start --open --server --build --run --no-config

cover:
	@echo "Running coverage"
	${GOPATH}/bin/goconvey

test:
	@echo "Running test"
	${GOPATH}/bin/richgo test ./... -v

# https://blog.codecentric.de/en/2017/08/gomock-tutorial/
generate:
	@echo "Generating files such as mock files"
	go generate ./...
