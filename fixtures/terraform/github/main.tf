# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    gitsync = {
      source = "hashicorp.com/ip812/gitsync"
    }
  }
}

provider "gitsync" {
  url   = "https://github.com/iypetrov/terraform-provider-gitsync-e2e-test.git"
}

resource "gitsync_values_yaml" "example_yaml" {
  branch  = "main"
  path    = "values/values.yaml"
  content = <<EOT
name: bar
replicas: 2
EOT
}

resource "gitsync_values_json" "example_json" {
  branch  = "main"
  path    = "values/values.json"
  content = <<EOT
{
  "name": "bar",
  "replicas": 2
}
EOT
}

resource "gitsync_values_file" "example_file" {
  branch  = "main"
  path    = "values/values.md"
  content = <<EOT
# Title

Name is bar and replicas are 2
EOT
}
