# Conformance Fixture Format

## Purpose

Keep fixture intent separate from test implementation while preserving Terraform
known, null, and unknown states.

## Fixture record

Each fixture has these fields:

```text
id
source
claim
input
expected_outcome
expected_value
expected_diagnostics
invariants
```

- `id` is stable and unique, such as `E8-partial-unknown-governance`.
- `source` identifies the design fixture or real requirement.
- `claim` identifies the falsification claim.
- `input` preserves Terraform value state.
- `expected_outcome` is `accept`, `reject`, `preserve_unknown`, or `pressure`.
- `expected_value` is required for successful canonical output.
- `expected_diagnostics` contains category and path, not incidental prose.
- `invariants` lists properties such as exact preservation or identity stability.

## State notation

Fixtures MUST represent each value as one of:

```text
known(value)
null(type)
unknown(type)
```

An object is known when its container is known even if one or more members are
unknown. A whole-object unknown retains its declared result type.

The fixture format MUST NOT encode unknown as a sentinel string such as
`"<unknown>"`, an empty value, or null.

## Example

```text
id: E8-partial-unknown-governance
source: context-contract-adversarial-fixtures.md#E8
claim: null-and-unknown-distinction
input:
  identity:
    name: known("api")
  ownership:
    owned_by: known("team:payments")
  governance:
    criticality: unknown(string)
    cost_center: null(string)
expected_outcome: preserve_unknown
expected_value:
  governance:
    criticality: unknown(string)
    cost_center: null(string)
expected_diagnostics: []
invariants:
  - known_siblings_preserved
  - unknown_not_null
```

The exact serialization format may be JSON, YAML, or Go test data. The chosen
format MUST preserve type information, unknown state, and diagnostic paths
without lossy conversion.

## Fixture ownership

Every fixture has one owner document. Tests reference the fixture ID instead of
copying prose descriptions. A behavior change requires updating the fixture,
its source rationale, and the relevant ADR before changing the test expectation.

## Related documents

- [Conformance matrix](context-contract-conformance.md)
- [Falsification plan](context-contract-falsification.md)
- [Test strategy](context-contract-test-strategy.md)
