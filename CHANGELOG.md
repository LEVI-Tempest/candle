# Changelog

All notable changes to this project are documented in this file.

## 2026-03-09

### Added
- Added repository-wide `AGENTS.md` and strengthened engineering rules:
  - architecture-first boundaries
  - prefer mature libraries over reimplementation
  - explicit delivery/validation expectations
- Added volume-price evidence pipeline in `pkg/identify`:
  - `PatternEvidence` model and factorized evidence output
  - indicator computation with `go-talib` (`OBV`, `MFI`, `AD`) plus `CMF`/`VPT`
  - context/volume scoring and contradiction handling
- Added `cmd/signal` JSON report CLI for offline/online analysis.
- Added signal output schema: `docs/signal.schema.json`.
- Added contract/behavior tests for signal and evidence paths.
- Added new reusable business layer `pkg/signal`:
  - config loading/merging (`--config`)
  - MA-based trend filter (`up/down/sideway/unknown`)
  - 100-point decision score (`strong/medium/weak`)
  - pattern-level reasons and volume states
  - forward return fields (`3/5/10`) and CSV logging support

### Changed
- Refactored `cmd/signal` to a thin entrypoint (I/O + flags), moving core logic into `pkg/signal`.
- Integrated structured evidences into `pkg/charting` detection flow.
- Updated `README.md` with signal CLI usage, config usage, and schema reference.
- Updated `docs/个人项目_功能指引_轻量版.md` with current implementation status.

### Fixed
- Fixed short-series panic risk when calling `go-talib` MFI by applying safe effective period fallback.

### Validation
- Passed package-level tests:
  - `go test ./cmd/signal ./pkg/signal ./pkg/identify ./pkg/charting`
- Passed full tests:
  - `go test ./...`
- Passed build:
  - `go build ./...`

### Related commits
- `2984f41` Add volume-price evidence engine with TA-Lib indicators
- `9938217` Add cmd/signal JSON report CLI
- `919b464` Add signal JSON schema and contract test
- `a2a5d9d` Refactor signal pipeline with configurable scoring and logging
