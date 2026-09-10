---
page_title: "Provider: PlatCtx"
description: |-
  Provider-agnostic service metadata for platform engineering.
---

# PlatCtx Provider

PlatCtx gives you a single, typed structure for service metadata that works across cloud providers and modules. Instead of duplicating identity, ownership, environment, and governance data for every provider, you define it once with `provider::platctx::canonicalize` and pass the result downstream.

If you only use one provider and a few modules, you probably do not need this. PlatCtx is built for teams that need consistent metadata across many providers, environments, and service catalogs.

## Requirements

| Dependency | Version |
| --- | --- |
| [Terraform](https://developer.hashicorp.com/terraform/downloads) | >= 1.8 |
| [Go](https://golang.org/doc/install) (for development) | >= 1.27 |

## Example Usage

A `provider "platctx"` block is not required. Only the `required_providers` block is needed.

```terraform
terraform {
  required_version = ">= 1.8"

  required_providers {
    platctx = {
      source  = "sojournerdev/platctx"
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
      owned_by = "team:platform"
    }
  })
}
```

## Security

This provider makes no network requests, reads no environment variables, and collects no telemetry. All computation is local and deterministic.

## Updating Your Lockfile

`terraform init` updates the lockfile for your current platform only. For teams that work across different operating systems or architectures, run `terraform providers lock` to pre-populate checksums for every platform your team uses. This prevents each person's `terraform init` from modifying the lockfile with platform-specific entries.

This provider is built for the following platforms:

| OS | Architectures |
| --- | --- |
| macOS | arm64 |
| Linux | amd64, arm64 |
| Windows | amd64, arm64 |
| FreeBSD | amd64, arm64 |

Run this command from your root module directory:

```shell
terraform providers lock \
    -platform=darwin_arm64 \
    -platform=linux_amd64 \
    -platform=linux_arm64 \
    -platform=windows_amd64 \
    -platform=windows_arm64 \
    sojournerdev/platctx
```

Omit platforms your team does not use. For example, if your team only uses macOS and Linux, run:

```shell
terraform providers lock \
    -platform=darwin_arm64 \
    -platform=linux_amd64 \
    -platform=linux_arm64 \
    sojournerdev/platctx
```

For more details, see the Terraform documentation on [`terraform providers lock`](https://developer.hashicorp.com/terraform/cli/commands/providers/lock).

## Argument Reference

This provider has no configuration arguments.
