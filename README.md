# Terraform Provider PlatCtx

## Overview

A Terraform provider that canonicalizes service metadata across cloud providers.

I was building a platform monorepo with shared modules and runtime infrastructure.
The pattern I was trying to implement was the concept of context inspired by [Cloud Posse's context pattern](https://docs.cloudposse.com/modules/library/provider/context/).
The context carried identity, ownership, and governance data, but different providers
expect different shapes. Azure wants tags, AWS wants different fields, and duplicating
the same data for each provider gets out of hand fast.

PlatCtx solves this by letting you define context once and normalizing it into a
provider-agnostic typed structure with `provider::platctx::canonicalize`. Downstream
modules get consistent, validated input without knowing which provider consumed it.

If you only use one provider and a few modules, you probably don't need this. PlatCtx
is built for teams managing metadata across multiple providers, environments, and
service catalogs.

See the [Terraform Registry documentation](https://registry.terraform.io/providers/sojournerdev/platctx/latest/docs) for function reference, examples, and configuration details.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.27 (for development)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See [LICENSE](LICENSE).
