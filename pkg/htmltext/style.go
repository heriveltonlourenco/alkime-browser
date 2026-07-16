package htmltext

import (
	"image/color"
	"strconv"
	"strings"
)

// parseInlineStyle parses a minimal subset of CSS out of a
// style="..." attribute value: semicolon-separated "property: value"
// declarations. Only color and background-color are understood; any
// other property, or a value this package can't parse, is silently
// ignored — this is a deliberately small, best-effort subset, not a
// spec-compliant CSS parser. There is no selector matching, no
// cascade: this only ever reads an element's own inline style.
func parseInlineStyle(style string) (fg, bg *color.RGBA) {
	for _, decl := range strings.Split(style, ";") {
		parts := strings.SplitN(decl, ":", 2)
		if len(parts) != 2 {
			continue
		}
		prop := strings.ToLower(strings.TrimSpace(parts[0]))
		c, ok := parseCSSColor(parts[1])
		if !ok {
			continue
		}
		switch prop {
		case "color":
			fg = &c
		case "background-color":
			bg = &c
		}
	}
	return fg, bg
}

// namedColors covers the handful of CSS color keywords a real page
// is likely to use; anything beyond this list falls back to hex.
var namedColors = map[string]color.RGBA{
	"black":  {R: 0, G: 0, B: 0, A: 255},
	"white":  {R: 255, G: 255, B: 255, A: 255},
	"red":    {R: 255, G: 0, B: 0, A: 255},
	"green":  {R: 0, G: 128, B: 0, A: 255},
	"blue":   {R: 0, G: 0, B: 255, A: 255},
	"yellow": {R: 255, G: 255, B: 0, A: 255},
	"orange": {R: 255, G: 165, B: 0, A: 255},
	"gray":   {R: 128, G: 128, B: 128, A: 255},
	"grey":   {R: 128, G: 128, B: 128, A: 255},
	"purple": {R: 128, G: 0, B: 128, A: 255},
	"pink":   {R: 255, G: 192, B: 203, A: 255},
	"cyan":   {R: 0, G: 255, B: 255, A: 255},
}

// parseCSSColor understands #rgb, #rrggbb, and the namedColors above.
// Everything else (rgb(), hsl(), currentColor, ...) reports false.
func parseCSSColor(v string) (color.RGBA, bool) {
	v = strings.ToLower(strings.TrimSpace(v))

	if strings.HasPrefix(v, "#") {
		hex := v[1:]
		if len(hex) == 3 {
			hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
		}
		if len(hex) != 6 {
			return color.RGBA{}, false
		}
		r, err1 := strconv.ParseUint(hex[0:2], 16, 8)
		g, err2 := strconv.ParseUint(hex[2:4], 16, 8)
		b, err3 := strconv.ParseUint(hex[4:6], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return color.RGBA{}, false
		}
		return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, true
	}

	if c, ok := namedColors[v]; ok {
		return c, true
	}
	return color.RGBA{}, false
}
