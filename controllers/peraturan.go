package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create
func CreatePeraturan(c *gin.Context) {
	var peraturan models.Peraturan
	if err := c.ShouldBindJSON(&peraturan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	peraturan.ID = primitive.NewObjectID()
	peraturan.Created = time.Now()

	_, err := config.PeraturanCollection.InsertOne(context.TODO(), peraturan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan peraturan"})
		return
	}

	c.JSON(http.StatusOK, peraturan)
}

// Get All
func GetPeraturans(c *gin.Context) {
	cursor, err := config.PeraturanCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(context.TODO())

	var peraturans []models.Peraturan
	if err := cursor.All(context.TODO(), &peraturans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
		return
	}

	c.JSON(http.StatusOK, peraturans)
}

// Update
func UpdatePeraturan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var updateData models.Peraturan
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"judul": updateData.Judul,
			"pasal": updateData.Pasal,
		},
	}

	_, err = config.PeraturanCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update peraturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil diupdate"})
}

// Delete
func DeletePeraturan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	_, err = config.PeraturanCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus peraturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil dihapus"})
}
