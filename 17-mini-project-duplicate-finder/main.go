package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicates")
	flag.Parse()

	fmt.Printf("📂 Scanning directory: %s\n", *dirPtr)
	fmt.Println("---------------------------")

	duplicates := make(map[string][]string)

	err := filepath.Walk(*dirPtr, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process regular files
		if info.Mode().IsRegular() {
			hash, err := calculateHash(path)
			if err != nil {
				fmt.Printf("❌ Error hashing %s: %v\n", path, err)
				return nil
			}
			duplicates[hash] = append(duplicates[hash], path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error walking path: %v\n", err)
		return
	}

	found := false
	for hash, paths := range duplicates {
		if len(paths) > 1 {
			found = true
			fmt.Printf("\n✨ Found %d duplicate files (Hash: %s...):\n", len(paths), hash[:10])
			for _, p := range paths {
				fmt.Printf("  📄 %s\n", p)
			}
		}
	}

	if !found {
		fmt.Println("✅ No duplicate files found.")
	} else {
		fmt.Println("\n---------------------------")
		fmt.Println("Scan complete. Review the files above.")
	}
}

func calculateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
