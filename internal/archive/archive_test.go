package archive

import (
	"testing"
)

func TestBuildArgs_NoPassword(t *testing.T) {
	req := ExtractRequest{
		SevenZipPath: "7z.exe",
		ArchivePath:  "test.zip",
		OutputDir:    "output",
	}
	args := buildArgs(req)
	expected := []string{"x", "test.zip", "-ooutput", "-y"}
	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(args))
	}
	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestBuildArgs_WithPassword(t *testing.T) {
	req := ExtractRequest{
		SevenZipPath: "7z.exe",
		ArchivePath:  "test.zip",
		OutputDir:    "output",
		Password:     "123456",
	}
	args := buildArgs(req)
	expected := []string{"x", "test.zip", "-ooutput", "-p123456", "-y"}
	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(args))
	}
	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestIsPasswordError(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{"english", "Wrong password", true},
		{"english lower case", "error: password is incorrect", true},
		{"7z encrypted data", "Data Error in encrypted file. Wrong password?", true},
		{"chinese", "密码错误", true},
		{"encrypted", "Cannot open encrypted archive", true},
		{"normal error", "File not found", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordError(tt.output); got != tt.want {
				t.Errorf("IsPasswordError(%q) = %v, want %v", tt.output, got, tt.want)
			}
		})
	}
}

func TestExtractResultErrorIncludesOutput(t *testing.T) {
	result := ExtractResult{
		Success: false,
		ExitErr: errForTest("exit status 2"),
		Output:  "ERROR: Wrong password",
	}

	err := result.Error()
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "exit status 2: ERROR: Wrong password" {
		t.Fatalf("unexpected error: %q", got)
	}
}

type errForTest string

func (e errForTest) Error() string {
	return string(e)
}
