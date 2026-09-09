# Terraform Provider PlatCtx

[![Terraform Registry](https://img.shields.io/terraform/v/sojournerdev/platctx?label=Terraform%20Registry)](https://registry.terraform.io/providers/sojournerdev/platctx/latest)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)

Canonicalize service metadata for platform engineering.

## Why Use This

Service metadata like identity, ownership, environment, and compliance facts often lives in different places depending on the cloud provider. When you move workloads across providers or environments, that context breaks.

PlatCtx gives you provider-agnostic service metadata that travels with your resources. The `canonicalize` function validates and normalizes this metadata into a consistent typed structure so downstream tools and modules can rely on it.

This metadata powers service catalogs like Backstage, where it drives service discovery, ownership tracking, and compliance reporting. Inspired by [Cloud Posse's context pattern](https://docs.cloudposse.com/modules/library/provider/context/), PlatCtx extends the idea from resource tagging to platform-level context.

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
