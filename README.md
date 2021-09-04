# terraform-provider-example

## Local Development

At the root of this repo:

```sh
make dev
cd sample

terraform init
terraform plan
terraform apply
```

Note:

- If you already have a `.terraform.lock.hcl` then it will need to be removed during development as you make changes to the provider.
- This sample assumes you have the sample API running + Redis instance
