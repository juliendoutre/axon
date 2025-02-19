terraform {
  required_providers {
    keycloak = {
      source  = "linz/keycloak"
      version = "4.4.1"
    }
  }
}

variable "keycloak_http_port" {
  type = number
}

variable "keycloak_admin_password" {
  type      = string
  sensitive = true
}

provider "keycloak" {
  url       = "http://localhost:${var.keycloak_http_port}"
  client_id = "admin-cli"
  username  = "admin"
  password  = var.keycloak_admin_password
}

data "keycloak_realm" "master" {
  realm = "master"
}

resource "keycloak_openid_client" "openid_client" {
  realm_id                     = data.keycloak_realm.master.id
  client_id                    = "axon"
  name                         = "axon"
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  direct_access_grants_enabled = true
  valid_redirect_uris          = ["/*"]
}

output "openid_client_secret" {
  value       = keycloak_openid_client.openid_client.client_secret
  description = "secret"
  sensitive   = true
}
