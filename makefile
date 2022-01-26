
PROVIDER_NAME=example
PROVIDER_VERSION=0.0.2
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

INSTALL_PATH=~/.local/share/terraform/plugins/localhost/providers/$(PROVIDER_NAME)/$(PROVIDER_VERSION)/$(OS)_$(ARCH)

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
