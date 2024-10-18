package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Şifreleme fonksiyonu
func encrypt(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Şifre çözme fonksiyonu
func decrypt(encryptedData string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("şifrelenmiş veri çok kısa")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func main() {
	key := []byte("this32bytekey234567890*jgsaw2453")
	text := "892627498883-dilm5tggiet8c895qr8vfo7dn83bgpj1.apps.googleusercontent.com"

	// Veriyi şifreleme
	encrypted, err := encrypt([]byte(text), key)
	if err != nil {
		fmt.Println("Şifreleme hatası:", err)
		return
	}
	fmt.Println("Şifrelenmiş veri:", encrypted)

	// Şifreli veriyi çözme
	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		fmt.Println("Şifre çözme hatası:", err)
		return
	}
	fmt.Println("Çözülen veri:", decrypted)
}
