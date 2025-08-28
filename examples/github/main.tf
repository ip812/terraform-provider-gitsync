terraform {
  required_providers {
    gitsync = {
      source = "hashicorp.com/ip812/gitsync"
    }
  }
}

provider "gitsync" {
  url = "https://github.com/iypetrov/terraform-provider-gitsync-e2e-test.git"
  token = "some-token"
}

resource "gitsync_values_yaml" "example" {
  branch = "main"
  path    = "values/values.yaml"
  content = <<EOT
foo: bar
replicas: 2
EOT
}
