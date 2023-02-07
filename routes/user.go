package routes

import (
	"errors"
	"to-do-backend/database"
	"to-do-backend/models"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	// This is not the model, more like a serializer
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
	PhoneNo  string `json:"phone_no"`
}

func CreateResponseUser(user models.User) User {
	return User{
		ID:       user.ID,
		Name:     user.Name,
		EmailID:  user.EmailID,
		Password: user.Password,
		PhoneNo:  user.PhoneNo,
	}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "email = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func findUserByEmail(email string, user *models.User) error {
	database.Database.Db.Find(&user, "email_id = ?", email)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := findUserByEmail(email, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("email")
	var user models.User
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		Name    string `json:"name"`
		EmailID string `json:"email_id"`
		PhoneNo string `json:"phone_no"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.Name = updateData.Name
	user.EmailID = updateData.EmailID
	user.PhoneNo = updateData.PhoneNo

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}
