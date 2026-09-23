# Agentic App Studio

Agentic App Studio is a local-first Studio for turning business intent and optional evidence into a Product Graph, executable Cards, agent sessions, verification, and a Ledger.

## Local Development

Initial requirements:

- Go 1.26.x or the version pinned in `go.mod`.
- Node.js and npm as specified in `apps/studio/frontend/package.json`.
- Wails v2.15.0 is the desktop runtime pinned in `go.mod`; the CLI can be installed locally or available in `PATH`.

Verification:

```text
make verify
```

The first stage does not run agents or commands against external repositories. It only validates the local runtime, contracts, and Studio shell.

To launch the native host during development, run `wails dev` from `apps/studio` after installing the Wails CLI.

## Principles

- The Product Graph and contracts are independent of any agent.
- The Studio is local-first, and its SQLite database is the source of truth for its own projects.
- Desktop Primary will be the only initial interface for writing to generated products.
- Agents only execute authorized Task Contracts.
- Every significant result includes provenance and evidence.
