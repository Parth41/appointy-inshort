package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB : Helper function to connect mongoDB
func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://Parth:NtNgdfNxQgKt1jSn@inshort.rcw2w.mongodb.net/<dbname>?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("appointy").Collection("inshort")

	// create text index to search result from keyword and apply weight to sort result according to weight
	opt := options.Index()
	opt.SetWeights(bson.M{
		"title":    5, // Word matches in the title are weighted 5× standard.
		"subtitle": 3,
		"content":  2,
	})
	index := mongo.IndexModel{Keys: bson.M{
		"title":    "text",
		"subtitle": "text",
		"content":  "text",
	}, Options: opt}
	_, err = collection.Indexes().CreateOne(context.TODO(), index)

	if err != nil {
		log.Println("Could not create text index:", err)
	}

	return collection
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
