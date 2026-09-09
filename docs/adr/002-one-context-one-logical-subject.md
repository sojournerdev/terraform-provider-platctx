# ADR-002: One Context describes one logical subject

## Status

Accepted

## Date

2026-08-31

## Context

A Terraform module may manage many related subjects. Combining them in one Context would make identity and ownership unclear.

## Decision

One Context describes exactly one logical subject.

The subject may be an application or infrastructure. It is not defined by a Terraform module, resource, deployment, or provider object. Related subjects use separate Context values.

## Alternatives

- **One Context per module:** rejected because a module may contain zero, one, or many subjects.
- **One Context per resource:** rejected because a subject may exist before infrastructure and have several representations.
- **One Context for a subject graph:** deferred until relationships require it.

## Consequences

Identity and ownership always refer to one subject. Multiple deployments may share that logical identity. Relationships remain a separate problem and do not justify `kind` in v0.1.
