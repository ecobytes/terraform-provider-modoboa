terraform {
  required_providers {
    modoboa = {
      source = "registry.terraform.io/ecobytes/modoboa"
    }
  }
}

provider "modoboa" {
  host = ""
  token = ""
}

data "modoboa_domains" "server" {}

output "server_domains" {
  value = data.modoboa_domains.server.domains
}
