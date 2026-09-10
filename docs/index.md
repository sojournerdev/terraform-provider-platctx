---
page_title: "Provider: PlatCtx"
description: |-
  Canonicalize service metadata for platform engineering.
---

# PlatCtx Provider

The PlatCtx provider canonicalizes service metadata for platform engineering.

## Goals

In enterprise environments with many teams, cloud providers, and modules, service metadata like identity, ownership, environment, and compliance facts is required by every resource, but each provider expects it in a different format. When you move workloads across providers or environments, that context breaks.

PlatCtx centralizes this metadata into a single source of truth. The `canonicalize` function validates and normalizes it into a consistent typed structure that passes cleanly downstream to modules, tooling, and service catalogs like Backstage.

If your environment only uses one provider and a handful of modules, you likely don't need this. PlatCtx is built for scale, where visibility and auditing across providers and teams is the hard problem.

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
