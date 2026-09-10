# Terraform Provider PlatCtx

## Overview

Here is a Terraform provider for an interesting problem I ran into. I was building a monorepo for an internal platform that contained modules and shared-infrastructure used as a runtime. I was inspired by [Cloud Posse's context pattern](https://docs.cloudposse.com/modules/library/provider/context/) and was using the concept as part of the platform to help with centralizing metadata that resources should have. The problem I saw was that I had context specific to the platform for an enterprise environment, and needed that data to be agnostic, but different providers need different shapes (azurerm needs tags, AWS needs different fields, etc). The same data from context could be used rather than duplicating.

Thus comes this provider solution. It is opinionated about some fields, but it's intended to simplify canonicalizing metadata that belongs to a platform. Basically we can rely on being programmatic with data being strongly typed, thus making the user interface of HCL code simpler.

## Goals

In enterprise environments with many teams, cloud providers, and modules, service metadata like identity, ownership, environment, and compliance facts is required by every resource, but each provider expects it in a different format. When you move workloads across providers or environments, that context breaks.

PlatCtx centralizes this metadata into a single source of truth. The `canonicalize` function validates and normalizes it into a consistent typed structure that passes cleanly downstream to modules, tooling, and service catalogs like Backstage.

If your environment only uses one provider and a handful of modules, you likely don't need this. PlatCtx is built for scale, where visibility and auditing across providers and teams is the hard problem.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.27 (for development)

## Usage

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
  # Returns provider-agnostic context with normalized
  # identity, ownership, environment, and governance
}
```

## Documentation

See the [Terraform Registry documentation](https://registry.terraform.io/providers/sojournerdev/platctx/latest/docs) for function reference, examples, and configuration details.

## Development

```bash
make dev             # build and install for local development
make check           # full quality suite (Go + docs)
make docs-generate   # regenerate documentation from templates
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See [LICENSE](LICENSE).
