package handlers

import (
	"nihon-vocabulary/configs"
	"nihon-vocabulary/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByEmail(email *string, c *fiber.Ctx) (*models.User, error) {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	user := &models.User{}

	query := bson.D{{Key: "email", Value: email}}

	err := userCollection.FindOne(c.Context(), query).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(id *primitive.ObjectID, c *fiber.Ctx) (*models.User, error) {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	user := &models.User{}

	query := bson.D{{Key: "_id", Value: id}}

	err := userCollection.FindOne(c.Context(), query).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed), err
}

// Login is a function to authentication user
// @Summary Login
// @Description Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginInput body models.LoginInput true "Login Form"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	//parser
	var input models.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"data":    err,
		})
	}
	email := input.Email
	password := input.Password

	//get user by email
	user, err := GetUserByEmail(&email, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err,
		})
	}

	//compare password
	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid password",
			"data":    err,
		})
	}

	//create jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data": fiber.Map{
			"token": t,
		},
	})
}

// Register is a function to register user
// @Summary Register
// @Description Register
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerInput body models.RegisterInput true "Register Form"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	//new user with user model
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err,
		})
	}

	//check exist user
	_, err := GetUserByEmail(&user.Email, c)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email is already used",
			"data":    err,
		})
	}

	//hash password
	hash, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't hash password",
			"data":    err,
		})
	}

	//create jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//add user to database
	user.ID = nil
	user.Password = hash
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't register",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success register",
		"data": fiber.Map{
			"token": t,
		},
	})
}

// ResetPassword is a function to reset user password
// @Summary Reset Password
// @Description Reset Password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param resetPasswordInput body models.ResetPasswordInput true "Reset Password Form"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/reset-password [post]
func ResetPassword(c *fiber.Ctx) error {
	return c.SendString("reset-password")
}
