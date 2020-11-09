package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	// Import Helper and Models package
	"github.com/parth41/inshort/helper"
	"github.com/parth41/inshort/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connection mongoDB with helper class
var collection = helper.ConnectDB()

func main() {

	// Set Route
	http.HandleFunc("/articles", handleArticles)
	http.HandleFunc("/articles/", getArticle)
	http.HandleFunc("/articles/search", searchArticles)

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// Fucntion for handle post and get method of route "/articles"
func handleArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		// Mongodb cursour, pagination - return mongo/option
		cur, err := collection.Find(context.TODO(), bson.M{}, pagination(r).SetSort(bson.M{"_id": -1}))

		if err != nil {
			helper.GetError(err, w)
			return
		}

		// function for get articles list from monogdb
		articles := getArticles(cur)

		// Serialize articles model data and retur
		json.NewEncoder(w).Encode(articles) // encode to serialize process.

	} else if r.Method == "POST" {

		var article models.Article
		article.CreationTimeStamp = time.Now()

		// we decode our body request params
		_ = json.NewDecoder(r.Body).Decode(&article)

		// insert our article model.
		result, err := collection.InsertOne(context.TODO(), article)

		if err != nil {
			helper.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)

	} else {
		message := "Method not allowed"
		fmt.Fprintf(w, message)
		http.Redirect(w, r, "/articles", http.StatusFound)
	}
}

// searchArticles function for handle search route get request
func searchArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Redirect(w, r, "/articles", http.StatusFound)
		}
		filter := bson.M{"$text": bson.M{"$search": query}}

		cur, err := collection.Find(context.TODO(), filter, pagination(r))

		if err != nil {
			helper.GetError(err, w)
			return
		}

		articles := getArticles(cur)

		json.NewEncoder(w).Encode(articles) // encode similar to serialize process.

	} else {
		message := "Method not allowed"
		fmt.Fprintf(w, message)
		http.Redirect(w, r, "/", http.StatusFound)

	}
}

// getArticles function return articles list
func getArticles(cur *mongo.Cursor) []models.Article {

	var articles []models.Article

	// Close the cursor once finished
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var article models.Article
		// & character returns the memory address of the following variable.
		err := cur.Decode(&article) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		articles = append(articles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return articles

}

// Handle get method of route "/articles/<article id>", get single article info from article id
func getArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		articleID := r.URL.Path[len("/articles/"):]
		if articleID != "" {

			w.Header().Set("Content-Type", "application/json")

			var article models.Article

			id, _ := primitive.ObjectIDFromHex(articleID)

			// Declare filtter to get article from ID
			filter := bson.M{"_id": id}
			err := collection.FindOne(context.TODO(), filter).Decode(&article)

			if err != nil {
				helper.GetError(err, w)
				return
			}

			json.NewEncoder(w).Encode(article)

		} else {
			// if id not provided in url redirect to articles route and show all article data
			http.Redirect(w, r, "/articles", http.StatusFound)
		}
	} else {
		message := "Method not allowed"
		fmt.Fprintf(w, message)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func pagination(r *http.Request) *options.FindOptions {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	getOptions := options.Find()
	if page == 0 {
		page = 1
	}
	if limit > 100 || limit == 0 {
		limit = 10
	}
	if page == 1 {
		getOptions.SetSkip(0)
		getOptions.SetLimit(limit)
	} else {
		getOptions.SetLimit(limit)
		getOptions.SetSkip((page - 1) * limit)
	}
	return getOptions
}
