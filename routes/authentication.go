package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"to-do-backend/database"
	"to-do-backend/models"

	"github.com/gofiber/fiber/v2"
)

func checkEmailExists(email string) bool {
	user := models.User{}
	result := database.Database.Db.Where("email_id = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func Signup(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	user.Password = hashPassword(user.Password, getSalt())
	if checkEmailExists(user.EmailID) {
		return c.Status(400).JSON("Email already exists")
	}
	return CreateUser(c)
}

func getSalt() []byte {
	return []byte("3q1ftcbddAxQ3Cav")
}

func hashPassword(password string, salt []byte) string {
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

func doPasswordsMatch(hashedPassword string, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)
	return hashedPassword == currPasswordHash
}

// Fix hardcoded password
func Login(c *fiber.Ctx) error {
	randomSalt := getSalt()
	hashedPassword := hashPassword("abcd", randomSalt)
	type UserCred struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	}
	var userCred UserCred

	if err := c.BodyParser(&userCred); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if doPasswordsMatch(hashedPassword, userCred.Password, randomSalt) {
		return c.SendString("Successful")
	} else {
		return c.SendString("Unsuccessful")
	}
}
