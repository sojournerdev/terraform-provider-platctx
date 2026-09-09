---
page_title: "Provider: PlatCtx"
description: |-
  Canonicalize service metadata for platform engineering.
---

# PlatCtx Provider

The PlatCtx provider canonicalizes service metadata for platform engineering.

## Why Use This

Service metadata like identity, ownership, environment, and compliance facts often lives in different places depending on the cloud provider. When you move workloads across providers or environments, that context breaks.

PlatCtx gives you provider-agnostic service metadata that travels with your resources. The `canonicalize` function validates and normalizes this metadata into a consistent typed structure so downstream tools and modules can rely on it.

This metadata powers service catalogs like Backstage, where it drives service discovery, ownership tracking, and compliance reporting.

## Example Usage

```terraform
terraform {
  required_providers {
    platctx = {
      source = "registry.terraform.io/sojournerdev/platctx"
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
      owned_by = "team:platform"
    }
  })
}
```

## Argument Reference

This provider has no configuration arguments.
