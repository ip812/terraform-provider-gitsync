terraform {
  required_providers {
    gitsync = {
      source = "hashicorp.com/ip812/gitsync"
    }
  }
}

provider "gitsync" {
  alias = "github"
  repository = "https://github.com/iypetrov/terraform-provider-gitsync-e2e-test.git"
}
