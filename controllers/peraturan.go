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

// Tambah peraturan baru
func CreatePeraturan(c *gin.Context) {
    var input models.Peraturan
    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.CreatedAt = time.Now()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    res, err := config.PeraturanCollection.InsertOne(ctx, input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan peraturan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil ditambahkan", "id": res.InsertedID})
}

// Ambil semua peraturan
func GetPeraturan(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := config.PeraturanCollection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
        return
    }
    defer cursor.Close(ctx)

    var results []models.Peraturan
    if err := cursor.All(ctx, &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
        return
    }

    c.JSON(http.StatusOK, results)
}

// Hapus peraturan by ID
func DeletePeraturan(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = config.PeraturanCollection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil dihapus"})
}
