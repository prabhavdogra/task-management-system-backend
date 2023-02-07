package routes

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSalt() []byte {
	return []byte("3q1ftcbddAxQ3Cav")
}

func HashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha256.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func DoPasswordsMatch(hashedPassword string, currPassword string, salt []byte) bool {
	var currPasswordHash = HashPassword(currPassword, salt)
	return hashedPassword == currPasswordHash
}
