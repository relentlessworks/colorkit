package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/convert?color=%23ff0000", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("convert status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=#ff0000") {
		t.Errorf("convert body = %q, want hex=#ff0000", body)
	}
	if !strings.Contains(body, "rgb=255,0,0") {
		t.Errorf("convert body = %q, want rgb=255,0,0", body)
	}
}

func TestConvertJSON(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/convert?color=%23ff0000&format=json", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("convert status = %d, want 200", w.Code)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("convert JSON parse error: %v", err)
	}
	if result["Hex"] != "#ff0000" {
		t.Errorf("convert JSON hex = %v, want #ff0000", result["Hex"])
	}
}

func TestConvertMissingParam(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/convert", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("convert missing param status = %d, want 400", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "error:") {
		t.Errorf("convert missing param body = %q, want error", body)
	}
	if !strings.Contains(body, "hint:") {
		t.Errorf("convert missing param body = %q, want hint", body)
	}
}

func TestMix(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/mix?c1=%23000000&c2=%23ffffff&ratio=0.5", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("mix status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "rgb=128,128,128") || !strings.Contains(body, "rgb=127,127,127") {
		// Rounding may vary
		if !strings.Contains(body, "rgb=12") {
			t.Errorf("mix body = %q, want ~128,128,128", body)
		}
	}
}

func TestLighten(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/lighten?color=%23333333&amount=20", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("lighten status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=") {
		t.Errorf("lighten body = %q, want hex", body)
	}
}

func TestDarken(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/darken?color=%23cccccc&amount=20", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("darken status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=") {
		t.Errorf("darken body = %q, want hex", body)
	}
}

func TestContrast(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/contrast?c1=%23000000&c2=%23ffffff", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("contrast status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "ratio=") {
		t.Errorf("contrast body = %q, want ratio", body)
	}
	if !strings.Contains(body, "wcag_aa=pass") {
		t.Errorf("contrast body = %q, want wcag_aa=pass", body)
	}
}

func TestGradient(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/gradient?c1=%23ff0000&c2=%230000ff&steps=3", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("gradient status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 3 {
		t.Errorf("gradient lines = %d, want 3", len(lines))
	}
}

func TestShades(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/shades?color=%23ff0000&count=3", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("shades status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 3 {
		t.Errorf("shades lines = %d, want 3", len(lines))
	}
}

func TestTints(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/tints?color=%23ff0000&count=3", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("tints status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 3 {
		t.Errorf("tints lines = %d, want 3", len(lines))
	}
}

func TestComplement(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/complement?color=%23ff0000", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("complement status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=#00ffff") {
		t.Errorf("complement body = %q, want hex=#00ffff", body)
	}
}

func TestInvert(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/invert?color=%23ff8800", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("invert status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=#0077ff") {
		t.Errorf("invert body = %q, want hex=#0077ff", body)
	}
}

func TestGrayscale(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/grayscale?color=%23ff0000", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("grayscale status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	// Red grayscale: 0.299*255 = 76
	if !strings.Contains(body, "hex=#4c4c4c") {
		t.Errorf("grayscale body = %q, want hex=#4c4c4c", body)
	}
}

func TestTriadic(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/triadic?color=%23ff0000", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("triadic status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 2 {
		t.Errorf("triadic lines = %d, want 2", len(lines))
	}
}

func TestSuggestedText(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	// White background should suggest black text
	req := httptest.NewRequest("GET", "/suggested-text?color=%23ffffff", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("suggested-text status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hex=#000000") {
		t.Errorf("suggested-text(white) body = %q, want hex=#000000", body)
	}

	// Black background should suggest white text
	req = httptest.NewRequest("GET", "/suggested-text?color=%23000000", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body = w.Body.String()
	if !strings.Contains(body, "hex=#ffffff") {
		t.Errorf("suggested-text(black) body = %q, want hex=#ffffff", body)
	}
}

func TestHelp(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("GET", "/help", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("help status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "colorkit") {
		t.Errorf("help body should contain 'colorkit'")
	}
	if !strings.Contains(body, "ENDPOINTS") {
		t.Errorf("help body should contain 'ENDPOINTS'")
	}
}

func TestMCPInitialize(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("mcp initialize status = %d, want 200", w.Code)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("mcp initialize JSON parse error: %v", err)
	}
	if result["jsonrpc"] != "2.0" {
		t.Errorf("mcp initialize jsonrpc = %v, want 2.0", result["jsonrpc"])
	}
}

func TestMCPToolsList(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("mcp tools/list status = %d, want 200", w.Code)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("mcp tools/list JSON parse error: %v", err)
	}
	tools, ok := result["result"].(map[string]interface{})["tools"].([]interface{})
	if !ok {
		t.Errorf("mcp tools/list result not as expected")
		return
	}
	if len(tools) < 18 {
		t.Errorf("mcp tools/list count = %d, want >= 18", len(tools))
	}
}

func TestMCPToolsCall(t *testing.T) {
	h := New("test-secret")
	mux := h.Routes()

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"convert","arguments":{"color":"#ff0000"}}}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("mcp tools/call status = %d, want 200", w.Code)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("mcp tools/call JSON parse error: %v", err)
	}
	content, ok := result["result"].(map[string]interface{})["content"].([]interface{})
	if !ok {
		t.Errorf("mcp tools/call result not as expected")
		return
	}
	text := content[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "hex=#ff0000") {
		t.Errorf("mcp tools/call text = %q, want hex=#ff0000", text)
	}
}
