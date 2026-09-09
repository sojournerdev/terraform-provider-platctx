# Provider Engineering Standards

## Purpose

Set the quality bar for a strong Terraform provider without expanding Context
Core beyond its purpose.

## Provider shape

The provider MUST focus on the Context problem domain. It MUST expose only
operations that are pure, offline, and useful to Terraform callers.

The provider MUST have:

- no required configuration;
- no resources or data sources in v0.1;
- stable provider metadata;
- documented function definitions and parameters; and
- a fixed, typed result for successful Context operations.

Each function MUST perform one computational operation. Options MUST NOT turn
one function into several unrelated behaviors.

## Implementation boundaries

Keep these layers separate:

1. **Terraform adapter:** value states, structural decoding, diagnostics, and
   result construction.
2. **Domain:** known validated facts, identity, and pure invariants.
3. **Provider wiring:** metadata, function registration, and protocol handling.

The domain MUST not import Terraform value types. The adapter MUST not invent
facts or bypass domain validation.

## API discipline

Public behavior MUST be specified before implementation. A new field or
operation requires:

- a documented use case;
- an ADR or an update to an accepted decision;
- conformance fixtures;
- unknown, null, and invalid-state tests; and
- compatibility analysis.

Do not add fields for anticipated projections, provider metadata, or speculative
future consumers.

## Terraform behavior

The provider MUST follow Terraform's known, null, and unknown model. It MUST
preserve unknown values and MUST never use Go zero values as semantic absence.

Every public consumer MUST validate structural and semantic input independently.
A canonicalization function is not a nominal type constructor.

## Code quality

The implementation SHOULD prefer small, explicit functions and standard library
behavior. It MUST avoid hidden global state, reflection where a typed operation
is sufficient, and unnecessary copies of large Terraform values.

Errors MUST be actionable and deterministic. Map iteration MUST NOT determine
observable diagnostic order.

## Testing and verification

Every public behavior requires:

- unit tests for domain invariants;
- Framework tests for value states and diagnostics;
- provider registration tests;
- Terraform CLI smoke coverage; and
- regression coverage for every fixed defect.

Tests MUST assert observable contracts, not private implementation structure.

## Documentation

A released function MUST have:

- a concise purpose statement;
- parameter and result documentation;
- valid and invalid examples;
- null and unknown behavior;
- diagnostic behavior;
- version requirements; and
- explicit non-goals.

## Compatibility

Follow Semantic Versioning for the provider. Treat changes to field meaning,
identity, result type, diagnostics relied upon by automation, or null/unknown
behavior as compatibility-sensitive changes.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Test strategy](context-contract-test-strategy.md)
- [Security and resource limits](security-and-resource-limits.md)
- [Support matrix](support-matrix.md)
