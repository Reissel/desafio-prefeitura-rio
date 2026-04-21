package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

var (
	encryptionKey = []byte("passphrasewhichneedstobe32bytes!")
	indexSalt     = []byte("your-secret-salt-for-blind-index")
)

func GenerateBlindIndex(cpf string) string {
	sig := hmac.New(sha256.New, indexSalt)
	sig.Write([]byte(cpf))
	return hex.EncodeToString(sig.Sum(nil))
}

func Encrypt(plainText string) ([]byte, error) {
	block, _ := aes.NewCipher(encryptionKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, []byte(plainText), nil), nil
}

func Decrypt(cipherText []byte) (string, error) {
	block, _ := aes.NewCipher(encryptionKey)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}
