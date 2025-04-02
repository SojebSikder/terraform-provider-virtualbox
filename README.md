# Description

Terraform virtualbox provider

# Usage

```hcl
terraform {
  required_providers {
    virtualbox = {
      source = "sojebsikder/virtualbox"
      version = "0.1.0"
    }
  }
}

provider "virtualbox" {}

resource "virtualbox_vm" "ubuntu" {
  name     = "UbuntuServer"
  iso_path = "/path/to/ubuntu-server.iso"
  memory   = 2048
  cpus     = 2
}
```

Run terraform

```bash
terraform init
terraform apply
```

## Build

```bash
go build -o terraform-provider-virtualbox
```

## Run (For development)
Create file named `terraform.rc` in this directory `C:\Users\USER\AppData\Roaming`

```bash
provider_installation {

  dev_overrides {
    "hashicorp.com/edu/virtualbox" = "C:/Users/USER/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Run
```
go install .
```

Executable file will be saved on `C:/Users/USER/go/bin`

Create `main.tf` at `examples/provider-install-verification` in project root directory

to test the `main.tf` go to `examples/provider-install-verification` then run
```
terraform plan
```