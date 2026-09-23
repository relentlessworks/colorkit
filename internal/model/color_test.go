package model

import "testing"

func TestParseHex(t *testing.T) {
	tests := []struct {
		in   string
		want RGB
	}{
		{"#ff0000", RGB{255, 0, 0}},
		{"ff0000", RGB{255, 0, 0}},
		{"#f00", RGB{255, 0, 0}},
		{"f00", RGB{255, 0, 0}},
		{"#00ff00", RGB{0, 255, 0}},
		{"#0000ff", RGB{0, 0, 255}},
		{"#ffffff", RGB{255, 255, 255}},
		{"#000000", RGB{0, 0, 0}},
	}
	for _, tt := range tests {
		got, err := ParseHex(tt.in)
		if err != nil {
			t.Errorf("ParseHex(%q) error: %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseHex(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseHexInvalid(t *testing.T) {
	invalid := []string{"", "#", "#ff", "#ffff", "#fffffff", "#gggggg", "xyz"}
	for _, s := range invalid {
		_, err := ParseHex(s)
		if err == nil {
			t.Errorf("ParseHex(%q) expected error, got nil", s)
		}
	}
}

func TestParseRGB(t *testing.T) {
	tests := []struct {
		in   string
		want RGB
	}{
		{"255,0,0", RGB{255, 0, 0}},
		{"0,255,0", RGB{0, 255, 0}},
		{"0, 0, 255", RGB{0, 0, 255}},
		{"128,128,128", RGB{128, 128, 128}},
	}
	for _, tt := range tests {
		got, err := ParseRGB(tt.in)
		if err != nil {
			t.Errorf("ParseRGB(%q) error: %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseRGB(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseRGBInvalid(t *testing.T) {
	invalid := []string{"255,0", "256,0,0", "0,-1,0", "0,0,300", "a,b,c", ""}
	for _, s := range invalid {
		_, err := ParseRGB(s)
		if err == nil {
			t.Errorf("ParseRGB(%q) expected error, got nil", s)
		}
	}
}

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		in   RGB
		want string
	}{
		{RGB{255, 0, 0}, "#ff0000"},
		{RGB{0, 255, 0}, "#00ff00"},
		{RGB{0, 0, 255}, "#0000ff"},
		{RGB{255, 255, 255}, "#ffffff"},
		{RGB{0, 0, 0}, "#000000"},
		{RGB{128, 64, 32}, "#804020"},
	}
	for _, tt := range tests {
		got := RGBToHex(tt.in)
		if got != tt.want {
			t.Errorf("RGBToHex(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRGBToHSL(t *testing.T) {
	tests := []struct {
		in      RGB
		wantH   float64
		wantS   float64
		wantL   float64
		epsilon float64
	}{
		{RGB{255, 0, 0}, 0, 100, 50, 0.1},
		{RGB{0, 255, 0}, 120, 100, 50, 0.1},
		{RGB{0, 0, 255}, 240, 100, 50, 0.1},
		{RGB{128, 128, 128}, 0, 0, 50.2, 0.5},
		{RGB{255, 255, 255}, 0, 0, 100, 0.1},
		{RGB{0, 0, 0}, 0, 0, 0, 0.1},
	}
	for _, tt := range tests {
		got := RGBToHSL(tt.in)
		if abs(got.H-tt.wantH) > tt.epsilon || abs(got.S-tt.wantS) > tt.epsilon || abs(got.L-tt.wantL) > tt.epsilon {
			t.Errorf("RGBToHSL(%v) = %.2f,%.2f,%.2f, want %.2f,%.2f,%.2f", tt.in, got.H, got.S, got.L, tt.wantH, tt.wantS, tt.wantL)
		}
	}
}

func TestHSLToRGB(t *testing.T) {
	tests := []struct {
		in   HSL
		want RGB
	}{
		{HSL{0, 100, 50}, RGB{255, 0, 0}},
		{HSL{120, 100, 50}, RGB{0, 255, 0}},
		{HSL{240, 100, 50}, RGB{0, 0, 255}},
		{HSL{0, 0, 0}, RGB{0, 0, 0}},
		{HSL{0, 0, 100}, RGB{255, 255, 255}},
		{HSL{0, 0, 50}, RGB{128, 128, 128}},
	}
	for _, tt := range tests {
		got := HSLToRGB(tt.in)
		if abs(float64(got.R-tt.want.R)) > 1 || abs(float64(got.G-tt.want.G)) > 1 || abs(float64(got.B-tt.want.B)) > 1 {
			t.Errorf("HSLToRGB(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestMix(t *testing.T) {
	c1 := RGB{0, 0, 0}
	c2 := RGB{255, 255, 255}

	got := Mix(c1, c2, 0)
	if got != c1 {
		t.Errorf("Mix ratio=0 = %v, want %v", got, c1)
	}
	got = Mix(c1, c2, 1)
	if got != c2 {
		t.Errorf("Mix ratio=1 = %v, want %v", got, c2)
	}
	got = Mix(c1, c2, 0.5)
	if got != (RGB{128, 128, 128}) {
		t.Errorf("Mix ratio=0.5 = %v, want {128,128,128}", got)
	}
}

func TestLightenDarken(t *testing.T) {
	rgb := RGB{100, 100, 100}
	lighter := Lighten(rgb, 20)
	if lighter == rgb {
		t.Error("Lighten did not change color")
	}
	hsl := RGBToHSL(lighter)
	if hsl.L <= 39 {
		t.Errorf("Lighten should increase lightness, got %.2f", hsl.L)
	}
	darker := Darken(rgb, 20)
	if darker == rgb {
		t.Error("Darken did not change color")
	}
	hsl = RGBToHSL(darker)
	if hsl.L >= 39 {
		t.Errorf("Darken should decrease lightness, got %.2f", hsl.L)
	}
}

func TestComplement(t *testing.T) {
	red := RGB{255, 0, 0}
	comp := Complement(red)
	hsl := RGBToHSL(comp)
	if abs(hsl.H-180) > 1 {
		t.Errorf("Complement of red should be ~180 hue, got %.2f", hsl.H)
	}
}

func TestContrastRatio(t *testing.T) {
	ratio := ContrastRatio(RGB{0, 0, 0}, RGB{255, 255, 255})
	if abs(ratio-21) > 0.5 {
		t.Errorf("Black/white contrast should be ~21, got %.2f", ratio)
	}
	ratio = ContrastRatio(RGB{255, 255, 255}, RGB{255, 255, 255})
	if abs(ratio-1) > 0.1 {
		t.Errorf("Same color contrast should be 1, got %.2f", ratio)
	}
}

func TestGradient(t *testing.T) {
	c1 := RGB{0, 0, 0}
	c2 := RGB{255, 255, 255}
	grad := Gradient(c1, c2, 5)
	if len(grad) != 5 {
		t.Fatalf("Gradient len = %d, want 5", len(grad))
	}
	if grad[0] != c1 {
		t.Errorf("Gradient[0] = %v, want %v", grad[0], c1)
	}
	if grad[4] != c2 {
		t.Errorf("Gradient[4] = %v, want %v", grad[4], c2)
	}
}

func TestInvert(t *testing.T) {
	got := Invert(RGB{100, 200, 50})
	want := RGB{155, 55, 205}
	if got != want {
		t.Errorf("Invert = %v, want %v", got, want)
	}
}

func TestGrayscale(t *testing.T) {
	got := Grayscale(RGB{255, 0, 0})
	if got.R != got.G || got.G != got.B {
		t.Errorf("Grayscale should have equal channels, got %v", got)
	}
}

func TestSuggestedTextColor(t *testing.T) {
	if tc := SuggestedTextColor(RGB{255, 255, 255}); tc != (RGB{0, 0, 0}) {
		t.Errorf("SuggestedTextColor(white) = %v, want black", tc)
	}
	if tc := SuggestedTextColor(RGB{0, 0, 0}); tc != (RGB{255, 255, 255}) {
		t.Errorf("SuggestedTextColor(black) = %v, want white", tc)
	}
}

func TestTriadic(t *testing.T) {
	red := RGB{255, 0, 0}
	c1, c2 := Triadic(red)
	h1 := RGBToHSL(c1)
	h2 := RGBToHSL(c2)
	if abs(h1.H-120) > 2 {
		t.Errorf("Triadic[0] hue = %.2f, want ~120", h1.H)
	}
	if abs(h2.H-240) > 2 {
		t.Errorf("Triadic[1] hue = %.2f, want ~240", h2.H)
	}
}

func TestAnalogous(t *testing.T) {
	red := RGB{255, 0, 0}
	c1, c2 := Analogous(red, 30)
	h1 := RGBToHSL(c1)
	h2 := RGBToHSL(c2)
	if abs(h1.H-30) > 2 {
		t.Errorf("Analogous[0] hue = %.2f, want ~30", h1.H)
	}
	if abs(h2.H-330) > 2 {
		t.Errorf("Analogous[1] hue = %.2f, want ~330", h2.H)
	}
}

func TestRotateHue(t *testing.T) {
	red := RGB{255, 0, 0}
	rotated := RotateHue(red, 120)
	h := RGBToHSL(rotated)
	if abs(h.H-120) > 2 {
		t.Errorf("RotateHue(120) hue = %.2f, want ~120", h.H)
	}
}

func TestShadesTints(t *testing.T) {
	rgb := RGB{128, 64, 32}
	shades := Shades(rgb, 3)
	if len(shades) != 3 {
		t.Fatalf("Shades len = %d, want 3", len(shades))
	}
	for i, s := range shades {
		hsl := RGBToHSL(s)
		origHSL := RGBToHSL(rgb)
		if hsl.L >= origHSL.L {
			t.Errorf("Shade[%d] should be darker, L=%.2f, orig L=%.2f", i, hsl.L, origHSL.L)
		}
	}
	tints := Tints(rgb, 3)
	if len(tints) != 3 {
		t.Fatalf("Tints len = %d, want 3", len(tints))
	}
	for i, ti := range tints {
		hsl := RGBToHSL(ti)
		origHSL := RGBToHSL(rgb)
		if hsl.L <= origHSL.L {
			t.Errorf("Tint[%d] should be lighter, L=%.2f, orig L=%.2f", i, hsl.L, origHSL.L)
		}
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
