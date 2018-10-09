.PHONY: all build prepare deps test generate

prepare:
	@echo "Installing ginkgo"
	go get github.com/onsi/ginkgo
	dep ensure

deps:
	@echo "Setting up the vendors folder..."
	@dep ensure -v
	@echo ""
	@echo "Resolved dependencies:"
	@dep status
	@echo ""

test:
	@echo "Running unit and integration test"
	docker-compose down
	docker-compose up -d
	MONGO_URL=mongodb://localhost:27017 MONGO_DATABASE=url-shortener MONGO_COLLECTION=urls ginkgo -r

build:
	docker build -t url-shortener-service .

# https://blog.codecentric.de/en/2017/08/gomock-tutorial/
generate:
	@echo "Generating files such as mock files"
	go generate ./...
