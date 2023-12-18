terraform {
  required_providers {
    densify = {
      source = "densify.com/provider/densify"
    }
  }
}

# Configuration-based authentication
provider "densify" {
  densify_instance = "instance.densify.com:443" # or can be passed in as env variable: DENSIFY_INSTANCE
  username         = "username"                 # or can be passed in as env variable: DENSIFY_USERNAME
  password         = "password"                 # or can be passed in as env variable: DENSIFY_PASSWORD
  tech_platform    = "aws"                      # or can be passed in as env variable: DENSIFY_TECH_PLATFORM
  account_number   = "account-num"
  system_name      = "system-name"
}
