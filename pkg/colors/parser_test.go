package colors

import (
	"os"
	"strings"
	"testing"
)

func TestConvertColor(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		format      string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:   "hex to rgb",
			input:  "#ff0000",
			format: "rgb",
			want:   "rgb(255, 0, 0)",
		},
		{
			name:   "hex to rgba",
			input:  "#ff0000",
			format: "rgba",
			want:   "rgba(255, 0, 0, 1.00)",
		},
		{
			name:   "rgb to hex",
			input:  "rgb(255, 0, 0)",
			format: "hex",
			want:   "#ff0000",
		},
		{
			name:   "rgba to hex with alpha",
			input:  "rgba(255, 0, 0, 0.5)",
			format: "hex",
			want:   "#ff000080",
		},
		{
			name:   "short hex to rgb",
			input:  "#f00",
			format: "rgb",
			want:   "rgb(255, 0, 0)",
		},
		{
			name:   "rgba to rgb",
			input:  "rgba(255, 0, 0, 0.5)",
			format: "rgb",
			want:   "rgb(255, 0, 0)",
		},
		{
			name:        "invalid format",
			input:       "#ff0000",
			format:      "invalid",
			wantErr:     true,
			errContains: "unsupported format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertColor(tt.input, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ConvertColor() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if got != tt.want {
				t.Errorf("ConvertColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateVariableName(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		{
			name:  "hex color",
			color: "#ff0000",
			want:  "--color-ff0000",
		},
		{
			name:  "rgb color",
			color: "rgb(255, 0, 0)",
			want:  "--color-rgb-255-0-0-",
		},
		{
			name:  "rgba color",
			color: "rgba(255, 0, 0, 0.5)",
			want:  "--color-rgba-255-0-0-0-5-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateVariableName(tt.color); got != tt.want {
				t.Errorf("GenerateVariableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseToRGBA(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		wantR    uint8
		wantG    uint8
		wantB    uint8
		wantA    float64
		wantZero bool
	}{
		{
			name:  "full hex",
			color: "#ff0000",
			wantR: 255,
			wantG: 0,
			wantB: 0,
			wantA: 1.0,
		},
		{
			name:  "short hex",
			color: "#f00",
			wantR: 255,
			wantG: 0,
			wantB: 0,
			wantA: 1.0,
		},
		{
			name:  "hex with alpha",
			color: "#ff000080",
			wantR: 255,
			wantG: 0,
			wantB: 0,
			wantA: 0.5,
		},
		{
			name:  "rgb",
			color: "rgb(255, 0, 0)",
			wantR: 255,
			wantG: 0,
			wantB: 0,
			wantA: 1.0,
		},
		{
			name:  "rgba",
			color: "rgba(255, 0, 0, 0.5)",
			wantR: 255,
			wantG: 0,
			wantB: 0,
			wantA: 0.5,
		},
		{
			name:     "invalid color",
			color:    "invalid",
			wantZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotG, gotB, gotA := ParseToRGBA(tt.color)
			if tt.wantZero {
				if gotR != 0 || gotG != 0 || gotB != 0 || gotA != 1.0 {
					t.Errorf("ParseToRGBA() = (%v, %v, %v, %v), want all zero values", gotR, gotG, gotB, gotA)
				}
				return
			}
			if gotR != tt.wantR || gotG != tt.wantG || gotB != tt.wantB || gotA != tt.wantA {
				t.Errorf("ParseToRGBA() = (%v, %v, %v, %v), want (%v, %v, %v, %v)",
					gotR, gotG, gotB, gotA, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
			}
		})
	}
}

func TestScanFile(t *testing.T) {
	content := `
.test {
	color: #ff0000;
	background: rgb(0, 255, 0);
	border-color: rgba(0, 0, 255, 0.5);
}
`
	tmpfile, err := os.CreateTemp("", "test*.css")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	matches, err := ScanFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ScanFile() error = %v", err)
	}

	expected := []struct {
		original string
		value    string
	}{
		{"#ff0000", "#ff0000"},
		{"rgb(0, 255, 0)", "rgb(0, 255, 0)"},
		{"rgba(0, 0, 255, 0.5)", "rgba(0, 0, 255, 0.5)"},
	}

	if len(matches) != len(expected) {
		t.Errorf("ScanFile() returned %d matches, want %d", len(matches), len(expected))
	}

	for i, exp := range expected {
		if i >= len(matches) {
			break
		}
		if matches[i].Original != exp.original || matches[i].Value != exp.value {
			t.Errorf("Match[%d] = {Original: %q, Value: %q}, want {Original: %q, Value: %q}",
				i, matches[i].Original, matches[i].Value, exp.original, exp.value)
		}
	}
}
