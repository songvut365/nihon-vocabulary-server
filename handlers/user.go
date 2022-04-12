package handlers

import (
	"nihon-vocabulary/configs"
	"nihon-vocabulary/middlewares"
	"nihon-vocabulary/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GET /api/v1/user
func GetUser(c *fiber.Ctx) error {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	//get id from token
	idFromToken := middlewares.GetIdFromToken(c)
	id, err := primitive.ObjectIDFromHex(idFromToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse user id",
			"data":    err,
		})
	}

	//find user by id
	user := &models.User{}

	query := bson.D{{Key: "_id", Value: id}}

	projection := bson.M{"_id": 0, "password": 0} //exclude fields
	opts := options.FindOne().SetProjection(projection)

	err = userCollection.FindOne(c.Context(), query, opts).Decode(user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success get user",
		"data": fiber.Map{
			"user": user,
		},
	})
}

//PUT /api/v1/user
func UpdateUser(c *fiber.Ctx) error {
	userCollection := configs.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	//get id from token
	idFromToken := middlewares.GetIdFromToken(c)
	id, err := primitive.ObjectIDFromHex(idFromToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse user id",
			"data":    err,
		})
	}

	//parser
	user := new(models.User)
	err = c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse json",
			"data":    err,
		})
	}

	//define vocabulary data to update
	query := bson.D{{Key: "_id", Value: id}}

	var userToUpdate bson.D

	userToUpdate = append(userToUpdate, bson.E{Key: "firstName", Value: user.FirstName})
	userToUpdate = append(userToUpdate, bson.E{Key: "lastName", Value: user.LastName})
	userToUpdate = append(userToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	//update user by id
	update := bson.D{{Key: "$set", Value: userToUpdate}}

	err = userCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot update user",
			"data":    err,
		})
	}

	//find user after updated
	user = &models.User{}

	query = bson.D{{Key: "_id", Value: id}}

	projection := bson.M{"_id": 0, "password": 0} //exclude fields
	opts := options.FindOne().SetProjection(projection)

	err = userCollection.FindOne(c.Context(), query, opts).Decode(user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success update user",
		"data": fiber.Map{
			"user": user,
		},
	})
}

//PUT /api/v1/user/password
func ChangePasswordUser(c *fiber.Ctx) error {
	return c.SendString("change password user")
}

//DELETE /api/v1/user
func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("delete user")
}
