# Day 34: Go-File-Cryptor - AES-256 File Security 🔒📂

## Overview

For Day 34, I built **Go-File-Cryptor**, a practical utility for file security. It uses industry-standard AES-256-GCM encryption to protect your data. This project demonstrates Go’s excellent standard library support for high-level cryptography, allowing for both data confidentiality and tamper-proof integrity checks.

## Features

- **AES-GCM Encryption**: Uses Galois/Counter Mode for authenticated encryption (it fails if the file was modified after encryption).
- **Zero-Guess Nonces**: Generates cryptographically secure random nonces for every encryption operation.
- **Passphrase Hashing**: Uses SHA-256 to derive a stable 32-byte encryption key from any user-provided string.
- **Dual Mode CLI**: Easy-to-use commands for both `encrypt` and `decrypt` operations.
- **Safe Outputs**: Appends `.enc` to encrypted files instead of overwriting the original by default.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 34-mini-project-file-cryptor
   ```
2. **Encrypt a file**:
   ```bash
   go run main.go encrypt mydata.txt mysecretpassword
   ```
   (Outputs `mydata.txt.enc`)

3. **Decrypt a file**:
   ```bash
   go run main.go decrypt mydata.txt.enc mysecretpassword
   ```

## Learning Reflection

- **Authenticated Encryption (AEAD)**: Discovered that AES-GCM is superior to standard AES-CTR/CBC because it includes an "authentication tag" that detects if a single bit of the file has been tampered with.
- **Cryptographic Randomness**: Used the `crypto/rand` package to ensure nonces are truly random and never reused.
- **Key Derivation**: While PBKDF2/Argon2 are better for passwords, I started with SHA-256 hashing to understand how static-length keys are mapped from variable-length user input.

---
*Locking down the web, one byte at a time. 🐹🔒*
