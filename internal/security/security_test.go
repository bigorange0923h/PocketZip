package security

import "testing"

func TestMaskPasswordArg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"with password", "7z.exe x test.zip -ooutput -p123456 -y", "7z.exe x test.zip -ooutput -p****** -y"},
		{"no password", "7z.exe x test.zip -ooutput -y", "7z.exe x test.zip -ooutput -y"},
		{"empty password", "7z.exe x test.zip -ooutput -p -y", "7z.exe x test.zip -ooutput -p -y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskPasswordArg(tt.input)
			if got != tt.want {
				t.Errorf("MaskPasswordArg(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}