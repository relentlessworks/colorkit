package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB represents a color in red-green-blue space (0-255 each).
type RGB struct {
	R, G, B int
}

// HSL represents a color in hue-saturation-lightness space.
// H: 0-360, S: 0-100, L: 0-100
type HSL struct {
	H, S, L float64
}

// Color is the full color representation with all formats.
type Color struct {
	Hex string
	RGB RGB
	HSL HSL
}

// ParseHex parses a hex color string (with or without #).
// Supports #rgb, #rrggbb, and #rrggbbaa formats.
func ParseHex(s string) (RGB, error) {
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		r, err := strconv.ParseInt(string(s[0])+string(s[0]), 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		g, err := strconv.ParseInt(string(s[1])+string(s[1]), 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		b, err := strconv.ParseInt(string(s[2])+string(s[2]), 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		return RGB{R: int(r), G: int(g), B: int(b)}, nil
	case 6:
		r, err := strconv.ParseInt(s[0:2], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		g, err := strconv.ParseInt(s[2:4], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		b, err := strconv.ParseInt(s[4:6], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		return RGB{R: int(r), G: int(g), B: int(b)}, nil
	case 8:
		r, err := strconv.ParseInt(s[0:2], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		g, err := strconv.ParseInt(s[2:4], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		b, err := strconv.ParseInt(s[4:6], 16, 0)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color")
		}
		return RGB{R: int(r), G: int(g), B: int(b)}, nil
	default:
		return RGB{}, fmt.Errorf("invalid hex color")
	}
}

// ParseRGB parses an "r,g,b" string into RGB.
func ParseRGB(s string) (RGB, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return RGB{}, fmt.Errorf("invalid rgb format, expected r,g,b")
	}
	r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || r < 0 || r > 255 {
		return RGB{}, fmt.Errorf("invalid red value, must be 0-255")
	}
	g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || g < 0 || g > 255 {
		return RGB{}, fmt.Errorf("invalid green value, must be 0-255")
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil || b < 0 || b > 255 {
		return RGB{}, fmt.Errorf("invalid blue value, must be 0-255")
	}
	return RGB{R: r, G: g, B: b}, nil
}

// ParseHSL parses an "h,s,l" string into HSL.
func ParseHSL(s string) (HSL, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return HSL{}, fmt.Errorf("invalid hsl format, expected h,s,l")
	}
	h, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil || h < 0 || h > 360 {
		return HSL{}, fmt.Errorf("invalid hue value, must be 0-360")
	}
	sat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil || sat < 0 || sat > 100 {
		return HSL{}, fmt.Errorf("invalid saturation value, must be 0-100")
	}
	l, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil || l < 0 || l > 100 {
		return HSL{}, fmt.Errorf("invalid lightness value, must be 0-100")
	}
	return HSL{H: h, S: sat, L: l}, nil
}

// RGBToHex converts RGB to a hex string (#rrggbb).
func RGBToHex(c RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// RGBToHSL converts RGB to HSL.
func RGBToHSL(c RGB) HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2.0

	var h, s float64
	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case r:
			if g < b {
				h = (g-b)/d + 6
			} else {
				h = (g - b) / d
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return HSL{H: h, S: s * 100, L: l * 100}
}

// HSLToRGB converts HSL to RGB.
func HSLToRGB(c HSL) RGB {
	h := c.H / 360.0
	s := c.S / 100.0
	l := c.L / 100.0

	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	return RGB{
		R: clampInt(r * 255),
		G: clampInt(g * 255),
		B: clampInt(b * 255),
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func clampInt(f float64) int {
	v := int(math.Round(f))
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

// FullColor creates a Color from RGB with all formats populated.
func FullColor(rgb RGB) Color {
	return Color{
		Hex: RGBToHex(rgb),
		RGB: rgb,
		HSL: RGBToHSL(rgb),
	}
}

// Mix blends two colors with a ratio (0.0 = color1, 1.0 = color2).
func Mix(c1, c2 RGB, ratio float64) RGB {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return RGB{
		R: clampInt(float64(c1.R) + (float64(c2.R)-float64(c1.R))*ratio),
		G: clampInt(float64(c1.G) + (float64(c2.G)-float64(c1.G))*ratio),
		B: clampInt(float64(c1.B) + (float64(c2.B)-float64(c1.B))*ratio),
	}
}

// Lighten increases the lightness of a color by the given percentage.
func Lighten(rgb RGB, amount float64) RGB {
	hsl := RGBToHSL(rgb)
	hsl.L += amount
	if hsl.L > 100 {
		hsl.L = 100
	}
	return HSLToRGB(hsl)
}

// Darken decreases the lightness of a color by the given percentage.
func Darken(rgb RGB, amount float64) RGB {
	hsl := RGBToHSL(rgb)
	hsl.L -= amount
	if hsl.L < 0 {
		hsl.L = 0
	}
	return HSLToRGB(hsl)
}

// Saturate increases the saturation of a color by the given percentage.
func Saturate(rgb RGB, amount float64) RGB {
	hsl := RGBToHSL(rgb)
	hsl.S += amount
	if hsl.S > 100 {
		hsl.S = 100
	}
	return HSLToRGB(hsl)
}

// Desaturate decreases the saturation of a color by the given percentage.
func Desaturate(rgb RGB, amount float64) RGB {
	hsl := RGBToHSL(rgb)
	hsl.S -= amount
	if hsl.S < 0 {
		hsl.S = 0
	}
	return HSLToRGB(hsl)
}

// Grayscale converts a color to grayscale (removes saturation).
func Grayscale(rgb RGB) RGB {
	gray := int(0.299*float64(rgb.R) + 0.587*float64(rgb.G) + 0.114*float64(rgb.B))
	return RGB{R: gray, G: gray, B: gray}
}

// Complement returns the complementary color (180 degrees opposite on the color wheel).
func Complement(rgb RGB) RGB {
	hsl := RGBToHSL(rgb)
	hsl.H = math.Mod(hsl.H+180, 360)
	return HSLToRGB(hsl)
}

// Analogous returns two colors adjacent to the given color on the color wheel.
func Analogous(rgb RGB, angle float64) (RGB, RGB) {
	hsl := RGBToHSL(rgb)
	h1 := math.Mod(hsl.H+angle, 360)
	h2 := math.Mod(hsl.H-angle+360, 360)
	return HSLToRGB(HSL{H: h1, S: hsl.S, L: hsl.L}), HSLToRGB(HSL{H: h2, S: hsl.S, L: hsl.L})
}

// Triadic returns two colors that form a triadic color scheme (120 degrees apart).
func Triadic(rgb RGB) (RGB, RGB) {
	hsl := RGBToHSL(rgb)
	h1 := math.Mod(hsl.H+120, 360)
	h2 := math.Mod(hsl.H+240, 360)
	return HSLToRGB(HSL{H: h1, S: hsl.S, L: hsl.L}), HSLToRGB(HSL{H: h2, S: hsl.S, L: hsl.L})
}

// Tetradic returns three colors that form a tetradic (rectangle) color scheme.
func Tetradic(rgb RGB) (RGB, RGB, RGB) {
	hsl := RGBToHSL(rgb)
	h1 := math.Mod(hsl.H+90, 360)
	h2 := math.Mod(hsl.H+180, 360)
	h3 := math.Mod(hsl.H+270, 360)
	return HSLToRGB(HSL{H: h1, S: hsl.S, L: hsl.L}),
		HSLToRGB(HSL{H: h2, S: hsl.S, L: hsl.L}),
		HSLToRGB(HSL{H: h3, S: hsl.S, L: hsl.L})
}

// RelativeLuminance calculates the WCAG relative luminance of a color.
func RelativeLuminance(rgb RGB) float64 {
	r := lin(float64(rgb.R) / 255.0)
	g := lin(float64(rgb.G) / 255.0)
	b := lin(float64(rgb.B) / 255.0)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func lin(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// ContrastRatio calculates the WCAG contrast ratio between two colors.
func ContrastRatio(c1, c2 RGB) float64 {
	l1 := RelativeLuminance(c1)
	l2 := RelativeLuminance(c2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// Gradient generates n evenly spaced colors between c1 and c2.
func Gradient(c1, c2 RGB, n int) []RGB {
	if n < 2 {
		return []RGB{c1}
	}
	result := make([]RGB, n)
	for i := 0; i < n; i++ {
		ratio := float64(i) / float64(n-1)
		result[i] = Mix(c1, c2, ratio)
	}
	return result
}

// Shades generates n shades (darker versions) of a color.
func Shades(rgb RGB, n int) []RGB {
	result := make([]RGB, n)
	for i := 0; i < n; i++ {
		amount := float64(i+1) / float64(n+1) * 100
		result[i] = Darken(rgb, amount)
	}
	return result
}

// Tints generates n tints (lighter versions) of a color.
func Tints(rgb RGB, n int) []RGB {
	result := make([]RGB, n)
	for i := 0; i < n; i++ {
		amount := float64(i+1) / float64(n+1) * 100
		result[i] = Lighten(rgb, amount)
	}
	return result
}

// Invert inverts a color (subtracts each channel from 255).
func Invert(rgb RGB) RGB {
	return RGB{R: 255 - rgb.R, G: 255 - rgb.G, B: 255 - rgb.B}
}

// RotateHue rotates the hue of a color by the given degrees.
func RotateHue(rgb RGB, degrees float64) RGB {
	hsl := RGBToHSL(rgb)
	hsl.H = math.Mod(hsl.H+degrees+360, 360)
	return HSLToRGB(hsl)
}

// IsLight returns true if the color is considered light (luminance > 0.5).
func IsLight(rgb RGB) bool {
	return RelativeLuminance(rgb) > 0.5
}

// IsDark returns true if the color is considered dark (luminance <= 0.5).
func IsDark(rgb RGB) bool {
	return !IsLight(rgb)
}

// SuggestedTextColor returns black or white, whichever has better contrast with the given color.
func SuggestedTextColor(rgb RGB) RGB {
	if IsLight(rgb) {
		return RGB{R: 0, G: 0, B: 0}
	}
	return RGB{R: 255, G: 255, B: 255}
}

// FormatRGB returns "r,g,b" string.
func FormatRGB(c RGB) string {
	return fmt.Sprintf("%d,%d,%d", c.R, c.G, c.B)
}

// FormatHSL returns "h,s,l" string (rounded to 1 decimal).
func FormatHSL(c HSL) string {
	return fmt.Sprintf("%.1f,%.1f,%.1f", c.H, c.S, c.L)
}
