package api

const helpText = `colorkit — agentic-first color manipulation service

Convert between hex/RGB/HSL, mix colors, generate shades/tints, complementary colors,
contrast checking, gradient generation, and color palette creation. No database needed.

== ENDPOINTS ==

GET  /convert?color=<hex|rgb>
  Convert a color between formats. Returns hex, rgb, hsl.
  Example: GET /convert?color=#ff0000
  Response: hex=#ff0000 rgb=255,0,0 hsl=0,100,50

GET  /mix?c1=<color>&c2=<color>&ratio=<0-1>
  Mix two colors with a ratio (0.0 = c1, 1.0 = c2, default 0.5).
  Example: GET /mix?c1=#000000&c2=#ffffff&ratio=0.5

GET  /lighten?color=<color>&amount=<0-100>
  Lighten a color by increasing lightness (default 10%).
  Example: GET /lighten?color=#336699&amount=20

GET  /darken?color=<color>&amount=<0-100>
  Darken a color by decreasing lightness (default 10%).
  Example: GET /darken?color=#336699&amount=20

GET  /saturate?color=<color>&amount=<0-100>
  Increase saturation (default 10%).
  Example: GET /saturate?color=#808080&amount=50

GET  /desaturate?color=<color>&amount=<0-100>
  Decrease saturation (default 10%).
  Example: GET /desaturate?color=#ff0000&amount=50

GET  /grayscale?color=<color>
  Convert a color to grayscale using luminance weights.
  Example: GET /grayscale?color=#ff0000

GET  /invert?color=<color>
  Invert a color (subtract each channel from 255).
  Example: GET /invert?color=#ff8800

GET  /complement?color=<color>
  Get the complementary color (180 degrees opposite on color wheel).
  Example: GET /complement?color=#ff0000

GET  /analogous?color=<color>&angle=<degrees>
  Get two analogous colors (adjacent on color wheel, default 30 degrees).
  Example: GET /analogous?color=#ff0000&angle=30

GET  /triadic?color=<color>
  Get two triadic colors (120 degrees apart on color wheel).
  Example: GET /triadic?color=#ff0000

GET  /tetradic?color=<color>
  Get three tetradic colors (90 degrees apart, forming a rectangle).
  Example: GET /tetradic?color=#ff0000

GET  /contrast?c1=<color>&c2=<color>
  Calculate WCAG contrast ratio between two colors.
  Example: GET /contrast?c1=#000000&c2=#ffffff
  Response: ratio=21.00 wcag_aa=pass wcag_aaa=pass

GET  /gradient?c1=<color>&c2=<color>&steps=<n>
  Generate n evenly spaced colors between c1 and c2 (default 5).
  Example: GET /gradient?c1=#ff0000&c2=#0000ff&steps=10

GET  /shades?color=<color>&count=<n>
  Generate n darker shades of a color (default 5).
  Example: GET /shades?color=#336699&count=5

GET  /tints?color=<color>&count=<n>
  Generate n lighter tints of a color (default 5).
  Example: GET /tints?color=#336699&count=5

GET  /rotate?color=<color>&degrees=<n>
  Rotate the hue of a color by degrees (default 30).
  Example: GET /rotate?color=#ff0000&degrees=120

GET  /suggested-text?color=<color>
  Get black or white, whichever has better contrast with the given color.
  Example: GET /suggested-text?color=#336699

POST /mcp
  MCP (Model Context Protocol) JSON-RPC 2.0 endpoint.
  Methods: initialize, tools/list, tools/call
  Tools: convert, mix, lighten, darken, saturate, desaturate, grayscale,
         invert, complement, analogous, triadic, tetradic, contrast,
         gradient, shades, tints, rotate, suggested_text

== COLOR FORMATS ==

Hex: #ff0000 or ff0000 (3 or 6 digit, with or without #)
RGB: 255,0,0 (comma-separated, 0-255 each)
HSL: 0,100,50 (comma-separated, H:0-360 S:0-100 L:0-100)

== RESPONSES ==

Plain text by default: one labeled line per color.
  hex=#ff0000 rgb=255,0,0 hsl=0,100,50

JSON on demand: send Accept: application/json or add ?format=json
  {"Hex":"#ff0000","RGB":{"R":255,"G":0,"B":0},"HSL":{"H":0,"S":100,"L":50}}

Errors are instructive:
  error: missing color parameter | hint: provide a color as hex (#ff0000) or rgb (255,0,0)

== CONFIG ==

Flags: -addr=:7100 -secret=<token>
Env:   COLORKIT_ADDR=:7100 COLORKIT_SECRET=<token>

No database needed. Pure stateless computation. Single Go binary.
`
