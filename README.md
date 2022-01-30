# terraform-provider-example

[![main](https://github.com/kiran94/terraform-provider-example/actions/workflows/main.yml/badge.svg)](https://github.com/kiran94/terraform-provider-example/actions/workflows/main.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/kiran94/terraform-provider-example.svg)]

## Local Development

At the root of this repo:

```sh
# Start an Example Api
# Do this in a different terminal
make example_server

# Build the provider
make dev

# Run the Sample
cd sample
terraform init
terraform apply
```

When making changes you want to clear the state in the sample. To do this use `make clean`
