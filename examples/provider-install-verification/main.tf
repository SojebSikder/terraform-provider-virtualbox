terraform {
  required_providers {
    virtualbox = {
      source = "hashicorp.com/edu/virtualbox"
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