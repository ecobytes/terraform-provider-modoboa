# terraform-provider-modoboa

*Experimental* Terraform provider for Modoboa.

## Build

Generate code

```sh
go generate ./...
sed 's/`json:/`tfsdk:/g' -i internal/client/client.gen.go
sed 's/,omitempty//g' -i internal/client/client.gen.go
sed 's/\*time\.Time/\*string/g' -i internal/client/client.gen.go
```

Install binary to system

```sh
go mod tidy
go install .
```

## Run

The binary is a plugin that can be installed into Terraform.

For a local installation prepare `~/.terraformrc`:

```hcl
provider_installation {

  dev_overrides {
    # Example GOBIN path, will need to be replaced with your own GOBIN path. Default is $GOPATH/bin
    "registry.terraform.io/ecobytes/modoboa" = "/home/$USER/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Replace `$USER` with the name of your home directory.

```sh
cd internal/test
TF_LOG=trace terraform plan
TF_LOG=debug terraform apply
terraform state list
```
