terraform {
  required_providers {
    platctx = {
      source  = "registry.terraform.io/sojournerdev/platctx"
      version = "~> 0.1.0"
    }
  }
}

output "context" {
  value = provider::platctx::canonicalize({
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
