package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectUrl := "mongodb+srv://" + dbUsername + ":" + dbPassword + "@testcluster.nhs5bd3.mongodb.net/?retryWrites=true&w=majority"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(connectUrl).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection("users")
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var users []primitive.M

	for cursor.Next(ctx) {
		var user bson.M
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	defer cursor.Close(ctx)

	r := gin.Default()
	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"users": users,
		})
	})
	r.Run()
}
