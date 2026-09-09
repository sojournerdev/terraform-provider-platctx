# Domain Model and Validation

## Purpose

Separate provider-neutral Context facts from Terraform state handling. The
pure domain contains only known, validated facts. The Terraform adapter owns
null, unknown, structural conversion, and diagnostics.

## Domain model

The domain model represents one logical subject:

```text
Context
├── Identity
│   ├── Name       (required)
│   └── Namespace  (optional)
├── Ownership
│   ├── OwnedBy    (required)
│   └── OperatedBy (optional)
├── Environment   (optional)
└── Governance    (optional)
    ├── Criticality
    ├── CostCenter
    └── DataClassification
```

The implementation SHOULD use distinct defined string types for semantically
different identifiers. Raw strings MUST NOT be interchangeable at domain API
boundaries.

The domain model MUST NOT contain Terraform unknown values. Optional domain
facts use absence; the adapter preserves unknown values before constructing a
known domain value.

## Identity

Logical identity is `(namespace?, name)`.

- `name` is required.
- `namespace` is an optional collision domain.
- Comparison is exact.
- Case, Unicode normalization, and byte representation are not changed.
- Ownership, operation, environment, and governance do not affect identity.
- No UUID, hash, canonical ID, or deployment ID is generated.

Changing `name` or `namespace` changes identity. Any other valid mutation keeps
the same identity.

## Validation rules

Every known string MUST:

1. be valid UTF-8;
2. contain at least one non-whitespace rune; and
3. have no leading or trailing Unicode whitespace.

The validator MUST preserve internal whitespace, punctuation, case, and Unicode.
It MUST NOT trim, normalize, parse, sanitize, truncate, or infer values.

`criticality` is the only closed domain:

```text
low | medium | high | critical
```

All other string fields are opaque organization-defined identifiers.

## Optional values

The domain distinguishes only present valid facts and absent facts. The adapter
maps both omitted and explicit-null optional Terraform values to absence.

An empty string is invalid. It is never an alias for absence.

Governance is absent when:

- it is omitted;
- it is null;
- it is an empty object; or
- all of its members are null.

Governance is present when at least one member is a known or unknown Terraform
value. Once in the domain, all members are known.

## Ownership

`owned_by` means ultimate accountability. `operated_by` means routine
operation. Equal values are valid and retain both meanings. A missing operator
does not imply that the owner operates the subject.

Ownership is singular in v0.1. A concrete jointly accountable subject that
cannot name one durable accountability principal is a documented falsification
case, not a reason to silently discard an owner.

## Domain operations

The implementation MUST provide behavior equivalent to these operations:

```text
Validate(Context) -> []ValidationIssue
Canonicalize(Context) -> Context
SameIdentity(Context, Context) -> bool
```

`Validate` reports every independently invalid known fact. Each issue contains
a stable attribute path and a stable rule identifier. The Terraform adapter
converts issues into user-facing diagnostics; the domain produces no error
strings or log output.

`Canonicalize` MUST be deterministic and idempotent. It MUST not perform I/O or
read process, provider, or environment state. It MUST not manufacture facts.

The Terraform adapter MAY combine decoding, validation, and canonical output,
but it MUST preserve the same domain rules and state semantics.

## Validation issues

Validation issues SHOULD have stable rules for:

- missing required value;
- empty or whitespace-padded string;
- invalid Criticality; and
- invalid UTF-8.

An issue SHOULD identify the field path and the violated rule. Path and rule
are the stable contract. Wording belongs in the Terraform adapter, not the
domain.

## Non-responsibilities

The domain does not:

- enforce organization-specific completeness;
- infer facts from other facts;
- select infrastructure;
- validate Azure, AWS, GCP, Kubernetes, or OpenTelemetry limits;
- manage relationships; or
- decide deployment identity.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Conformance matrix](context-contract-conformance.md)
- [ADR-003: Context contains facts](../adr/003-context-contains-facts.md)
- [ADR-004: Opaque semantic identifiers](../adr/004-opaque-semantic-identifiers.md)
