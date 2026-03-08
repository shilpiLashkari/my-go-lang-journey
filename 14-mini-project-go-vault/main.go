package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const vaultFile = "vault.enc"

type Vault struct {
	Secrets map[string]string `json:"secrets"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-vault <command> [args]")
		fmt.Println("Commands: set, get, list")
		return
	}

	command := os.Args[1]

	switch command {
	case "set":
		if len(os.Args) < 4 {
			fmt.Println("Usage: set <key> <value>")
			return
		}
		setSecret(os.Args[2], os.Args[3])
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: get <key>")
			return
		}
		getSecret(os.Args[2])
	case "list":
		listSecrets()
	default:
		fmt.Println("Unknown command")
	}
}

func getPassword() []byte {
	fmt.Print("Enter Master Password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(password)
	return hash[:]
}

func encrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(encryptedStr string, key []byte) ([]byte, error) {
	data, err := hex.DecodeString(encryptedStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func loadVault(key []byte) (*Vault, error) {
	if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
		return &Vault{Secrets: make(map[string]string)}, nil
	}

	encryptedData, err := os.ReadFile(vaultFile)
	if err != nil {
		return nil, err
	}

	decryptedData, err := decrypt(string(encryptedData), key)
	if err != nil {
		return nil, fmt.Errorf("incorrect password or corrupted file")
	}

	var vault Vault
	err = json.Unmarshal(decryptedData, &vault)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func saveVault(vault *Vault, key []byte) error {
	data, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(data, key)
	if err != nil {
		return err
	}

	return os.WriteFile(vaultFile, []byte(encryptedData), 0600)
}

func setSecret(key, value string) {
	masterKey := getPassword()
	vault, err := loadVault(masterKey)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	vault.Secrets[key] = value
	err = saveVault(vault, masterKey)
	if err != nil {
		fmt.Printf("❌ Error saving vault: %v\n", err)
		return
	}

	fmt.Printf("✅ Secret '%s' saved successfully.\n", key)
}

func getSecret(key string) {
	masterKey := getPassword()
	vault, err := loadVault(masterKey)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	value, ok := vault.Secrets[key]
	if !ok {
		fmt.Printf("❓ Key '%s' not found.\n", key)
		return
	}

	fmt.Printf("🔑 %s: %s\n", key, value)
}

func listSecrets() {
	masterKey := getPassword()
	vault, err := loadVault(masterKey)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n📑 Your Secrets:")
	fmt.Println(strings.Repeat("-", 20))
	for k := range vault.Secrets {
		fmt.Printf("• %s\n", k)
	}
}
