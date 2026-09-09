# ADR-004: Organization-defined identifiers remain opaque

## Status

Accepted

## Date

2026-08-31

## Context

Organizations define their own owner references, environments, cost centers, data classifications, and identity namespaces. Context Core does not own those taxonomies.

## Decision

The Go domain model uses distinct opaque string types for these fields.

Known values must:

- contain valid UTF-8;
- be non-empty;
- have no leading or trailing Unicode whitespace;
- preserve case, Unicode, and internal characters.

Context Core does not parse, normalize, sanitize, truncate, or infer meaning from them.

## Alternatives

- **Raw strings throughout Go:** rejected because unrelated fields could be mixed accidentally.
- **Provider-defined enums:** rejected because these domains belong to each organization.
- **Provider-safe syntax:** rejected because downstream restrictions are projection concerns.

## Consequences

The domain gains compile-time separation without imposing a taxonomy. Future projections may reject or encode a value, but cannot change the canonical fact.
