package colors

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ColorMatch struct {
	Original string
	Variable string
	Value    string
	Line     int
}

var (
	hexColorRegex  = regexp.MustCompile(`(?i)#([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})(?:[)\s;,}]|$)`)
	rgbColorRegex  = regexp.MustCompile(`rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)`)
	rgbaColorRegex = regexp.MustCompile(`rgba\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(0|1|0?\.\d+)\s*\)`)
)

func ScanFile(filepath string) ([]ColorMatch, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var matches []ColorMatch
	colorMap := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		hexMatches := hexColorRegex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range hexMatches {
			if len(match) >= 4 {
				start := match[0] // Start of entire match
				end := match[3]   // End of the hex color part
				if start >= 0 && end <= len(line) {
					colorPart := line[start:end]
					colorPart = strings.TrimSpace(colorPart)
					if !colorMap[colorPart] {
						colorMap[colorPart] = true
						varName := GenerateVariableName(colorPart)
						matches = append(matches, ColorMatch{
							Original: colorPart,
							Variable: varName,
							Value:    colorPart,
							Line:     lineNum,
						})
					}
				}
			}
		}

		rgbMatches := rgbColorRegex.FindAllString(line, -1)
		for _, match := range rgbMatches {
			if !colorMap[match] {
				colorMap[match] = true
				varName := GenerateVariableName(match)
				matches = append(matches, ColorMatch{
					Original: match,
					Variable: varName,
					Value:    match,
					Line:     lineNum,
				})
			}
		}

		rgbaMatches := rgbaColorRegex.FindAllString(line, -1)
		for _, match := range rgbaMatches {
			if !colorMap[match] {
				colorMap[match] = true
				varName := GenerateVariableName(match)
				matches = append(matches, ColorMatch{
					Original: match,
					Variable: varName,
					Value:    match,
					Line:     lineNum,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return matches, nil
}

func GenerateVariableName(color string) string {
	name := strings.ToLower(color)
	name = strings.ReplaceAll(name, "#", "")
	name = strings.ReplaceAll(name, "(", "-")
	name = strings.ReplaceAll(name, ")", "-")
	name = strings.ReplaceAll(name, ", ", "-")
	name = strings.ReplaceAll(name, ",", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")

	if !strings.HasPrefix(name, "--color-") {
		name = "--color-" + name
	}

	if strings.Contains(name, "-rgb") {
		name = strings.TrimSuffix(name, "-") + "-"
	}

	return name
}

func ConvertColor(color, format string) (string, error) {
	r, g, b, a := ParseToRGBA(color)

	switch format {
	case "hex":
		if a < 1 {
			alphaHex := uint8(math.Round(a * 255))
			return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, alphaHex), nil
		}
		return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
	case "rgb":
		return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b), nil
	case "rgba":
		return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, a), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func ParseToRGBA(color string) (r, g, b uint8, a float64) {
	a = 1.0

	if strings.HasPrefix(color, "rgba(") {
		parts := rgbaColorRegex.FindStringSubmatch(color)
		if len(parts) == 5 {
			r = uint8(parseDecimal(parts[1]))
			g = uint8(parseDecimal(parts[2]))
			b = uint8(parseDecimal(parts[3]))
			a, _ = strconv.ParseFloat(parts[4], 64)
			a = float64(int(a*100)) / 100
			return
		}
	}

	if strings.HasPrefix(color, "rgb(") {
		parts := rgbColorRegex.FindStringSubmatch(color)
		if len(parts) == 4 {
			r = uint8(parseDecimal(parts[1]))
			g = uint8(parseDecimal(parts[2]))
			b = uint8(parseDecimal(parts[3]))
			return
		}
	}

	if strings.HasPrefix(color, "#") {
		hex := strings.TrimPrefix(color, "#")
		switch len(hex) {
		case 3: // #RGB
			r = parseHex(string(hex[0]) + string(hex[0]))
			g = parseHex(string(hex[1]) + string(hex[1]))
			b = parseHex(string(hex[2]) + string(hex[2]))
		case 4: // #RGBA
			r = parseHex(string(hex[0]) + string(hex[0]))
			g = parseHex(string(hex[1]) + string(hex[1]))
			b = parseHex(string(hex[2]) + string(hex[2]))
			a = float64(parseHex(string(hex[3])+string(hex[3]))) / 255
			// Round alpha to 2 decimal places
			a = float64(int(a*100)) / 100
		case 6: // #RRGGBB
			r = parseHex(hex[0:2])
			g = parseHex(hex[2:4])
			b = parseHex(hex[4:6])
		case 8: // #RRGGBBAA
			r = parseHex(hex[0:2])
			g = parseHex(hex[2:4])
			b = parseHex(hex[4:6])
			a = float64(parseHex(hex[6:8])) / 255
			// Round alpha to 2 decimal places
			a = float64(int(a*100)) / 100
		}
		return
	}

	return
}

func parseHex(hex string) uint8 {
	val, _ := strconv.ParseUint(hex, 16, 8)
	return uint8(val)
}

func parseDecimal(dec string) int {
	val, _ := strconv.Atoi(dec)
	if val > 255 {
		val = 255
	} else if val < 0 {
		val = 0
	}
	return val
}
