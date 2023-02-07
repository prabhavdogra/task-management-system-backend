package routes

import (
	"time"
	"to-do-backend/database"
	"to-do-backend/models"

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
	return c.Status(400).JSON(fiber.Map{
		"token": token,
	})
}

// func fn(c *fiber.Ctx, x string) error {
// 	return c.JSON(fiber.Map{"k": x})
// }

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

func Logout(c *fiber.Ctx) error {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return c.Status(400).SendString("No JWT Token Found")
	}
	token := models.BlacklistedToken{}
	result := database.Database.Db.Where("blacklisted_token = ?", authData.Authorization).First(&token)
	if result.RowsAffected == 1 {
		return c.Status(400).SendString("Token already blacklisted")
	} else {
		token.BlacklistedToken = authData.Authorization
		database.Database.Db.Save(&token)
		return c.Status(200).SendString("Added to blacklisted tokens")
	}
}

func DeleteUser(c *fiber.Ctx) error {
	database.Database.Db.Delete(&models.Task{}, 1)
	database.Database.Db.Delete(&models.Task{}, 2)
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
