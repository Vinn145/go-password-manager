# 🔐 go-password-manager

Simple CLI Password Manager built with Golang.  
Passwords are encrypted with **AES-256** and stored locally in `vault.json`.

---

## ✨ Features
- ➕ Add new password with alias  
- 🔓 Retrieve password by alias  
- 📋 List all saved aliases  
- 🗑️ Delete password entry  
- 🔒 AES-256 encryption  

---

## ⚙️ Installation & Run

Clone repo:
```bash
git clone https://github.com/Vinn145/go-password-manager.git
cd go-password-manager

# Add password
go run main.go add github mysecret123

# List saved aliases
go run main.go list

# Get password
go run main.go get github

# Delete password
go run main.go delete github
