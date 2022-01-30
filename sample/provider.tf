terraform {
  required_providers {
    planner = {
      source  = "terraform-example.com/exampleprovider/example"
      version = "~> 1.0.0"
    }
  }
}

provider "planner" {}

resource "planner_todo_note" "item" {
  title    = "my cool title"
  message  = "my cool message"
  priority = 8
}
