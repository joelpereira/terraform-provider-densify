terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    densify = {
      source = "densify.com/provider/densify"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}
