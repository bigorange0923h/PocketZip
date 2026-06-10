package security

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("test password 123456")

	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

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