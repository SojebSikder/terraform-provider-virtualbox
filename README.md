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