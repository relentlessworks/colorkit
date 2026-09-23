package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/relentlessworks/colorkit/internal/model"
)

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *mcpError   `json:"error,omitempty"`
}

type mcpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "POST JSON-RPC 2.0 to /mcp")
		return
	}

	var req mcpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMCPError(w, nil, -32700, "parse error")
		return
	}

	switch req.Method {
	case "initialize":
		writeMCPResult(w, req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    "colorkit",
				"version": "0.1.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		})

	case "tools/list":
		writeMCPResult(w, req.ID, map[string]interface{}{
			"tools": mcpTools(),
		})

	case "tools/call":
		var params struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			writeMCPError(w, req.ID, -32602, "invalid params")
			return
		}
		result, err := h.handleMCPTool(params.Name, params.Arguments)
		if err != nil {
			writeMCPError(w, req.ID, -32603, err.Error())
			return
		}
		writeMCPResult(w, req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": result},
			},
		})

	default:
		writeMCPError(w, req.ID, -32601, "method not found")
	}
}

func writeMCPResult(w http.ResponseWriter, id interface{}, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mcpResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func writeMCPError(w http.ResponseWriter, id interface{}, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mcpResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &mcpError{Code: code, Message: msg},
	})
}

type mcpTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func mcpTools() []mcpTool {
	strType := map[string]interface{}{"type": "string"}
	return []mcpTool{
		{Name: "convert", Description: "Convert a color between hex, RGB, and HSL formats", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "mix", Description: "Mix two colors with a ratio", InputSchema: schemaProps(map[string]interface{}{"c1": strType, "c2": strType, "ratio": map[string]interface{}{"type": "string", "description": "0.0 to 1.0, default 0.5"}}, "c1", "c2")},
		{Name: "lighten", Description: "Lighten a color by increasing lightness", InputSchema: schemaProps(map[string]interface{}{"color": strType, "amount": map[string]interface{}{"type": "string", "description": "0-100, default 10"}}, "color")},
		{Name: "darken", Description: "Darken a color by decreasing lightness", InputSchema: schemaProps(map[string]interface{}{"color": strType, "amount": map[string]interface{}{"type": "string", "description": "0-100, default 10"}}, "color")},
		{Name: "saturate", Description: "Increase saturation of a color", InputSchema: schemaProps(map[string]interface{}{"color": strType, "amount": map[string]interface{}{"type": "string", "description": "0-100, default 10"}}, "color")},
		{Name: "desaturate", Description: "Decrease saturation of a color", InputSchema: schemaProps(map[string]interface{}{"color": strType, "amount": map[string]interface{}{"type": "string", "description": "0-100, default 10"}}, "color")},
		{Name: "grayscale", Description: "Convert a color to grayscale", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "invert", Description: "Invert a color", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "complement", Description: "Get the complementary color (180 degrees opposite)", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "analogous", Description: "Get analogous colors (adjacent on color wheel)", InputSchema: schemaProps(map[string]interface{}{"color": strType, "angle": map[string]interface{}{"type": "string", "description": "degrees, default 30"}}, "color")},
		{Name: "triadic", Description: "Get triadic colors (120 degrees apart)", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "tetradic", Description: "Get tetradic colors (90 degrees apart)", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
		{Name: "contrast", Description: "Calculate WCAG contrast ratio between two colors", InputSchema: schemaProps(map[string]interface{}{"c1": strType, "c2": strType}, "c1", "c2")},
		{Name: "gradient", Description: "Generate a gradient between two colors", InputSchema: schemaProps(map[string]interface{}{"c1": strType, "c2": strType, "steps": map[string]interface{}{"type": "string", "description": "number of steps, default 5"}}, "c1", "c2")},
		{Name: "shades", Description: "Generate darker shades of a color", InputSchema: schemaProps(map[string]interface{}{"color": strType, "count": map[string]interface{}{"type": "string", "description": "default 5"}}, "color")},
		{Name: "tints", Description: "Generate lighter tints of a color", InputSchema: schemaProps(map[string]interface{}{"color": strType, "count": map[string]interface{}{"type": "string", "description": "default 5"}}, "color")},
		{Name: "rotate", Description: "Rotate the hue of a color", InputSchema: schemaProps(map[string]interface{}{"color": strType, "degrees": map[string]interface{}{"type": "string", "description": "default 30"}}, "color")},
		{Name: "suggested_text", Description: "Get suggested text color (black or white) for best contrast", InputSchema: schemaProps(map[string]interface{}{"color": strType}, "color")},
	}
}

func schemaProps(props map[string]interface{}, required ...string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

func (h *Handler) handleMCPTool(name string, args map[string]string) (string, error) {
	switch name {
	case "convert":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		return formatColorText(rgb), nil

	case "mix":
		c1, err := parseColorString(args["c1"])
		if err != nil {
			return "", err
		}
		c2, err := parseColorString(args["c2"])
		if err != nil {
			return "", err
		}
		ratio := 0.5
		if v, ok := args["ratio"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				ratio = f
			}
		}
		return formatColorText(model.Mix(c1, c2, ratio)), nil

	case "lighten":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		amount := 10.0
		if v, ok := args["amount"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				amount = f
			}
		}
		return formatColorText(model.Lighten(rgb, amount)), nil

	case "darken":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		amount := 10.0
		if v, ok := args["amount"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				amount = f
			}
		}
		return formatColorText(model.Darken(rgb, amount)), nil

	case "saturate":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		amount := 10.0
		if v, ok := args["amount"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				amount = f
			}
		}
		return formatColorText(model.Saturate(rgb, amount)), nil

	case "desaturate":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		amount := 10.0
		if v, ok := args["amount"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				amount = f
			}
		}
		return formatColorText(model.Desaturate(rgb, amount)), nil

	case "grayscale":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		return formatColorText(model.Grayscale(rgb)), nil

	case "invert":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		return formatColorText(model.Invert(rgb)), nil

	case "complement":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		return formatColorText(model.Complement(rgb)), nil

	case "analogous":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		angle := 30.0
		if v, ok := args["angle"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				angle = f
			}
		}
		c1, c2 := model.Analogous(rgb, angle)
		return formatColorList("analogous", []model.RGB{c1, c2}), nil

	case "triadic":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		c1, c2 := model.Triadic(rgb)
		return formatColorList("triadic", []model.RGB{c1, c2}), nil

	case "tetradic":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		c1, c2, c3 := model.Tetradic(rgb)
		return formatColorList("tetradic", []model.RGB{c1, c2, c3}), nil

	case "contrast":
		c1, err := parseColorString(args["c1"])
		if err != nil {
			return "", err
		}
		c2, err := parseColorString(args["c2"])
		if err != nil {
			return "", err
		}
		ratio := model.ContrastRatio(c1, c2)
		aa := "fail"
		if ratio >= 4.5 {
			aa = "pass"
		}
		aaa := "fail"
		if ratio >= 7.0 {
			aaa = "pass"
		}
		return fmt.Sprintf("ratio=%.2f wcag_aa=%s wcag_aaa=%s", ratio, aa, aaa), nil

	case "gradient":
		c1, err := parseColorString(args["c1"])
		if err != nil {
			return "", err
		}
		c2, err := parseColorString(args["c2"])
		if err != nil {
			return "", err
		}
		n := 5
		if v, ok := args["steps"]; ok {
			if i, err := strconv.Atoi(v); err == nil && i > 0 {
				n = i
			}
		}
		return formatColorList("gradient", model.Gradient(c1, c2, n)), nil

	case "shades":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		n := 5
		if v, ok := args["count"]; ok {
			if i, err := strconv.Atoi(v); err == nil && i > 0 {
				n = i
			}
		}
		return formatColorList("shade", model.Shades(rgb, n)), nil

	case "tints":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		n := 5
		if v, ok := args["count"]; ok {
			if i, err := strconv.Atoi(v); err == nil && i > 0 {
				n = i
			}
		}
		return formatColorList("tint", model.Tints(rgb, n)), nil

	case "rotate":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		degrees := 30.0
		if v, ok := args["degrees"]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				degrees = f
			}
		}
		return formatColorText(model.RotateHue(rgb, degrees)), nil

	case "suggested_text":
		rgb, err := parseColorString(args["color"])
		if err != nil {
			return "", err
		}
		return formatColorText(model.SuggestedTextColor(rgb)), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func formatColorText(rgb model.RGB) string {
	hsl := model.RGBToHSL(rgb)
	return fmt.Sprintf("hex=%s rgb=%s hsl=%s", model.RGBToHex(rgb), model.FormatRGB(rgb), model.FormatHSL(hsl))
}

func formatColorList(prefix string, colors []model.RGB) string {
	var sb strings.Builder
	for i, c := range colors {
		hsl := model.RGBToHSL(c)
		sb.WriteString(fmt.Sprintf("%s[%d] hex=%s rgb=%s hsl=%s\n", prefix, i, model.RGBToHex(c), model.FormatRGB(c), model.FormatHSL(hsl)))
	}
	return sb.String()
}
