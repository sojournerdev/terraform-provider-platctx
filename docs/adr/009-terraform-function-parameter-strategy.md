# ADR-009: Use a dynamic input for strict Context validation

## Status

Accepted

## Date

2026-08-31

## Context

The Context contract requires all of the following:

- callers may omit optional attributes;
- unsupported attributes are rejected;
- numbers and booleans are not silently converted to strings;
- nested null and unknown states are preserved; and
- malformed structures produce Context-specific diagnostics.

A statically typed Terraform object parameter cannot guarantee these properties.
Terraform object parameters require declared attributes, and Terraform object
conversion can discard extra attributes and coerce compatible primitive values.

## Decision

The candidate function will use a Terraform Framework dynamic parameter with
both null and unknown handling enabled.

The implementation will:

1. receive the original Terraform value without type conversion;
2. reject null top-level input with a Context diagnostic;
3. preserve an unknown top-level input as an unknown result of the fixed Context
   result type;
4. inspect the underlying value type and require an object;
5. reject unsupported attributes;
6. validate nested types, nulls, unknowns, and semantic values; and
7. return one fixed typed Context object on success.

This choice preserves the Context contract at the cost of moving input shape
validation from Terraform's static schema into provider logic.

## Alternatives

- **Static object parameter:** rejected because optional omission, unsupported
  attribute rejection, and exact primitive-type rejection cannot all be
  preserved.
- **Static object parameter with explicit nulls:** rejected because it changes
  the caller contract and still permits lossy extra-attribute conversion.
- **Accept Terraform coercion:** rejected because it turns invalid facts into
  different valid-looking facts.
- **Accept any dynamic value without validation:** rejected because dynamic
  input is useful only with explicit structural validation.

## Consequences

The provider must maintain recursive structural validation and comprehensive
Framework tests. The result remains statically typed for downstream Terraform
expressions. The function definition must document that the input is dynamic
because strict preservation is intentional, not because the Context shape is
unspecified.

Top-level null and unknown behavior is provider-controlled. Nested null and
unknown values remain visible inside known object values and must be handled by
the adapter.

## References

- [Dynamic function parameters](https://developer.hashicorp.com/terraform/plugin/framework/functions/parameters/dynamic)
- [Object function parameters](https://developer.hashicorp.com/terraform/plugin/framework/functions/parameters/object)
- [Terraform type conversion](https://developer.hashicorp.com/terraform/language/expressions/type-constraints)
- [Terraform boundary contract](../design/terraform-boundary-contract.md)
