package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("❌ Usage:")
		fmt.Println("  Encrypt: go run main.go encrypt <filename> <password>")
		fmt.Println("  Decrypt: go run main.go decrypt <filename> <password>")
		return
	}

	mode := os.Args[1]
	filename := os.Args[2]
	password := os.Args[3]

	// Derive a 32-byte key from the password
	key := sha256.Sum256([]byte(password))

	switch mode {
	case "encrypt":
		err := encryptFile(filename, key[:])
		if err != nil {
			fmt.Printf("❌ Encryption failed: %v\n", err)
		} else {
			fmt.Printf("✅ File encrypted: %s.enc\n", filename)
		}
	case "decrypt":
		err := decryptFile(filename, key[:])
		if err != nil {
			fmt.Printf("❌ Decryption failed: %v (Is the password correct?)\n", err)
		} else {
			fmt.Printf("✅ File decrypted successfully!\n")
		}
	default:
		fmt.Println("❌ Invalid mode. Use 'encrypt' or 'decrypt'.")
	}
}

func encryptFile(filename string, key []byte) error {
	plaintext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Create a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Encrypt: Nonce + Ciphertext (GCM appends the tag at the end)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return os.WriteFile(filename+".enc", ciphertext, 0644)
}

func decryptFile(filename string, key []byte) error {
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return err
	}

	// Save decrypted content (removes .enc if present)
	outPath := filename
	if len(filename) > 4 && filename[len(filename)-4:] == ".enc" {
		outPath = filename[:len(filename)-4]
	} else {
		outPath = "decrypted_" + filename
	}

	return os.WriteFile(outPath, plaintext, 0644)
}
