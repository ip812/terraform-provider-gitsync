# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    gitsync = {
      source = "hashicorp.com/ip812/gitsync"
    }
  }
}

provider "gitsync" {
  url   = "https://gitlab.com/iypetrov/terraform-provider-gitsync-e2e-test.git"
  token = "your-token"
}

resource "gitsync_values_yaml" "example" {
  branch  = "main"
  path    = "values/values.yaml"
  content = <<EOT
name: bar
replicas: 2
EOT
}
