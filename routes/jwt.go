package routes

import (
	"time"
	"to-do-backend/database"
	"to-do-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var sampleSecretKey = []byte("GoLinuxCloudKey")

func CreateJWT(email string, expTime int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = expTime

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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

func IsTokenBlackedOut(c *fiber.Ctx) bool {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return false
	}
	token := authData.Authorization
	var blacktoken models.BlacklistedToken
	result := database.Database.Db.Where("blacklisted_token = ?", token).First(&blacktoken)
	if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func verifyClaims(claims jwt.MapClaims) (int, string) {
	if claims["authorized"] == false {
		return 400, "You are not authorized"
	}
	if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
		return 400, "Session expired!"
	}
	return 200, "JWT Token Authenticated"
}

func AuthenticateRequest(c *fiber.Ctx) bool {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil || IsTokenBlackedOut(c) {
		return false
	}
	claims, validToken := authenticateHelper(authData.Authorization)
	if !validToken {
		return false
	}
	statusCode, _ := verifyClaims(claims)
	return statusCode == 200
}

func AuthenticateJWTToken(c *fiber.Ctx) error {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return c.SendString(err.Error())
	}
	JWTToken := authData.Authorization
	claims, validToken := authenticateHelper(JWTToken)
	if !validToken || IsTokenBlackedOut(c) {
		return c.Status(201).SendString("Invalid Token")
	}
	statusCode, message := verifyClaims(claims)
	return c.Status(statusCode).SendString(message)
}

func FindUserByJWT(c *fiber.Ctx, user *models.User) error {
	type AuthData struct {
		Authorization string `json:"authorization"`
	}
	var authData AuthData
	err := c.ReqHeaderParser(&authData)
	if err != nil {
		return c.SendString(err.Error())
	}
	claims, validToken := authenticateHelper(authData.Authorization)
	if !validToken || IsTokenBlackedOut(c) {
		return c.SendString(err.Error())
	}
	findUserByEmail(claims["email"].(string), user)
	return nil
}
