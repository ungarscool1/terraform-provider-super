terraform {
  required_providers {
    super = {
      source = "ungarscool1/super-terraform-provider"
      version = "1.0.0"
    }
  }
}

provider "super" {}

data "super_api_gateway_api_keys" "example" {}

output "api_keys" {
  value = data.super_api_gateway_api_keys.example
}
