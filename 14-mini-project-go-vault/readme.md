# Day 14: Go-Vault - Secure Secret Store 🔐

## Overview

For Day 14, I've built **Go-Vault**, a CLI tool for securely storing secrets. After working with network scanners and RSS readers, I wanted to dive into how Go handles security and cryptography. This project uses the `crypto` package to encrypt data so that even if my vault file is stolen, the secrets remain safe without my master password.

## Features

- **Authenticated Encryption**: Uses AES-GCM (Advanced Encryption Standard with Galois/Counter Mode) to ensure both confidentiality and integrity.
- **Master Password Protection**: Derives a 32-byte key from my master password using SHA-256.
- **Local Persistence**: Saves all my secrets to an encrypted `vault.enc` file.
- **Simple CLI**: Easy commands to manage keys and values.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 14-mini-project-go-vault
   ```
2. Set a secret (it will prompt for a master password):
   ```bash
   go run main.go set github mypassword123
   ```
3. Retrieve a secret:
   ```bash
   go run main.go get github
   ```
4. List all secrets:
   ```bash
   go run main.go list
   ```

## Learning Reflection

- **Cryptography in Go**: Learned that Go's `crypto` packages are designed to be "vaguely safe by default," but still require careful handling of nonces and keys.
- **Binary Data Handling**: Practiced working with byte slices and hex encoding for storing encrypted payloads.
- **Key Derivation**: Understood why I need to hash or derive keys from passwords rather than using the password directly.

---

_My secrets are now safe behind Go's powerful crypto tools!_
