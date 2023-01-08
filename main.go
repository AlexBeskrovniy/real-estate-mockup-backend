package main

import (
	"log"
	"os"

	"github.com/AlexBlacksmith/real-estate-mockup-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	client, ctx, cancel, err := db.Connect(connectUrl)

	defer cancel()

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

	port := os.Getenv("PORT")
	r := gin.Default()
	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"users": users,
		})
	})
	r.Run(":" + port)
}
