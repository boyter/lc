# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`lc` (licensechecker) is a CLI tool that recursively scans directories to identify software licenses in files. It uses SPDX license identification with multiple detection strategies: SPDX header parsing, filename matching, and content analysis (keyword matching, vector space, Levenshtein distance). Written in Go, currently at v2.0.0 alpha.

## Build & Development Commands

```bash
go build                        # Build binary
go test -cover -race ./...      # Run all tests with coverage and race detection
go test -v -run TestName ./processor/  # Run a single test
gofmt -s -w ./..                # Format code
golangci-lint run --enable gofmt ./...  # Lint
./check.sh                      # Full verification (fmt, test, lint, race, cross-compile)
```

### Regenerating the License Database

```bash
./generate_database.sh          # Build DB, copy JSON, run go generate, test
```

This builds `assets/database/`, produces `database_keywords.json`, then `go generate` (via `scripts/include.go`) embeds it as base64 in `processor/constants.go`.

## Architecture

**Entry point:** `main.go` — Cobra CLI that creates `processor.NewProcess(".")` and calls `StartProcess()`.

**`processor/` package** (active v2 code):
- `processor.go` — Orchestrator: walks files via `gocodewalker`, reads content (max 100KB), routes through detection pipeline
- `detector_spdx.go` — Parses `SPDX-License-Identifier:` headers from source files (100% confidence)
- `detector_license.go` — Detects licenses in dedicated license files (LICENSE, COPYING, etc.) using filename regex matching
- `guesser.go` — `LicenceGuesser` interface and framework; two instances: common licenses (fast path) and full database
- `guesser_keyword.go` — Keyword-based license detection using Aho-Corasick
- `guesser_vectorspace.go` — TF-IDF vector space similarity matching
- `guesser_blended.go` — Combines keyword + vector space scores
- `constants.go` — Auto-generated; contains base64-encoded license database (do not edit manually)
- `structs.go` — Core data types (`FileResult`, `LicenseMatch`, etc.)
- `common.go` — Shared utilities and compiled regex patterns for license filename detection

**`parsers/` and `pkg/`** — Legacy v1 code, deprecated and scheduled for removal.

**`assets/database/`** — Database builder that processes 425+ SPDX license definition files into `database_keywords.json`.

## Detection Pipeline

1. Check if file is binary (null byte detection) — skip unless `--binary` flag
2. For files matching license filename patterns (license, copying, mit, apache, etc.): run through `LicenceGuesser` (keyword → vector space → blended)
3. For all other files: scan for `SPDX-License-Identifier:` headers
