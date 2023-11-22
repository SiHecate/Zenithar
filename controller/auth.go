package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"Zenithar/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

/*
Signup:

	Signup, yeni bir kullanıcı kaydı oluşturur.
	İlgili JSON isteği (signupRequest) kullanılarak yeni bir kullanıcı oluşturulur.
	Şifre, güvenli bir şekilde hashlenir.
	Eğer aynı e-posta adresiyle kayıtlı bir kullanıcı varsa hata döner.
*/
func Signup(c *fiber.Ctx) error {
	var signupRequest struct {
		Name             string `json:"name"`
		Email            string `json:"email"`
		Password         string `json:"password"`
		Password_confirm string `json:"password_confirm"`
	}

	signupRequest.Name = strings.TrimSpace(signupRequest.Name)
	signupRequest.Email = strings.TrimSpace(signupRequest.Email)

	if err := c.BodyParser(&signupRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	if signupRequest.Password != signupRequest.Password_confirm {
		return c.Status(400).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	var existingUser models.User
	database.Conn.Where("email = ?", signupRequest.Email).First(&existingUser)

	if existingUser.ID != 0 {
		return c.Status(400).JSON(fiber.Map{"error": "user already exists"})
	}

	var err error
	signupRequest.Password, err = utils.GenerateHashPassword(signupRequest.Password)
	if err != nil {
		return err
	}

	response := models.User{
		Name:     signupRequest.Name,
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
		Role:     "User",
	}

	if err := database.Conn.Create(&response).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User creation failed"})
	}

	return c.JSON(response)
}

/*
Login:

	Login, kullanıcı girişi yapar ve JWT (Json Web Token) üretilir.
	Kullanıcının adı ve şifresi alınır.
	Veritabanında bu bilgilerle kullanıcı aranır.
	Şifre kontrolü yapılır ve eğer doğruysa JWT üretilir.
	Üretilen JWT, bir HTTP only cookie ile kullanıcıya gönderilir.
*/
func Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingUser models.User
	database.Conn.Where("name = ?", loginRequest.Name).First(&existingUser)

	if existingUser.ID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "User not found"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	// JWT ************

	expirationTime := time.Now().Add(12 * time.Hour)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(existingUser.ID)),
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var jwtKey = []byte("my_secret_key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not generate token"})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Minute * 60),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "user login success!",
		"token":   tokenString, // Include the token in the response
	})
}

/*
Logout:

	Logout, kullanıcının çıkışını sağlar.
	Kullanıcının tarayıcısındaki JWT içeren cookie silinir.
*/
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
