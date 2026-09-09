# ADR-006: Missing facts remain missing

## Status

Accepted

## Date

2026-08-31

## Context

Context Core cannot know why an optional fact was omitted and must not invent one. Terraform distinguishes known, null, and unknown values, but a fixed object represents omitted optional fields as typed nulls.

## Decision

- Omitted input and explicit null both mean semantic absence and canonicalize to typed null.
- Empty or surrounding-whitespace strings are invalid.
- Unknown values remain unknown.
- Known values are preserved exactly after validation.
- An optional object with no facts canonicalizes to null.
- An optional object containing an unknown member does not canonicalize to null.
- No semantic defaults are allowed.

## Alternatives

- **Keep omission distinct from null:** rejected because a fixed Terraform object cannot preserve that distinction without a tagged representation.
- **Apply defaults:** rejected because defaults would manufacture facts.
- **Treat empty strings as absent:** rejected because that would hide invalid input.

## Consequences

Canonical absence has one representation: typed null. Terraform null and unknown handling stays at the adapter boundary; the pure domain contains only validated known facts.

## References

- [Optional object attributes](https://developer.hashicorp.com/terraform/language/expressions/type-constraints)
- [Function argument handling](https://developer.hashicorp.com/terraform/plugin/framework/functions)
