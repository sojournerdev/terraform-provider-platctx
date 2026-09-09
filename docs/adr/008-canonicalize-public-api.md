# ADR-008: Treat `canonicalize` as a release-gated candidate API

## Status

Accepted

## Date

2026-08-31

## Context

The Context contract has reusable validation and normalization rules, but the
fixtures do not prove that a public Terraform function is better than ordinary
Terraform object constraints and consumer validation. Shipping an API before
that comparison would add surface area without evidence.

## Decision

Implement `canonicalize(context)` as a candidate operation so its behavior can
be tested at the domain, Terraform Framework, provider, and CLI boundaries.

Do not treat the operation as a supported public API until the falsification
plan proves all of the following:

- canonical output is deterministic and idempotent;
- known facts are preserved exactly;
- null and unknown states survive the Framework boundary;
- invalid structural and semantic values produce useful diagnostics;
- at least two independent consumers can reuse the behavior; and
- the operation provides material value beyond ordinary Terraform constraints
  and validation.

If the evidence fails, retain useful internal validation and do not ship a
provider function solely to preserve this design.

## Alternatives

- **Ship `canonicalize` immediately:** rejected because usefulness is unproven.
- **Build no operation:** deferred until the boundary and reuse experiments are
  complete.
- **Use a Terraform module:** remains valid when a provider function does not
  earn a public API, but does not provide shared provider-side validation.
- **Add resources or data sources:** rejected because Context has no remote
  lifecycle or lookup.

## Consequences

Implementation can proceed without prematurely promising a public API. Tests
must compare the candidate with a baseline of ordinary Terraform constraints
and validation. Release documentation must record the evidence and final
status.

## References

- [Context API v0.1](../design/context-api-v0.1.md)
- [Falsification plan](../design/context-contract-falsification.md)
- [ADR-001: Function-only Terraform provider](001-function-only-terraform-provider.md)
- [ADR-005: Structural Terraform contract](005-structural-terraform-contract.md)
