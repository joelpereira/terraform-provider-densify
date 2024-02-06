# credentials can be passed in as environment variables, DENSIFY_INSTANCE, DENSIFY_USERNAME, DENSIFY_PASSWORD, DENSIFY_TECH_PLATFORM, DENSIFY_ANALYSIS_NAME, DENSIFY_ENTITY_NAME
provider "densify" {
  tech_platform   = "kubernetes"
  cluster         = "<cluster_name>"
  namespace       = "<namespace>"
  controller_type = "deployment"
  pod_name        = "<pod_name>"
  container_name  = "my-container"

  fallback_cpu_req = "1200m"
  fallback_cpu_lim = "4000m"
  fallback_mem_req = "4000Mi"
  fallback_mem_lim = "5120Mi"

  # continue_if_error = true
}

data "densify_container" "optimized" {}

resource "kubernetes_deployment" "den-web" {
  metadata {
    name = "sample-webserver"
    labels = {
      app = "sample-webserver"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "sample-webserver"
      }
    }
    template {
      metadata {
        labels = {
          app = "sample-webserver"
        }
      }
      spec {
        container {
          image = "nginx:latest"
          name  = "my-container"

          port {
            container_port = 80
          }

          resources {
            requests = {
              # original resource settings
              # cpu     = "1200m"
              # memory  = "4000Mi"

              # utilize Densify recommendations instead
              cpu    = data.densify_container.optimized.containers.my-container.rec_cpu_req
              memory = data.densify_container.optimized.containers.my-container.rec_mem_req
            }
            limits = {
              # original resource settings
              # cpu     = "4000m"
              # memory  = "5120Mi"

              # utilize Densify recommendations instead
              cpu    = data.densify_container.optimized.containers.my-container.rec_cpu_lim
              memory = data.densify_container.optimized.containers.my-container.rec_mem_lim
            }
          }


          command = ["sleep"]
          args    = ["infinity"]
        }
      }
    }
  }
}