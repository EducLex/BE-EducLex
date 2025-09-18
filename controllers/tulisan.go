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

func CreateTulisan(c *gin.Context) {
	var tulisan models.Tulisan   
	if err := c.ShouldBindJSON(&tulisan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tulisan.ID = primitive.NewObjectID()
	tulisan.Created = time.Now()

	_, err := config.ArticleCollection.InsertOne(context.TODO(), tulisan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tulisan"})
		return
	}

	c.JSON(http.StatusOK, tulisan)
}

func GetTulisans(c *gin.Context) {
	cursor, err := config.ArticleCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(context.TODO())

	var tulisans []models.Tulisan   
	if err := cursor.All(context.TODO(), &tulisans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
		return
	}

	c.JSON(http.StatusOK, tulisans)
}
