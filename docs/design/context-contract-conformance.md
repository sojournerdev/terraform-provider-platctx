# Context Contract Conformance

## Purpose

Turn the design fixtures into an executable oracle. A conformance test passes
when the observed result matches the outcome and invariant in this document.

## Outcomes

- **Accept:** return the canonical Context and no error diagnostics.
- **Reject:** return error diagnostics; no canonical result is published.
- **Preserve unknown:** return no validation error for unknown information and
  preserve its location and state.
- **Pressure:** record the result and design impact; do not silently broaden the
  v0.1 contract.

## Expected result matrix

| IDs | Input condition | Expected outcome |
| --- | --- | --- |
| H1-H5 | Valid minimal, full, shared, Unicode, or equal-owner Context | Accept; preserve known values exactly |
| B1 | Empty or whitespace-only known string | Reject; path identifies the string |
| B2 | Leading or trailing Unicode whitespace | Reject; input is not trimmed |
| B3 | Null or omitted required Context, object, or field | Reject |
| B4 | Wrong type, malformed object, or unsupported attribute | Reject; retain full path when available |
| B5 | `criticality` in the four-value domain | Accept |
| B5 | Empty, case-variant, or unrecognized `criticality` | Reject at `governance.criticality` |
| B6 | Missing optional fact | Accept; canonical typed null |
| B7 | One Context proposed for multiple logical subjects | Pressure/reject the model; do not invent an aggregate fact |
| B8 | Provider, deployment, or representation field | Reject as unsupported structure |
| E1 | Exact and near-equal identity tuples | Compare exact values; only exact tuples match |
| E2 | Same identity with different environments | Accept; environment is not identity |
| E3-E4 | Owner or operator change/removal | Accept when values remain valid; identity is unchanged |
| E5 | Contributors versus genuine joint accountability | Accept contributors outside the contract; record joint accountability as pressure |
| E6 | Omitted, null, empty, or valid optional string | Null/null/reject/preserve respectively |
| E7 | Empty, null, partially known, or unknown governance | Null/null/preserve object/preserve unknown object |
| E8 | Whole or partially unknown input | Preserve unknown; validate known siblings immediately |
| E9 | Each field mutation | Only name or namespace changes identity |
| E10 | Azure topology changes | Context identity and facts remain independent of resources |
| E11 | Repeated or changed execution conditions | Same input produces same result; `C(C(x)) = C(x)` |
| E12 | Schema changes | Apply compatibility rules; provider-only fields remain rejected |

## Canonical output rules

For accepted known input:

- all required fields are known and present;
- known optional values are unchanged;
- absent optional values are typed nulls;
- empty governance is typed null;
- equal owner and operator values remain distinct fields; and
- no additional field or generated value appears.

For unknown input:

- unknown remains unknown;
- unknown does not become null;
- known valid siblings remain known; and
- a known invalid sibling still produces an error.

## Diagnostic oracle

Rejecting tests MUST verify:

- error severity;
- rule category;
- affected attribute path; and
- absence of a successful canonical result.

Tests SHOULD avoid matching incidental diagnostic prose. If the Framework does
not expose stable categories, define stable summaries before implementing the
public operation.

## Required Framework cases

The adapter test suite MUST include at least:

1. unknown top-level Context;
2. unknown `identity.name`;
3. unknown optional `environment`;
4. one unknown governance member with known siblings;
5. one known invalid member beside one unknown member;
6. omitted and explicit-null optional attributes;
7. empty governance;
8. malformed nested objects;
9. unsupported attributes; and
10. inspection of every typed-null result.

## Pressure handling

Pressure cases are evidence about model limits. They MUST produce a written
result in the falsification report. They MUST NOT be resolved by adding a
field, default, provider-specific value, or silent lossy conversion.

## Source fixtures

The detailed examples and facts remain in:

- [Normal fixtures](context-contract-fixtures.md)
- [Adversarial fixtures](context-contract-adversarial-fixtures.md)

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Test strategy](context-contract-test-strategy.md)
- [Conformance fixture format](conformance-fixture-format.md)
- [Diagnostic catalog](diagnostic-catalog.md)
- [Falsification plan](context-contract-falsification.md)
