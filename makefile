PROVIDER_NAME=example
PROVIDER_VERSION=1.0.0
PROVIDER_HOST=terraform-example.com
PROVIDER_NAMESPACE=exampleprovider
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

INSTALL_PATH=~/.terraform.d/plugins/$(PROVIDER_HOST)/$(PROVIDER_NAMESPACE)/$(PROVIDER_NAME)/$(PROVIDER_VERSION)/$(OS)_$(ARCH)/

dev:
	mkdir -p $(INSTALL_PATH)
	go build -o $(INSTALL_PATH)/terraform-provider-$(PROVIDER_NAME) main.go

build:
	go build

test:
	go test

clean:
	rm ./sample/.terraform.lock.hcl
	rm ./sample/terraform.tfstate

example_server:
	go run ./cmd/server/main.go
