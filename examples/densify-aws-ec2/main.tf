terraform {
  required_providers {
    densify = {
      source = "densify.com/provider/densify"
    }
  }
  # required_version = ">= 1.1.0"
}

# credentials can be passed in as environment variables, DENSIFY_INSTANCE, DENSIFY_USERNAME, DENSIFY_PASSWORD, DENSIFY_TECH_PLATFORM, DENSIFY_ANALYSIS_NAME, DENSIFY_ENTITY_NAME

provider "densify" {
  tech_platform = "aws"
  # account_name = "account-name"
  account_number = "777999222999"
  system_name    = "ea-pro-duct-259"
  # skip_errors = false
}

data "densify_cloudreco" "sample" {}

output "sample_cloudreco" {
  value = data.densify_cloudreco.sample
}
