resource "komodo_swarm_config" "example" {
  swarm = "my-swarm"
  name  = "app-config"
  data  = "replace-with-config-content"
}

# Optional labels and template driver
resource "komodo_swarm_config" "with_options" {
  swarm           = "my-swarm"
  name            = "nginx-conf"
  data            = "replace-with-nginx-config"
  template_driver = ""
  labels          = ["environment=production", "managed-by=terraform"]
}
