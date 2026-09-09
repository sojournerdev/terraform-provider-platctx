# ADR-005: Terraform Context is a structural contract

## Status

Accepted

## Date

2026-08-31

## Context

Terraform is structurally typed. Users can create or modify any Context-shaped object, so shape does not prove prior validation.

## Decision

Context is a validated structural contract, not a nominal Terraform type.

Every public Context consumer must validate its input independently:

- structure and field types;
- required, optional, null, and unknown states;
- semantic string rules;
- closed domains and Context invariants.

No consumer may assume that input came from `canonicalize()`.

## Alternatives

- **Treat `canonicalize()` as a constructor:** rejected because Terraform does not preserve nominal provenance.
- **Add a marker or hash:** rejected because structural values can be copied or forged.
- **Trust module output:** rejected because modules can return manually built objects.

## Consequences

Validation is reusable internally but repeated at every public boundary. Canonical output can be deterministic without becoming a nominal type.

## References

- [Terraform types and values](https://developer.hashicorp.com/terraform/language/expressions/types)
- [Terraform type constraints](https://developer.hashicorp.com/terraform/language/expressions/type-constraints)
