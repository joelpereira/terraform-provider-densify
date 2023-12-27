terraform {
  required_providers {
    densify = {
      source = "densify.com/provider/densify"
    }
  }
}

provider "densify" {
  tech_platform  = "aws" # or can be passed in as env variable: DENSIFY_TECH_PLATFORM
  account_number = var.account_number
  system_name    = var.name

  # continue_if_error = true
  fallback_instance_type = "m4.large" // backup/fallback instance type until there is a recommendation
}
data "densify_cloud" "optimization" {}

provider "aws" {
  region = "us-east-2"
}

resource "aws_instance" "create" {

  # legacy way of creating an instance;  hardcoding the instance type
  # instance_type = "m4.large"

  # new self-optimizing instance type from Densify
  instance_type = data.densify_cloud.optimization.recommended_type


  ami = "ami-00eeedc4036573771" # Ubuntu 22.04 LTS

  # tag instance to make it Self-Aware these tags are optional and can set as few or as many as you like.
  tags = {
    Name = var.name
    #Should match the densify_unique_id value as this is how Densify references the system as unique
    "Provisioning ID"                 = var.name
    Densify-optimal-instance-type     = data.densify_cloud.optimization.recommended_type
    Densify-potential-monthly-savings = format("%s effort and estimated savings of $%f", data.densify_cloud.optimization.effort_estimate, data.densify_cloud.optimization.savings_estimate)
  }
}

