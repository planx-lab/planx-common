# AI RULES — PLANX COMMON (MANDATORY)

## Authority Documents

Before working here, read:
1. [planx-architecture.md](../planx-architecture.md)
2. [planx-ai-guardrails.md](../planx-ai-guardrails.md)
3. [AI_CONTRACT.md](../AI_CONTRACT.md)

---

## SCOPE

This repository provides Engine-side infrastructure utilities ONLY.
planx-common is NOT a shared foundation library.
It is an engine-side infrastructure utility set.

---

## IMPORT RESTRICTIONS
- planx-sdk-*
- planx-plugin-*
- planx-proto

It contains no runtime logic, no SPI, and no protocol definitions.
