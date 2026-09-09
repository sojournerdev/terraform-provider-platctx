terraform {
  required_version = ">= 1.8"

  required_providers {
    platctx = {
      source = "registry.terraform.io/sojournerdev/platctx"
    }
  }
}

locals {
  context = provider::platctx::canonicalize({
    identity = {
      name      = "payments"
      namespace = "platform"
    }
    ownership = {
      owned_by    = "team:platform"
      operated_by = "team:sre"
    }
    environment = "production"
    governance = {
      criticality         = "critical"
      cost_center         = "CC-FIN-001"
      data_classification = "internal"
    }
  })
}

output "context" {
  value = local.context
}
