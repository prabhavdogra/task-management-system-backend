package routes

import (
	"errors"
	"to-do-backend/database"
	"to-do-backend/models"

	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Heading  string `json:"heading"`
	Content  string `json:"content"`
	Progress int    `json:"progress"`
}

func CreateResponseTask(task models.Task) Task {
	return Task{
		ID:       task.ID,
		Heading:  task.Heading,
		Content:  task.Content,
		Progress: task.Progress,
	}
}

func getEmailFromJWT(c *fiber.Ctx) string {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	c.ReqHeaderParser(&authData)
	claims, _ := authenticateHelper(authData.Authorization)
	return claims["email"].(string)
}

func CreateTask(c *fiber.Ctx) error {
	if !AuthenticateRequest(c) {
		return c.Status(400).SendString("JWT Token couldn't be authenticated")
	}
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	var user models.User
	if err := findUserByEmail(getEmailFromJWT(c), &user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	database.Database.Db.Model(&user).Association("Tasks").Append(&task)
	responseTask := CreateResponseTask(task)
	return c.Status(200).JSON(responseTask)
}

func GetAllTasks(c *fiber.Ctx) error {
	if !AuthenticateRequest(c) {
		return c.Status(400).SendString("JWT Token couldn't be authenticated")
	}
	var user models.User
	if err := findUserByEmail(getEmailFromJWT(c), &user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	var tasks []Task
	database.Database.Db.Model(&user).Association("Tasks").Find(&tasks)
	return c.Status(200).JSON(tasks)
}

func FindTaskById(id uint, task *models.Task) error {
	database.Database.Db.Find(&task, "id = ?", id)
	if task.ID == 0 {
		return errors.New("Task does not exist")
	}
	return nil
}

func UpdateTask(c *fiber.Ctx) error {
	if !AuthenticateRequest(c) {
		return c.Status(400).SendString("JWT Token couldn't be authenticated")
	}
	var user models.User
	if err := findUserByEmail(getEmailFromJWT(c), &user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err := FindTaskById(task.ID, &task)
	if err != nil {
		return errors.New("Task doesn't exist")
	}
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Save(&task)
	return c.Status(200).JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid JSON")
	}
	database.Database.Db.Delete(&models.Task{}, id)
	return c.Status(200).SendString("Deleted Successfully")
}
