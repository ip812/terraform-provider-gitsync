terraform {
  required_providers {
    gitsync = {
      source = "hashicorp.com/ip812/gitsync"
    }
  }
}

provider "gitsync" {
  token      = "foo"
  repository = "https://github.com/ip812/terraform-provider-gitsync.git"
}

resource "gitsync_values_yaml" "example" {
  path    = "coffees.yaml"
  content = "foo: bar"
}

resource "gitsync_values_yaml" "example-1" {
  path    = "beer.yaml"
  content = "foo: buzz"
}
