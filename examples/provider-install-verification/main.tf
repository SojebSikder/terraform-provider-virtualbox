terraform {
  required_providers {
    virtualbox = {
      source  = "hashicorp.com/edu/virtualbox"
      version = "0.0.2"
    }
  }
}

provider "virtualbox" {}

resource "virtualbox_vm" "ubuntu" {
  count    = 2
  name     = "UbuntuServerTest${count.index}"
  iso_path = "D:/All Programs/Software/OS/ubuntu-24.04.2-live-server-amd64.iso"
  memory   = 2048
  cpus     = 2

  network_adapter = {
    type = "nat"
  }
}
