# colorkit

Agentic-first color manipulation service. Convert between hex, RGB, and HSL. Mix colors, generate shades/tints, complementary colors, contrast checking, gradient generation, and color palette creation. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
# Build
make build

# Run (defaults to :7100)
./colorkit

# Use it
curl "http://localhost:7100/convert?color=#ff0000"
# hex=#ff0000 rgb=255,0,0 hsl=0,100,50

curl "http://localhost:7100/mix?c1=#ff0000&c2=#0000ff&ratio=0.5"
# hex=#800080 rgb=128,0,128 hsl=270,100,25

curl "http://localhost:7100/contrast?c1=#000000&c2=#ffffff"
# ratio=21.00 wcag_aa=pass wcag_aaa=pass
```

## Principles

- **The agent IS the interface** — No UI, no SDK. The API is the product.
- **Plain text by default** — One labeled, grepable line per record. JSON on demand via `Accept: application/json` or `?format=json`.
- **Instructive errors** — Every 4xx includes a hint telling the agent what to do next.
- **Self-documenting** — `GET /help` returns a one-page operating manual.
- **Single static binary** — Go, CGO_ENABLED=0, zero external dependencies.
- **Zero config defaults** — Runs out of the box. Config: defaults < env < flags.
- **MCP connector** — Speaks Model Context Protocol at `POST /mcp`.

## API Reference

### Color Conversion

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/convert` | `color` | Convert between hex/RGB/HSL formats |

### Color Manipulation

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/mix` | `c1`, `c2`, `ratio` | Mix two colors (ratio 0.0-1.0, default 0.5) |
| GET | `/lighten` | `color`, `amount` | Lighten color (amount 0-100, default 10) |
| GET | `/darken` | `color`, `amount` | Darken color (amount 0-100, default 10) |
| GET | `/saturate` | `color`, `amount` | Increase saturation (amount 0-100, default 10) |
| GET | `/desaturate` | `color`, `amount` | Decrease saturation (amount 0-100, default 10) |
| GET | `/grayscale` | `color` | Convert to grayscale |
| GET | `/invert` | `color` | Invert color |
| GET | `/rotate` | `color`, `degrees` | Rotate hue (degrees, default 30) |

### Color Schemes

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/complement` | `color` | Complementary color (180° opposite) |
| GET | `/analogous` | `color`, `angle` | Two analogous colors (default 30°) |
| GET | `/triadic` | `color` | Two triadic colors (120° apart) |
| GET | `/tetradic` | `color` | Three tetradic colors (90° apart) |

### Color Analysis

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/contrast` | `c1`, `c2` | WCAG contrast ratio + AA/AAA pass/fail |
| GET | `/suggested-text` | `color` | Black or white for best text contrast |

### Color Generation

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/gradient` | `c1`, `c2`, `steps` | Gradient between two colors (default 5 steps) |
| GET | `/shades` | `color`, `count` | Darker shades (default 5) |
| GET | `/tints` | `color`, `count` | Lighter tints (default 5) |

### Color Formats

- **Hex**: `#ff0000`, `ff0000`, `#f00`, `f00`
- **RGB**: `255,0,0` (comma-separated, 0-255)
- **HSL**: `0,100,50` (comma-separated, H:0-360 S:0-100 L:0-100)

### MCP

`POST /mcp` — JSON-RPC 2.0 endpoint with 18 tools.

## Configuration

| Source | Key | Default | Description |
|--------|-----|---------|-------------|
| Flag | `-addr` | `:7100` | Listen address |
| Flag | `-secret` | (random) | Auth token signing secret |
| Env | `COLORKIT_ADDR` | `:7100` | Listen address |
| Env | `COLORKIT_SECRET` | (random) | Auth token signing secret |

## Build

```bash
make build    # CGO_ENABLED=0 go build -trimpath ./cmd/colorkit
make test     # go test -race ./...
make vet      # go vet ./...
```

No database needed. Pure stateless computation. Single Go binary.

## License

MIT
