terraform {
  required_providers {
    example = {
      source  = "localhost/providers/example"
      version = "~> 0.0.2"
    }
  }
}

provider "example" {

}

resource "example_item" "item" {
  unique_id = "46"
  name = "TMNT"
  rating = "500"
}
