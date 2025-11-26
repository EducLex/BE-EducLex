package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

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
	JaksaCollection          *mongo.Collection
	CategoryCollection       *mongo.Collection
	BidangCollection         *mongo.Collection
)

func ConnectDB() {
	uri := "mongodb+srv://educlexUser:Dewi201202@educlex.fupsgp1.mongodb.net/?retryWrites=true&w=majority"

	tlsConfig := &tls.Config{}
	clientOptions := options.Client().ApplyURI(uri).SetTLSConfig(tlsConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Gagal konek ke MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB tidak bisa diakses:", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")

	// Initialize collections
	UserCollection = client.Database("EducLex").Collection("users")
	QuestionCollection = client.Database("EducLex").Collection("questions")
	ArticleCollection = client.Database("EducLex").Collection("articles")
	TulisanCollection = client.Database("EducLex").Collection("tulisan")
	PeraturanCollection = client.Database("EducLex").Collection("peraturan")
	TokenBlacklistCollection = client.Database("EducLex").Collection("token_blacklist")
	JaksaCollection = client.Database("EducLex").Collection("jaksa")
	CategoryCollection = client.Database("EducLex").Collection("categories")
	BidangCollection = client.Database("EducLex").Collection("bidang")
}
