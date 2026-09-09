# Context API v0.1

## Purpose

Define the smallest useful API for Context Core. This document is normative for
v0.1 scope; implementation details remain in the domain and Terraform boundary
documents.

## Candidate operation

Context Core will implement one candidate operation:

```text
canonicalize(context) -> canonical context
```

The operation accepts one structural Context value, validates its known facts,
and returns the same facts in canonical form. It does not look up, create, or
manage infrastructure.

## Caller-local composition

The intended caller keeps facts in Terraform locals and calls the operation
before passing the canonical value to one or more modules. The caller owns the
facts; the function only validates and canonicalizes them.

This provides deterministic behavior at a reusable module boundary:

```text
local facts → canonicalize(facts) → module inputs
```

The function MUST NOT discover, infer, default, or manage anything. The
provider-qualified Terraform name remains illustrative until the provider
source address is fixed.

The operation is a candidate public API. It is not approved for release until
the falsification plan shows reusable value beyond ordinary Terraform object
constraints and validation.

## Context shape

```hcl
{
  identity = {
    name      = string
    namespace = optional(string)
  }
  ownership = {
    owned_by    = string
    operated_by = optional(string)
  }
  environment = optional(string)
  governance = optional({
    criticality         = optional(string)
    cost_center         = optional(string)
    data_classification = optional(string)
  })
}
```

The function receives a dynamic Terraform input so it can validate the original
structure without coercion or loss. Its successful result is the fixed Context
object shown above.

Logical identity is the exact tuple `(namespace?, name)`. All other fields are
descriptive facts and do not affect identity.

## Semantics

- Known strings MUST be valid UTF-8, non-empty, and free of leading or trailing
  Unicode whitespace.
- Known values MUST be preserved exactly. No case folding, Unicode
  normalization, trimming, parsing, or sanitizing is allowed.
- `criticality` MUST be one of `low`, `medium`, `high`, or `critical`.
- Required known values MUST be present and valid.
- Optional omission and explicit null mean absence and canonicalize to typed
  null.
- Unknown values remain unknown. Unknown is never treated as null or a Go zero
  value.
- An empty governance object canonicalizes to typed null. Unknown governance
  information is preserved.
- No default creates a missing fact.
- Canonicalization MUST be deterministic and idempotent.

## Non-goals

v0.1 does not include:

- resources or data sources;
- provider configuration or network access;
- policy or organization-specific completeness checks;
- Azure, AWS, GCP, Kubernetes, or OpenTelemetry projections;
- relationships between subjects;
- deployment identity or generated IDs;
- provider-specific fields or extension maps;
- plural ownership.

## Release gate

The provider MUST NOT expose the candidate operation as a supported public API
unless all of these are true:

1. every conformance case has an executable expected result;
2. unknown and null values survive the Terraform Framework boundary;
3. diagnostics identify invalid paths without hiding other known invalid values;
4. determinism and idempotence are proven;
5. at least two concrete caller workflows, or one repeated cross-module
   boundary, reuse the behavior; and
6. comparison with ordinary Terraform constraints shows a material benefit.

If the gate fails, keep the domain validation useful internally and do not ship
a provider only to preserve the design.

## Related documents

- [Domain model and validation](domain-model-and-validation.md)
- [Consumer use cases](consumer-use-cases.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Conformance matrix](context-contract-conformance.md)
- [Falsification plan](context-contract-falsification.md)
- [ADR-008: Candidate canonicalize API](../adr/008-canonicalize-public-api.md)
- [ADR-009: Terraform function parameter strategy](../adr/009-terraform-function-parameter-strategy.md)
