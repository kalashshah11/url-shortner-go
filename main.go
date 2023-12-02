package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type shortenedURL struct {
	ID  string `bson:"_id", json:"id"`
	URL string `bson:"url", json:"url"`
}

var collection *mongo.Collection
var ctx context.Context

func getUrlFromId(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{{"_id", id}}
	var result shortenedURL
	collection.FindOne(ctx, filter).Decode(&result)
	c.IndentedJSON(http.StatusOK, result)
}

func shortenUrl(c *gin.Context) {
	var newUrl shortenedURL

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newUrl); err != nil {
		log.Println("Hello1")
		return
	}
	newUrl = shortenedURL{ID: string(uuid.New().String()), URL: newUrl.URL}
	log.Println(newUrl, "Hello2")
	// Add the new album to the slice.
	result, err := collection.InsertOne(ctx, newUrl)

	if err == nil {
		log.Println("Hello3", result, err)
		c.IndentedJSON(http.StatusCreated, result.InsertedID)
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, err)
}

func main() {
	router := gin.Default()
	router.POST("/shorten-url", shortenUrl)
	router.GET("/original-url/:id", getUrlFromId)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection = client.Database("url-shortner").Collection("urlmap")
	if err != nil {
		log.Fatal(err)
	}
	router.Run("localhost:8080")

	defer client.Disconnect(ctx)
}
