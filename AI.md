# planx-common

This repository provides Engine-side infrastructure utilities ONLY.
planx-common is NOT a shared foundation library.
It is an engine-side infrastructure utility set.

It MUST NOT be imported by:
- planx-sdk-*
- planx-plugin-*
- planx-proto

It contains no runtime logic, no SPI, and no protocol definitions.
