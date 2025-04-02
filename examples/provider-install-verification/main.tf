terraform {
  required_providers {
    virtualbox = {
      source = "hashicorp.com/edu/virtualbox"
      version = "0.0.1"
    }
  }
}

provider "virtualbox" {}

resource "virtualbox_vm" "ubuntu" {
  name     = "UbuntuServerTest"
  iso_path = "D:/All Programs/Software/OS/ubuntu-24.04.2-live-server-amd64.iso"
  memory   = 2048
  cpus     = 2
}