# Planx Common

## Collaboration and references

Use [workspace guidance](../AGENTS.md) for task completion, authorization and
proportionate checks. Read relevant clauses of the [canonical contract](../planx-spec/AI_CONTRACT.md),
[architecture](../planx-spec/planx-architecture.md) and accepted
[ADR-017](../planx-spec/adr/017-builtin-typed-data-integration.md) when their
subject changes. Use [repo.lock](repo.lock) for source ownership. Do not reread
the entire specification for an unrelated edit.

## Ownership

Provide Engine-side infrastructure utilities only. This is not a foundation
library for all modules. No runtime, SPI, protocol or domain logic.
SDK and plugins must not import Common. Common must not depend on SDK, plugins
or Proto. Metrics/telemetry ownership remains in `repo.lock`; do not move
pipeline/session semantics into infrastructure initialization.
Use affected package checks, then `go test ./...` for module changes.
