//go:build !windows

package security

func Encrypt(plaintext []byte) ([]byte, error) {
	// 非 Windows 平台，直接返回明文（仅用于开发测试）
	return plaintext, nil
}

func Decrypt(ciphertext []byte) ([]byte, error) {
	// 非 Windows 平台，直接返回密文（仅用于开发测试）
	return ciphertext, nil
}
