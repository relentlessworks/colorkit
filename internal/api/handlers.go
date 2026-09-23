package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/colorkit/internal/model"
)

// Handler holds all HTTP handlers for the colorkit service.
type Handler struct {
	secret string
}

// New creates a new API handler.
func New(secret string) *Handler {
	return &Handler{secret: secret}
}

// Routes returns the HTTP mux with all routes registered.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.help)
	mux.HandleFunc("/convert", h.convert)
	mux.HandleFunc("/mix", h.mix)
	mux.HandleFunc("/lighten", h.lighten)
	mux.HandleFunc("/darken", h.darken)
	mux.HandleFunc("/saturate", h.saturate)
	mux.HandleFunc("/desaturate", h.desaturate)
	mux.HandleFunc("/grayscale", h.grayscale)
	mux.HandleFunc("/invert", h.invert)
	mux.HandleFunc("/complement", h.complement)
	mux.HandleFunc("/analogous", h.analogous)
	mux.HandleFunc("/triadic", h.triadic)
	mux.HandleFunc("/tetradic", h.tetradic)
	mux.HandleFunc("/contrast", h.contrast)
	mux.HandleFunc("/gradient", h.gradient)
	mux.HandleFunc("/shades", h.shades)
	mux.HandleFunc("/tints", h.tints)
	mux.HandleFunc("/rotate", h.rotate)
	mux.HandleFunc("/suggested-text", h.suggestedText)
	mux.HandleFunc("/mcp", h.mcp)
	return mux
}

// wantsJSON checks if the client wants JSON output.
func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

// writeError writes an instructive error response.
func writeError(w http.ResponseWriter, status int, msg, hint string) {
	w.WriteHeader(status)
	if hint != "" {
		fmt.Fprintf(w, "error: %s | hint: %s\n", msg, hint)
	} else {
		fmt.Fprintf(w, "error: %s\n", msg)
	}
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// writeColor writes a color in plain text or JSON format.
func writeColor(w http.ResponseWriter, r *http.Request, rgb model.RGB) {
	if wantsJSON(r) {
		c := model.FullColor(rgb)
		writeJSON(w, c)
		return
	}
	hsl := model.RGBToHSL(rgb)
	fmt.Fprintf(w, "hex=%s rgb=%s hsl=%s\n", model.RGBToHex(rgb), model.FormatRGB(rgb), model.FormatHSL(hsl))
}

// writeColorList writes a list of colors.
func writeColorList(w http.ResponseWriter, r *http.Request, colors []model.RGB, prefix string) {
	if wantsJSON(r) {
		result := make([]model.Color, len(colors))
		for i, c := range colors {
			result[i] = model.FullColor(c)
		}
		writeJSON(w, result)
		return
	}
	for i, c := range colors {
		hsl := model.RGBToHSL(c)
		fmt.Fprintf(w, "%s[%d] hex=%s rgb=%s hsl=%s\n", prefix, i, model.RGBToHex(c), model.FormatRGB(c), model.FormatHSL(hsl))
	}
}

// parseColorParam extracts a color from query params or form body.
func parseColorParam(r *http.Request, key string) (model.RGB, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		val = r.FormValue(key)
	}
	if val == "" {
		return model.RGB{}, fmt.Errorf("missing %s parameter", key)
	}
	return parseColorString(val)
}

// parseColorString parses a color from hex or rgb string.
func parseColorString(s string) (model.RGB, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") || isHexString(s) {
		return model.ParseHex(s)
	}
	if strings.Contains(s, ",") {
		return model.ParseRGB(s)
	}
	return model.ParseHex(s)
}

func isHexString(s string) bool {
	if len(s) == 0 || len(s) > 8 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// --- Handlers ---

func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, helpText)
}

func (h *Handler) convert(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex (#ff0000 or ff0000) or rgb (255,0,0)")
		return
	}
	writeColor(w, r, rgb)
}

func (h *Handler) mix(w http.ResponseWriter, r *http.Request) {
	c1, err := parseColorParam(r, "c1")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c1 as hex or rgb (e.g. c1=#ff0000)")
		return
	}
	c2, err := parseColorParam(r, "c2")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c2 as hex or rgb (e.g. c2=#0000ff)")
		return
	}
	ratio := 0.5
	if v := r.URL.Query().Get("ratio"); v != "" {
		if f, e := parseFloat(v); e == nil {
			ratio = f
		}
	}
	result := model.Mix(c1, c2, ratio)
	writeColor(w, r, result)
}

func (h *Handler) lighten(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	amount := 10.0
	if v := r.URL.Query().Get("amount"); v != "" {
		if f, e := parseFloat(v); e == nil {
			amount = f
		}
	}
	writeColor(w, r, model.Lighten(rgb, amount))
}

func (h *Handler) darken(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	amount := 10.0
	if v := r.URL.Query().Get("amount"); v != "" {
		if f, e := parseFloat(v); e == nil {
			amount = f
		}
	}
	writeColor(w, r, model.Darken(rgb, amount))
}

func (h *Handler) saturate(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	amount := 10.0
	if v := r.URL.Query().Get("amount"); v != "" {
		if f, e := parseFloat(v); e == nil {
			amount = f
		}
	}
	writeColor(w, r, model.Saturate(rgb, amount))
}

func (h *Handler) desaturate(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	amount := 10.0
	if v := r.URL.Query().Get("amount"); v != "" {
		if f, e := parseFloat(v); e == nil {
			amount = f
		}
	}
	writeColor(w, r, model.Desaturate(rgb, amount))
}

func (h *Handler) grayscale(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	writeColor(w, r, model.Grayscale(rgb))
}

func (h *Handler) invert(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	writeColor(w, r, model.Invert(rgb))
}

func (h *Handler) complement(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	writeColor(w, r, model.Complement(rgb))
}

func (h *Handler) analogous(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	angle := 30.0
	if v := r.URL.Query().Get("angle"); v != "" {
		if f, e := parseFloat(v); e == nil {
			angle = f
		}
	}
	c1, c2 := model.Analogous(rgb, angle)
	colors := []model.RGB{c1, c2}
	writeColorList(w, r, colors, "analogous")
}

func (h *Handler) triadic(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	c1, c2 := model.Triadic(rgb)
	colors := []model.RGB{c1, c2}
	writeColorList(w, r, colors, "triadic")
}

func (h *Handler) tetradic(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	c1, c2, c3 := model.Tetradic(rgb)
	colors := []model.RGB{c1, c2, c3}
	writeColorList(w, r, colors, "tetradic")
}

func (h *Handler) contrast(w http.ResponseWriter, r *http.Request) {
	c1, err := parseColorParam(r, "c1")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c1 as hex or rgb")
		return
	}
	c2, err := parseColorParam(r, "c2")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c2 as hex or rgb")
		return
	}
	ratio := model.ContrastRatio(c1, c2)
	if wantsJSON(r) {
		writeJSON(w, map[string]interface{}{
			"ratio":    ratio,
			"wcag_aa":  ratio >= 4.5,
			"wcag_aaa": ratio >= 7.0,
			"c1":       model.FullColor(c1),
			"c2":       model.FullColor(c2),
		})
		return
	}
	aa := "fail"
	if ratio >= 4.5 {
		aa = "pass"
	}
	aaa := "fail"
	if ratio >= 7.0 {
		aaa = "pass"
	}
	fmt.Fprintf(w, "ratio=%.2f wcag_aa=%s wcag_aaa=%s\n", ratio, aa, aaa)
}

func (h *Handler) gradient(w http.ResponseWriter, r *http.Request) {
	c1, err := parseColorParam(r, "c1")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c1 as hex or rgb")
		return
	}
	c2, err := parseColorParam(r, "c2")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide c2 as hex or rgb")
		return
	}
	n := 5
	if v := r.URL.Query().Get("steps"); v != "" {
		if i, e := parseInt(v); e == nil && i > 0 {
			n = i
		}
	}
	colors := model.Gradient(c1, c2, n)
	writeColorList(w, r, colors, "gradient")
}

func (h *Handler) shades(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	n := 5
	if v := r.URL.Query().Get("count"); v != "" {
		if i, e := parseInt(v); e == nil && i > 0 {
			n = i
		}
	}
	colors := model.Shades(rgb, n)
	writeColorList(w, r, colors, "shade")
}

func (h *Handler) tints(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	n := 5
	if v := r.URL.Query().Get("count"); v != "" {
		if i, e := parseInt(v); e == nil && i > 0 {
			n = i
		}
	}
	colors := model.Tints(rgb, n)
	writeColorList(w, r, colors, "tint")
}

func (h *Handler) rotate(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	degrees := 30.0
	if v := r.URL.Query().Get("degrees"); v != "" {
		if f, e := parseFloat(v); e == nil {
			degrees = f
		}
	}
	writeColor(w, r, model.RotateHue(rgb, degrees))
}

func (h *Handler) suggestedText(w http.ResponseWriter, r *http.Request) {
	rgb, err := parseColorParam(r, "color")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "provide a color as hex or rgb")
		return
	}
	tc := model.SuggestedTextColor(rgb)
	writeColor(w, r, tc)
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
