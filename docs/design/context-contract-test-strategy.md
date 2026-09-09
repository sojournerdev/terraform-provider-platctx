# Context Contract Test Strategy

## Purpose

Verify behavior at the layer where each contract is defined. Tests MUST check
observable behavior, boundaries, invariants, and transitions rather than source
structure or incidental implementation details.

## Test layers

### Domain tests

Test known values without Terraform dependencies:

- required fields and optional absence;
- exact identity comparison;
- opaque Unicode and punctuation preservation;
- invalid empty and surrounding-whitespace values;
- Criticality vocabulary;
- owner and operator distinction;
- governance normalization;
- identity mutation behavior;
- determinism; and
- idempotence.

### Adapter tests

Use real Terraform Framework values and diagnostics to test:

- dynamic input with no type coercion;
- fixed typed result construction;
- omitted and explicit-null values;
- whole-object unknown values;
- partially unknown nested values;
- known invalid values beside unknown values;
- malformed nested objects;
- unsupported attributes; and
- typed-null canonical output.

Do not replace unknown values with test zero values. The test input must carry
the same state a Terraform configuration can produce.

### Provider tests

Exercise the provider through its Framework function surface. Verify:

- provider metadata and protocol compatibility;
- the candidate function argument and result types;
- diagnostics crossing the public boundary;
- no provider configuration is required; and
- repeated calls return equivalent results.

### CLI smoke tests

Run a minimal Terraform configuration that invokes the candidate function with:

1. a valid minimal Context;
2. a fully populated Context;
3. an omitted optional value;
4. an invalid identifier; and
5. an unknown value produced by a plan-time expression.

The smoke test verifies wiring, not every domain rule.

## Test oracle

Use [Context Contract Conformance](context-contract-conformance.md) as the
single expected-result matrix. Detailed input facts remain in the normal and
adversarial fixture documents.

Tests MUST distinguish:

- a valid input that is accepted;
- an invalid input that is correctly rejected;
- unknown information that is correctly preserved; and
- a pressure case that requires design review.

## Properties

The implementation MUST satisfy these properties for accepted values:

```text
C(C(x)) = C(x)
```

Canonicalization MUST not depend on time, host, process environment, provider
configuration, credentials, filesystem state, network access, randomness, or
prior calls.

Known strings MUST be byte-for-byte unchanged after successful canonicalization.

## Diagnostics assertions

Tests MUST assert error severity, rule category, and attribute path. Tests SHOULD
avoid exact matching of prose unless the wording is explicitly part of the
public contract.

A rejected input MUST NOT produce a successful canonical result. An unknown
input MUST NOT produce a validation error solely because it is unknown.

## Test naming and fixtures

Test names SHOULD include the stable fixture ID, for example:

```text
TestCanonicalize_E8_PartiallyUnknownGovernance
```

Each fixture test SHOULD contain only the input, expected outcome, and relevant
invariant. Shared helpers MAY construct values, but helpers MUST NOT hide the
state transitions under test.

## Failure handling

When a test exposes a contract mismatch:

1. preserve the failing input as a fixture;
2. classify it as implementation defect, documentation defect, or contract
   falsifier;
3. update the relevant ADR and conformance row before changing behavior; and
4. rerun the affected layer and the full documentation/test checks.

Do not weaken a test to match an implementation detail.

## Coverage gate

Before release, the suite MUST cover every conformance row, every required
Framework case, every falsifier experiment, and every public diagnostic path.
A coverage percentage alone is not evidence of contract coverage.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Conformance matrix](context-contract-conformance.md)
- [Conformance fixture format](conformance-fixture-format.md)
- [Diagnostic catalog](diagnostic-catalog.md)
- [Provider engineering standards](provider-engineering-standards.md)
- [Falsification plan](context-contract-falsification.md)
