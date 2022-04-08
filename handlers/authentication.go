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

func getUserByEmail(email *string, c *fiber.Ctx) (*models.User, error) {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	user := &models.User{}

	query := bson.D{{Key: "email", Value: email}}

	err := userCollection.FindOne(c.Context(), query).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed), err
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//database
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	var input LoginInput

	//parser
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
	user, err := getUserByEmail(&email, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err,
		})
	}

	//compare password
	if !checkPasswordHash(password, user.Password) {
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

	//prepare update data to append token
	id, _ := primitive.ObjectIDFromHex(*user.ID)
	user.Token = append(user.Token, t)

	var userToUpdate bson.D
	userToUpdate = append(userToUpdate, bson.E{Key: "updatedAt", Value: time.Now()}) //updatedAt
	userToUpdate = append(userToUpdate, bson.E{Key: "token", Value: user.Token})     //token

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: userToUpdate}}

	//update
	err = userCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't login",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data": fiber.Map{
			"token": t,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	return c.SendString("logout")
}

func Register(c *fiber.Ctx) error {
	//database
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	//new user with user model
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err,
		})
	}

	//check exist user
	_, err := getUserByEmail(&user.Email, c)
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Email is already used",
			"data":    err,
		})
	}

	//hash password
	hash, err := hashPassword(user.Password)
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
	user.Token = append(user.Token, t)

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

func ResetPassword(c *fiber.Ctx) error {
	return c.SendString("reset-password")
}
