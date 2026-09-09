# Diagnostic Catalog

## Purpose

Define stable, useful diagnostics for invalid Context values. Diagnostics explain
what failed, where it failed, and how to correct it without exposing raw input
values unnecessarily.

## Categories

| Category | Summary | Typical path |
| --- | --- | --- |
| `context.invalid_type` | Invalid Context type | `context` |
| `context.required` | Required value is missing | Required field |
| `context.unsupported_attribute` | Unsupported Context attribute | Unknown field |
| `context.invalid_string` | Invalid identifier value | String field |
| `context.invalid_criticality` | Unsupported criticality | `governance.criticality` |
| `context.invalid_structure` | Invalid nested Context structure | Nested object |
| `context.invalid_invariant` | Context invariant is violated | Context or affected field |

The category names are implementation-facing stable identifiers. The provider
MUST keep them stable within a major release.

## Rules

- Errors use error severity.
- A diagnostic identifies the complete Terraform attribute path when available.
- A diagnostic names the expected rule in concise language.
- Diagnostics MUST NOT depend on network, provider configuration, or external
  state.
- Diagnostics SHOULD avoid echoing raw identifiers, which may contain sensitive
  organizational information.
- Unknown values do not produce an error merely because they are unknown.
- Known invalid values produce diagnostics immediately.
- Independently discoverable errors SHOULD be returned together in deterministic
  path order.
- A failed operation MUST NOT publish a successful canonical result.

## Required cases

| Input | Category | Path |
| --- | --- | --- |
| Null or non-object Context | `context.invalid_type` | `context` |
| Missing or null `identity` | `context.required` | `identity` |
| Missing or null `ownership` | `context.required` | `ownership` |
| Missing or null `identity.name` | `context.required` | `identity.name` |
| Missing or null `ownership.owned_by` | `context.required` | `ownership.owned_by` |
| Non-string known field | `context.invalid_type` | Affected field |
| Empty or padded known string | `context.invalid_string` | Affected field |
| Unsupported attribute | `context.unsupported_attribute` | Full unknown path |
| Invalid governance shape | `context.invalid_structure` | `governance` |
| Unsupported Criticality | `context.invalid_criticality` | `governance.criticality` |
| Violated cross-field rule | `context.invalid_invariant` | Context or affected field |

## Stable ordering

When multiple diagnostics are returned, sort by attribute path using a documented
lexicographic order. If two diagnostics share a path, sort by category. The
order MUST NOT depend on map iteration order.

## Testing

Tests MUST assert severity, category, and path. Tests MAY assert summary text
when it is part of the published user contract. Tests SHOULD avoid matching
unstructured detail text.

The catalog does not define Terraform's own static type errors. It defines
provider diagnostics after the dynamic parameter has delivered the value.

## Related documents

- [Terraform boundary contract](terraform-boundary-contract.md)
- [Conformance matrix](context-contract-conformance.md)
- [Test strategy](context-contract-test-strategy.md)
