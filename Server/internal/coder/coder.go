package coder

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var (
	key   = []byte{248, 48, 135, 120, 118, 157, 242, 205, 202, 4, 151, 69, 142, 14, 146, 124, 159, 70, 24, 162, 31, 209, 250, 178, 15, 153, 83, 13, 28, 21, 217, 192}
	nonce = []byte{201, 164, 211, 227, 211, 34, 224, 13, 99, 11, 232, 220}
)

func Encrypt(inp []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	dst := aesgcm.Seal(nil, nonce, inp, nil) // зашифровываем
	return dst, nil
}

func Derypt(inp []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	dst, err := aesgcm.Open(nil, nonce, inp, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	return dst, nil
}

// func main() {
// 	src := []byte("Ключ от сердца") // данные, которые хотим зашифровать
// 	fmt.Printf("original: %s\n", src)

// 	// NewCipher создает и возвращает новый cipher.Block.
// 	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
// 	// для выбора AES-128, AES-192 или AES-256.
// 	aesblock, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	// NewGCM возвращает заданный 128-битный блочный шифр
// 	aesgcm, err := cipher.NewGCM(aesblock)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
// 	fmt.Printf("encrypted: %x\n", dst)

// 	src2, err := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("decrypted: %s\n", src2)
// }
