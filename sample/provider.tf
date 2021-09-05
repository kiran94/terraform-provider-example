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

resource "example_tv_show" "item" {
  unique_id = 50
  name = "TMNT"
  rating = 9
}
