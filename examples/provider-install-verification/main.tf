terraform {
  required_providers {
    gitoutput = {
      source = "hashicorp.com/iypetrov/gitoutput"
    }
  }
}

provider "gitoutput" {
  token = "foo"
  owner = "iypetrov"
}

resource "gitoutput_values_file" "example" {
  file_path = "coffees.json"
  values    = "foo: bar"
}

resource "gitoutput_values_file" "example-1" {
  file_path = "beer.json"
  values    = "foo: buzz"
}
