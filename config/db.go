package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserCollection      *mongo.Collection
	QuestionCollection  *mongo.Collection
	ArticleCollection   *mongo.Collection
	TulisanCollection   *mongo.Collection
	PeraturanCollection *mongo.Collection
)

func ConnectDB() {
	uri := "mongodb+srv://educlexUser:Dewi201202@educlex.fupsgp1.mongodb.net/?retryWrites=true&w=majority&appName=EducLex"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")

	UserCollection = client.Database("EducLex").Collection("users")
	QuestionCollection = client.Database("EducLex").Collection("questions")
	ArticleCollection = client.Database("articles").Collection("articles")
	TulisanCollection = client.Database("EducLex").Collection("tulisan")
	PeraturanCollection = client.Database("peraturan").Collection("peraturan")

}
