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
)

// Create Vocabulary is a function to get vocabulary from database
// @Summary Create Vocabulary
// @Description Create Vocabulary
// @Tags Vocabulary
// @Accept json
// @Produce json
// @Param vocabularyInput body models.VocabularyInput true "Vocabulary ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vocabulary [post]
func CreateVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	//parser
	vocabulary := new(models.Vocabulary)

	err := c.BodyParser(&vocabulary)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse JSON",
			"data":    err,
		})
	}

	//get id from token
	idFromToken := middlewares.GetIdFromToken(c)

	//create new vocabulary
	vocabulary.ID = nil
	vocabulary.Owner = idFromToken
	vocabulary.IsShow = true
	vocabulary.CreatedAt = time.Now()
	vocabulary.UpdatedAt = time.Now()

	//insert new vocabulary
	result, err := vocabularyCollection.InsertOne(c.Context(), vocabulary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot insert vocabulary",
			"data":    err,
		})
	}

	//find vocabulary after created
	newVocabulary := &models.Vocabulary{}

	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	vocabularyCollection.FindOne(c.Context(), query).Decode(newVocabulary)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success create vocabulary",
		"data": fiber.Map{
			"vocabulary": newVocabulary,
		},
	})
}

// Get Vocabularies is a function to get all vocabulary from database
// @Summary Get Vocabularies
// @Description Get Vocabularies
// @Tags Vocabulary
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vocabulary/ [get]
func GetVocabularies(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	//get id from token
	idFromToken := middlewares.GetIdFromToken(c)

	//find all vocabulary by owner
	query := bson.D{{Key: "owner", Value: idFromToken}}

	filterCursor, err := vocabularyCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Vocabulary not found",
			"data":    err,
		})
	}

	//push each vocabulary to array
	var vocabularies []bson.M
	if err = filterCursor.All(c.Context(), &vocabularies); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot insert each vocabulary",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success get vocabularies",
		"data": fiber.Map{
			"vocabularies": vocabularies,
		},
	})
}

// Get Vocabulary is a function to get vocabulary from database
// @Summary Get Vocabulary
// @Description Get Vocabulary
// @Tags Vocabulary
// @Accept json
// @Produce json
// @Param id path int true "Vocabulary ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vocabulary/{id} [get]
func GetVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	//get id from params
	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	//get user id from token
	userIdFromToken := middlewares.GetIdFromToken(c)

	//find vocabulary by id
	vocabulary := &models.Vocabulary{}

	query := bson.D{{Key: "_id", Value: id}, {Key: "owner", Value: userIdFromToken}}

	err = vocabularyCollection.FindOne(c.Context(), query).Decode(vocabulary)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Vocabualry not found",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success get vocabulary",
		"data": fiber.Map{
			"vocabularies": vocabulary,
		},
	})
}

// Update Vocabulary is a function to update vocabulary to database
// @Summary Update Vocabulary
// @Description Update Vocabulary
// @Tags Vocabulary
// @Accept json
// @Produce json
// @Param vocabularyInput body models.VocabularyInput true "Vocabulary ID"
// @Param id path int true "Vocabulary ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vocabulary/{id} [put]
func UpdateVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	//get id from params
	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	//parser
	vocabulary := new(models.Vocabulary)
	err = c.BodyParser(&vocabulary)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse json",
			"data":    err,
		})
	}

	//get user id from token
	userIdFromToken := middlewares.GetIdFromToken(c)

	//define vocabulary data to update
	query := bson.D{{Key: "_id", Value: id}, {Key: "owner", Value: userIdFromToken}}

	var vocabularyToUpdate bson.D

	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "japanese", Value: vocabulary.Japanese})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "thai", Value: vocabulary.Thai})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "english", Value: vocabulary.English})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "examples", Value: vocabulary.Examples})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "image", Value: vocabulary.Image})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "voice", Value: vocabulary.Voice})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "type", Value: vocabulary.Type})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "tags", Value: vocabulary.Tags})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "isShow", Value: vocabulary.IsShow})
	vocabularyToUpdate = append(vocabularyToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	//update vocabulary by id
	update := bson.D{{Key: "$set", Value: vocabularyToUpdate}}

	err = vocabularyCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot update vocabulary",
			"data":    err,
		})
	}

	//find vocabulary after updated
	vocabulary = &models.Vocabulary{}
	vocabularyCollection.FindOne(c.Context(), query).Decode(vocabulary)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success update vocabulary",
		"data": fiber.Map{
			"vocabularies": vocabulary,
		},
	})
}

// Delete Vocabulary is a function to delete vocabulary from database
// @Summary Delete Vocabulary
// @Description Delete Vocabulary
// @Tags Vocabulary
// @Accept json
// @Produce json
// @Param id path int true "Vocabulary ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vocabulary/{id} [delete]
func DeleteVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	//get id from params
	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	//get id from token
	idFromToken := middlewares.GetIdFromToken(c)

	//delete vocabulary by id
	query := bson.D{{Key: "_id", Value: id}, {Key: "owner", Value: idFromToken}}

	err = vocabularyCollection.FindOneAndDelete(c.Context(), query).Err()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot delete vocabulary",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success delete vocabulary",
		"data": fiber.Map{
			"id": id,
		},
	})
}
