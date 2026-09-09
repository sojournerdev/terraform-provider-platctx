# Context Contract Falsification

## Purpose

Try to disprove the Context contract before releasing code. A passing fixture
shows compatibility with the proposal; it does not prove the proposal is
necessary or complete.

## Claims under test

1. `(namespace?, name)` is sufficient logical identity.
2. One Context can describe both application and infrastructure subjects.
3. Singular `owned_by` is sufficient for v0.1.
4. `environment` is coherent as deployment-context information outside identity.
5. Null and unknown values can survive the Terraform Framework boundary.
6. Validation and canonicalization are deterministic and idempotent.
7. Canonicalization provides reusable value beyond ordinary Terraform object
   constraints and validation.
8. Provider-neutral facts can remain unchanged when a downstream projection has
   stricter limits.

## Methods

Use four evidence sources:

- pure domain tests;
- Terraform Framework boundary tests;
- Terraform CLI smoke tests; and
- concrete subject and consumer examples.

Use the normal and adversarial fixtures as the initial corpus. Add a fixture
only when it represents a real requirement, failure, or boundary. Every added
fixture records its source, expected outcome, and the claim it tests.

## Falsifiers

The contract is falsified, or must be revised, when any of these occurs:

- a truthful subject cannot be represented without a fake, overloaded, or
  provider-specific field;
- two distinct subjects are forced to share identity;
- a concrete jointly accountable subject cannot name one durable accountability
  principal and no accepted organizational group represents it;
- a real operation requires `kind`, relationships, deployment identity, or a
  generated ID;
- a valid unknown or null state is changed, lost, or misclassified at the
  Terraform boundary;
- a known invalid value is accepted, or a known invalid sibling is hidden by an
  unknown value;
- canonicalization changes a known value, depends on external state, or is not
  idempotent;
- ordinary Terraform constraints and validation provide the same reusable
  behavior with less public API and less semantic risk; or
- a downstream projection requires changing the canonical fact rather than
  defining an explicit mapping and loss policy.

Rejecting deliberately invalid input is not a falsification. It is evidence of
correct validation.

## Public API comparison

Compare the candidate operation with a baseline consisting of:

1. a caller-owned Terraform local;
2. a fixed Terraform object type where practical;
3. ordinary attribute validation; and
4. consumer-local validation where needed.

The candidate earns a public API only if it provides a clear, reusable benefit
at a caller boundary, such as one canonical representation shared across
multiple module inputs or two concrete caller workflows. It must also preserve
all required structural and semantic errors. Code reduction alone is not
sufficient if it relies on nominal provenance that Terraform does not preserve.

## Required experiments

| Experiment | Evidence | Failure condition |
| --- | --- | --- |
| Run all conformance cases | Test report | Any unexpected outcome |
| Repeat accepted canonicalization | Equal serialized values | Result changes between runs |
| Canonicalize canonical output | Equal output | Idempotence fails |
| Vary host, time, credentials, and provider config | Equal output | External state changes result |
| Pass unknown top-level and nested values | State-aware result | Unknown becomes null, zero, or error |
| Mix known invalid and unknown values | Diagnostics | Invalid known value is deferred or hidden |
| Pass unsupported attributes | Structural result | Extra fields are silently accepted when rejection is required |
| Compare baseline and candidate | Design record | Candidate adds API without material benefit |
| Exercise caller workflows | Shared canonical output | No reusable caller boundary is demonstrated |

## Decision rule

Record one of these outcomes:

- **Retain:** evidence supports the contract and no falsifier occurred.
- **Revise:** a falsifier identifies a bounded change; update the ADR and
  fixtures before implementation continues.
- **Reject public API:** the domain remains useful internally, but the provider
  operation does not earn release.
- **Defer:** evidence is insufficient; identify the missing experiment and do
  not claim acceptance.

The report MUST link each conclusion to tests, fixtures, or a concrete example.
Assertions without evidence remain proposals.

## Current expected pressure

The known pressure points are:

- `environment` has deployment scope while identity has subject scope;
- singular ownership may not represent genuine joint accountability; and
- Terraform structural conversion may not expose unsupported attributes for
  rejection.

These are explicit experiments, not reasons to add speculative fields.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Conformance matrix](context-contract-conformance.md)
- [Test strategy](context-contract-test-strategy.md)
- [Consumer use cases](consumer-use-cases.md)
- [ADR-008: Candidate canonicalize API](../adr/008-canonicalize-public-api.md)
