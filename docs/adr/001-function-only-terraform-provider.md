# ADR-001: Function-only Terraform provider

## Status

Accepted

## Date

2026-08-31

## Context

Context Core needs local, deterministic operations. It does not manage or query remote systems.

## Decision

If Context Core exposes a Terraform API, it will use a function-only provider:

- no provider configuration;
- no resources or data sources;
- no network, filesystem, process environment, time, or random inputs;
- Terraform 1.8 or later.

This decision does not justify `canonicalize()` by itself.

## Alternatives

- **Resource:** rejected because Context has no remote lifecycle.
- **Data source:** rejected because Context needs no remote lookup.
- **Configured provider:** rejected because hidden inputs would make results unclear.
- **Terraform module:** remains valid if no provider function proves useful.

## Consequences

Operations remain pure and explicit. If no function earns a public API, we should not ship a provider merely to preserve this design.

## References

- [Provider-defined function concepts](https://developer.hashicorp.com/terraform/plugin/framework/functions/concepts)
- [Function-only provider tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-function-only)
