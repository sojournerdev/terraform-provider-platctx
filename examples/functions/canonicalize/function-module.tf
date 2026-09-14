terraform {
  required_providers {
    platctx = {
      source  = "registry.terraform.io/sojournerdev/platctx"
      version = "~> 0.1.0"
    }
  }
}

locals {
  ctx = provider::platctx::canonicalize({
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

module "service" {
  source = "./modules/platform/service"

  name                = local.ctx.identity.name
  namespace           = local.ctx.identity.namespace
  owner               = local.ctx.ownership.owned_by
  operator            = local.ctx.ownership.operated_by
  environment         = local.ctx.environment
  criticality         = local.ctx.governance.criticality
  cost_center         = local.ctx.governance.cost_center
  data_classification = local.ctx.governance.data_classification
}
