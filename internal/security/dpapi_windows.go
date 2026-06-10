//go:build windows

package security

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Encrypt 使用 Windows DPAPI 加密数据
func Encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return plaintext, nil
	}

	// 转换为 UTF-16 字符串
	utf16, err := windows.UTF16FromString(string(plaintext))
	if err != nil {
		return nil, err
	}

	// 将 UTF-16 转换为字节
	dataIn := make([]byte, len(utf16)*2)
	for i, v := range utf16 {
		dataIn[i*2] = byte(v)
		dataIn[i*2+1] = byte(v >> 8)
	}

	// 创建 DATA_BLOB
	var dataInBlob windows.DataBlob
	dataInBlob.Size = uint32(len(dataIn))
	if len(dataIn) > 0 {
		dataInBlob.Data = &dataIn[0]
	}

	var dataOutBlob windows.DataBlob

	// 调用 CryptProtectData
	err = windows.CryptProtectData(
		&dataInBlob,
		nil, // 描述
		nil, // 可选熵
		0,   // 保留
		nil, // 提示
		windows.CRYPTPROTECT_UI_FORBIDDEN,
		&dataOutBlob,
	)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(dataOutBlob.Data)))

	// 复制输出数据
	encrypted := make([]byte, dataOutBlob.Size)
	copy(encrypted, unsafe.Slice(dataOutBlob.Data, dataOutBlob.Size))

	return encrypted, nil
}

// Decrypt 使用 Windows DPAPI 解密数据
func Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return ciphertext, nil
	}

	// 创建 DATA_BLOB
	var dataInBlob windows.DataBlob
	dataInBlob.Size = uint32(len(ciphertext))
	if len(ciphertext) > 0 {
		dataInBlob.Data = &ciphertext[0]
	}

	var dataOutBlob windows.DataBlob

	// 调用 CryptUnprotectData
	err := windows.CryptUnprotectData(
		&dataInBlob,
		nil, // 描述
		nil, // 可选熵
		0,   // 保留
		nil, // 提示
		windows.CRYPTPROTECT_UI_FORBIDDEN,
		&dataOutBlob,
	)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(dataOutBlob.Data)))

	// 复制输出数据（UTF-16 格式）
	decryptedUTF16 := make([]byte, dataOutBlob.Size)
	copy(decryptedUTF16, unsafe.Slice(dataOutBlob.Data, dataOutBlob.Size))

	// 转换 UTF-16 字节为 Go 字符串
	if len(decryptedUTF16)%2 != 0 {
		return nil, errors.New("invalid UTF-16 data")
	}

	utf16Slice := make([]uint16, len(decryptedUTF16)/2)
	for i := range utf16Slice {
		utf16Slice[i] = uint16(decryptedUTF16[i*2]) | uint16(decryptedUTF16[i*2+1])<<8
	}

	// 去除末尾的 null 字符
	for len(utf16Slice) > 0 && utf16Slice[len(utf16Slice)-1] == 0 {
		utf16Slice = utf16Slice[:len(utf16Slice)-1]
	}

	return []byte(windows.UTF16ToString(utf16Slice)), nil
}
