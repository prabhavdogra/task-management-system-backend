package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"to-do-backend/database"
	"to-do-backend/models"

	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v2"
)

func CheckEmailExists(email string) bool {
	user := models.User{}
	result := database.Database.Db.Where("email_id = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func CheckAgainstLoggedOutTokens(token string) bool {
	return false
}

var sampleSecretKey = []byte("GoLinuxCloudKey")

func CreateJWT(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = user.EmailID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Signup(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if CheckEmailExists(user.EmailID) {
		return c.Status(400).JSON("Email already exists")
	}
	user.Password = HashPassword(user.Password, GetSalt())
	database.Database.Db.Save(&user)
	token, err := CreateJWT(user)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).SendString(token)
}

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

// func fn(c *fiber.Ctx, x string) error {
// 	return c.JSON(fiber.Map{"k": x})
// }

func AuthenticateJWTToken(c *fiber.Ctx) error {
	// CheckAgainstLoggedOutTokens()
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	token, err := jwt.Parse(authData.Authorization[7:], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return sampleSecretKey, nil
	})

	if token == nil {
		return c.Status(400).SendString(err.Error())
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if !token.Valid {
		return c.Status(400).SendString("Invalid token")
	}

	if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
		return c.Status(400).SendString("Session expired!")
	}
	return c.Status(400).SendString("Valid Token!")
	// }
	// if token == nil {
	// 	return false
	// }
	// return true
}

func Login(c *fiber.Ctx) error {
	randomSalt := GetSalt()
	type UserCred struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	}
	var userCred UserCred
	if err := c.BodyParser(&userCred); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	user := models.User{}
	database.Database.Db.Where("email_id = ?", userCred.EmailID).First(&user)
	hashedPassword := user.Password
	if DoPasswordsMatch(hashedPassword, userCred.Password, randomSalt) {
		token, err := CreateJWT(user)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.Status(400).JSON(fiber.Map{
			"token": token,
		})
	} else {
		return c.SendString("Login Unsuccessful")
	}
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzU3MDExMzcsInVzZXJuYW1lIjoiam9lbWFtYTEifQ.Yhr7EfdLxqqwQEfyA2otHjN0yycSx-WPf0OWHrM0Uc4

func DeleteUser(c *fiber.Ctx) error {
	// database.Database.Db.Delete(&User{}, 1)
	// database.Database.Db.Delete(&User{}, 2)
	// database.Database.Db.Delete(&User{}, 3)
	// database.Database.Db.Delete(&User{}, 4)
	// database.Database.Db.Delete(&User{}, 5)
	// database.Database.Db.Delete(&User{}, 6)
	// database.Database.Db.Delete(&User{}, 7)
	// database.Database.Db.Delete(&User{}, 8)
	// database.Database.Db.Delete(&User{}, 9)
	// database.Database.Db.Delete(&User{}, 10)
	// database.Database.Db.Delete(&User{}, 11)
	// database.Database.Db.Delete(&User{}, 12)
	// database.Database.Db.Delete(&User{}, 14)
	// database.Database.Db.Delete(&User{}, 13)
	// database.Database.Db.Delete(&User{}, 15)
	// database.Database.Db.Delete(&User{}, 16)
	return nil
}
