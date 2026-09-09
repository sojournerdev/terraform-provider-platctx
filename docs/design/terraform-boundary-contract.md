# Terraform Boundary Contract

## Purpose

Define how Context Core validates and canonicalizes Terraform values. Terraform
is structurally typed, so every public consumer validates its own input. No
consumer may trust that a value came from `canonicalize()`.

## Public value type

The Context value is a fixed object with these nested object attributes:

```hcl
object({
  identity = object({
    name      = string
    namespace = optional(string)
  })
  ownership = object({
    owned_by    = string
    operated_by = optional(string)
  })
  environment = optional(string)
  governance = optional(object({
    criticality         = optional(string)
    cost_center         = optional(string)
    data_classification = optional(string)
  }))
})
```

The provider MUST use one declared result type for every successful result.
Absent optional attributes in that result are typed nulls.

## State handling

The adapter MUST distinguish three states:

| State | Required value | Optional value | Action |
| --- | --- | --- | --- |
| Known valid | Preserve | Preserve | Continue |
| Known invalid | Diagnostic | Diagnostic | Reject |
| Null | Diagnostic | Typed null | Continue |
| Unknown | Preserve unknown | Preserve unknown | Defer semantic validation |

A known invalid value MUST be rejected even when another value is unknown. An
unknown value MUST NOT be replaced by null, an empty string, or a zero value.

A null top-level Context, identity object, ownership object, required name, or
required owner is invalid.

## Structural validation

Each public consumer MUST validate:

- top-level object shape;
- required nested objects;
- required and optional attribute types;
- null and unknown states;
- supported attributes;
- semantic string rules;
- the Criticality vocabulary; and
- Context invariants.

The implementation MUST verify how Terraform converts extra object attributes.
If conversion discards them before validation, the public input type or adapter
MUST change so that unsupported attributes are rejected as required by the
contract. This behavior requires a Framework test; it MUST NOT be inferred from
HCL type syntax alone.

## Canonicalization

For a known valid input:

- known values are copied without transformation;
- omitted and explicit-null optional values become typed nulls;
- an empty governance object becomes typed null;
- governance containing an unknown member remains an object;
- equal owner and operator values remain two attributes;
- no default or generated value is added.

Canonicalization MUST have no network, filesystem, process, time, credential,
random, cache, or provider-configuration dependency.

## Diagnostics

Diagnostics MUST include:

- error severity;
- a stable summary identifying the violated rule;
- a concise detail naming the expected condition; and
- the complete Terraform attribute path when a path exists.

Tests SHOULD assert severity, rule category, and path. Tests SHOULD NOT depend on
incidental prose when a stable diagnostic code or category can be used.

The adapter MUST return all independently discoverable known errors in one
operation where the Framework permits it. Unknown values do not produce an
error merely because they are unknown.

## Provider operation

The candidate operation has one dynamic argument interpreted as a Context and
one fixed typed result:

```text
canonicalize(context: dynamic Context input) -> Context
```

The dynamic input is intentional. It prevents Terraform from discarding
unsupported attributes or coercing invalid primitive values before provider
validation. The provider MUST validate the original value and MUST return the
fixed Context result type on success.

The function definition MUST opt into null and unknown argument handling. A
top-level null produces a provider diagnostic. A top-level unknown produces an
unknown result of the fixed Context result type. A known object containing
unknown members reaches provider logic and preserves those members.

The provider MUST have no configuration, resources, data sources, or remote
side effects. The provider-facing qualified function name follows the final
provider source address and is recorded in the support matrix before release.

## Boundary invariants

For every accepted value `x`:

```text
C(C(x)) = C(x)
```

For every known value, canonicalization preserves exact string content. For
every unknown value, canonicalization preserves its unknown state and does not
claim knowledge that Terraform has not provided.

## Required spike

Before freezing the adapter, implement a small Framework experiment covering:

1. an unknown top-level object;
2. an object with one unknown nested member;
3. known invalid input beside an unknown member;
4. omitted and explicit-null optional members;
5. empty governance;
6. an unsupported nested attribute; and
7. typed-null result inspection.

Record observed Framework behavior in the conformance tests and update this
document if the proposed structural boundary cannot express the contract.

## Related documents

- [Diagnostic catalog](diagnostic-catalog.md)
- [ADR-009: Terraform function parameter strategy](../adr/009-terraform-function-parameter-strategy.md)
- [Context API v0.1](context-api-v0.1.md)
- [Domain model and validation](domain-model-and-validation.md)
- [Conformance matrix](context-contract-conformance.md)
- [ADR-005: Structural Terraform contract](../adr/005-structural-terraform-contract.md)
- [ADR-006: Missing facts remain missing](../adr/006-missing-facts-remain-missing.md)
