package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserCollection           *mongo.Collection
	QuestionCollection       *mongo.Collection
	ArticleCollection        *mongo.Collection
	TulisanCollection        *mongo.Collection
	PeraturanCollection      *mongo.Collection
	TokenBlacklistCollection *mongo.Collection
)

func ConnectDB() {
	uri := "mongodb+srv://educlexUser:Dewi201202@educlex.fupsgp1.mongodb.net/?retryWrites=true&w=majority&appName=EducLex"

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("❌ Gagal konek ke MongoDB:", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB tidak bisa diakses:", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")

	UserCollection = client.Database("EducLex").Collection("users")
	QuestionCollection = client.Database("EducLex").Collection("questions")
	ArticleCollection = client.Database("articles").Collection("articles")
	TulisanCollection = client.Database("EducLex").Collection("tulisan")
	PeraturanCollection = client.Database("EducLex").Collection("peraturan")
	TokenBlacklistCollection = client.Database("EducLex").Collection("token_blacklist")
}
