package routes

import (
	"crypto/sha256"
	"encoding/hex"
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

func CreateJWT(email string, expTime int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = email
	claims["exp"] = expTime

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
	token, err := CreateJWT(user.EmailID, time.Now().Add(time.Minute*30).Unix())
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

func authenticateHelper(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func AuthenticateJWTToken(c *fiber.Ctx) error {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return c.Status(400).SendString("No JWT Token Found")
	}
	claims, validToken := authenticateHelper(authData.Authorization)
	if !validToken {
		return c.Status(400).SendString("Claims not verified")
	}
	if claims["authorized"] == false {
		return c.Status(400).SendString("You are not authorized")
	}
	if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
		return c.Status(400).SendString("Session expired!")
	}
	return c.Status(200).SendString("JWT Token Authenticated")
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
		token, err := CreateJWT(user.EmailID, time.Now().Add(time.Minute*30).Unix())
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
