terraform {
  required_providers {
    git-output = {
      source = "hashicorp.com/iypetrov/ip812-git-output"
    }
  }
}

provider "ip812-git-output" {
    token = "foo"
    owner = "iypetrov"
}

resource "ip812_git_output_values_file" "example" {
    file_path = "coffees.json"
    values    = "foo: bar"
}
