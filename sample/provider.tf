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
  unique_id = "47"
  name = "TMNT3"
  rating = "600"
}
