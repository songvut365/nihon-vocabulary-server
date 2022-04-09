package handlers

import (
	"nihon-vocabulary/configs"
	"nihon-vocabulary/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	vocabulary := new(models.Vocabulary)

	err := c.BodyParser(&vocabulary)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse JSON",
			"data":    err,
		})
	}

	vocabulary.ID = nil
	vocabulary.Owner = "public"
	vocabulary.IsShow = true
	vocabulary.CreatedAt = time.Now()
	vocabulary.UpdatedAt = time.Now()

	result, err := vocabularyCollection.InsertOne(c.Context(), vocabulary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot insert vocabulary",
			"data":    err,
		})
	}

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

func GetVocabularies(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	query := bson.M{"owner": "public"}

	filterCursor, err := vocabularyCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Vocabulary not found",
			"data":    err,
		})
	}

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

func GetVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	vocabulary := &models.Vocabulary{}

	query := bson.D{{Key: "_id", Value: id}}

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

func UpdateVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	vocabulary := new(models.Vocabulary)
	err = c.BodyParser(&vocabulary)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse json",
			"data":    err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

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

	update := bson.D{{Key: "$set", Value: vocabularyToUpdate}}

	err = vocabularyCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot update vocabulary",
			"data":    err,
		})
	}

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

func DeleteVocabulary(c *fiber.Ctx) error {
	vocabularyCollection := configs.MI.DB.Collection(os.Getenv("VOCABULARY_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse id",
			"data":    err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

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
