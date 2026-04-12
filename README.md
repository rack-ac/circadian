# circadian

Circadian is a small MCP server for the WHOOP Developer API.

I built it partly because I wanted a WHOOP integration for my own MCP setup, and partly because I wanted to properly learn how MCP servers work in practice.

It is intentionally narrow in scope. The goal is not to mirror the full WHOOP API. The goal is to provide a small, usable local server that exposes the pieces of WHOOP data that are most useful in day-to-day MCP workflows.

## What It Does

Circadian exposes WHOOP data over MCP streamable HTTP.

Right now it includes tools for:

- profile
- body measurements
- cycles
- recoveries
- sleeps
- workouts

It also includes a few lifecycle tools so you can inspect auth state, reauthorize, and revoke access without dealing with token files manually.

## How It Runs

The server runs locally and exposes an MCP endpoint like:

```text
http://127.0.0.1:8080/mcp
```

Your MCP client connects to that endpoint.

When a WHOOP-protected tool is used for the first time, the server opens a browser window so you can authorize against the WHOOP Developer API. Tokens are then stored locally at:

```text
~/.circadian/tokens.json
```

The server handles refresh tokens automatically when it can.

## Current Behavior

A few implementation details are worth knowing:

- responses are returned as a small `summary` plus the raw upstream WHOOP payload under `data`
- `401` responses trigger token renewal and retry
- `429` and transient upstream failures are retried with backoff
- a small local request interval helps avoid accidental bursts against WHOOP
- token files are written with `0600` permissions
- MCP request logging is written to stderr

This is an MCP endpoint, not a browser UI. If you hit `/mcp` directly with a plain `GET`, that is not a health page.

## WHOOP Setup

Create a WHOOP developer application at <https://developer.whoop.com/api/>.

Unless you override it, use this redirect URI in the WHOOP app:

```text
http://127.0.0.1:8976/callback
```

## Configuration

Required:

```bash
export WHOOP_CLIENT_ID=your_client_id
export WHOOP_CLIENT_SECRET=your_client_secret
```

Optional:

```bash
export WHOOP_REDIRECT_URI=http://127.0.0.1:8976/callback
export WHOOP_API_ROOT=https://api.prod.whoop.com
export WHOOP_MIN_INTERVAL_MS=350
export CIRCADIAN_HOST=127.0.0.1
export CIRCADIAN_PORT=8080
export CIRCADIAN_PATH=/mcp
```

Defaults:

- MCP endpoint: `http://127.0.0.1:8080/mcp`
- OAuth callback: `http://127.0.0.1:8976/callback`
- token file: `~/.circadian/tokens.json`

## Running It

Run it directly with Go:

```bash
go run ./cmd/circadian
```

Or build a local binary:

```bash
make build
./circadian
```

## Connecting a Client

Example MCP client config:

```json
{
  "mcpServers": {
    "whoop": {
      "url": "http://127.0.0.1:8080/mcp"
    }
  }
}
```

If you change `CIRCADIAN_HOST`, `CIRCADIAN_PORT`, or `CIRCADIAN_PATH`, update the URL to match.

## Tool Inputs

Collection-style tools accept:

- `limit`
- `start`
- `end`
- `next_token`

Single-record lookups:

- `whoop_cycles` accepts `cycle_id`
- `whoop_recoveries` accepts `cycle_id`
- `whoop_sleeps` accepts `sleep_id`
- `whoop_workouts` accepts `workout_id`

The server validates a few things locally before calling WHOOP:

- `limit` must be between `0` and `25`
- `start` and `end` must be RFC3339 timestamps
- single-record IDs cannot be mixed with collection filters

## Example Response Shape

```json
{
  "summary": {
    "type": "sleep",
    "sleep_id": "uuid",
    "duration": "7h52m0s",
    "sleep_efficiency_percentage": 94.2
  },
  "data": {
    "...": "raw WHOOP response"
  }
}
```

The `summary` field is there to make the MCP output easier to read quickly. The `data` field keeps the full upstream payload available.

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build ./...
```

Build release binaries locally:

```bash
make dist
```

Current release output is macOS only:

- `circadian-darwin-arm64`
- `circadian-darwin-amd64`

## Automation

CI currently runs:

- tests
- build checks
- `gofmt`
- `govulncheck`
- `gosec`
- commitlint on pull requests

Dependency updates are handled through Renovate.

Releases are driven by conventional commits using `release-please`, so versioning is based on the actual commit history rather than manually choosing versions each time.
