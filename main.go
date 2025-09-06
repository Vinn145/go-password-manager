package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ==========================
// ENCRYPT / DECRYPT
// ==========================
var key = []byte("12345678901234567890123456789012") // ✅ 32 byte

func encrypt(plain string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plain))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// ==========================
// STRUCT & VAULT FUNCTIONS
// ==========================
type Entry struct {
	Alias     string `json:"alias"`
	Encrypted string `json:"encrypted"`
}

func loadVault() ([]Entry, error) {
	var entries []Entry
	data, err := os.ReadFile("vault.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil // kalau file belum ada, return kosong
		}
		return nil, err
	}
	json.Unmarshal(data, &entries)
	return entries, nil
}

func saveVault(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("vault.json", data, 0644)
}

// ==========================
// CRUD FUNCTIONS
// ==========================
func addPassword(alias, plain string) error {
	encrypted, err := encrypt(plain)
	if err != nil {
		return err
	}

	entries, _ := loadVault()
	entries = append(entries, Entry{Alias: alias, Encrypted: encrypted})
	return saveVault(entries)
}

func getPassword(alias string) (string, error) {
	entries, _ := loadVault()
	for _, e := range entries {
		if e.Alias == alias {
			return decrypt(e.Encrypted)
		}
	}
	return "", fmt.Errorf("alias %s tidak ditemukan", alias)
}

func listPasswords() {
	entries, _ := loadVault()
	if len(entries) == 0 {
		fmt.Println("⚠️ Vault masih kosong")
		return
	}
	for _, e := range entries {
		fmt.Println("🔑", e.Alias)
	}
}

func deletePassword(alias string) error {
	entries, _ := loadVault()
	var newEntries []Entry
	found := false
	for _, e := range entries {
		if e.Alias != alias {
			newEntries = append(newEntries, e)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("alias %s tidak ditemukan", alias)
	}
	return saveVault(newEntries)
}

// ==========================
// MAIN CLI
// ==========================
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [add|get|list|delete] ...")
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go add <alias> <password>")
			return
		}
		err := addPassword(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println("❌ Error saving password:", err)
			return
		}
		fmt.Println("✅ Password added for", os.Args[2])

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go get <alias>")
			return
		}
		pass, err := getPassword(os.Args[2])
		if err != nil {
			fmt.Println("❌", err)
			return
		}
		fmt.Println("🔓", pass)

	case "list":
		listPasswords()

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go delete <alias>")
			return
		}
		err := deletePassword(os.Args[2])
		if err != nil {
			fmt.Println("❌", err)
			return
		}
		fmt.Println("🗑️ Deleted", os.Args[2])

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
