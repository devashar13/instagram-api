package controllers

import (
	"context"
	"encoding/json"
	"github.com/devashar13/instagram-api/models"
	"github.com/devashar13/instagram-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"strconv"
	"net/http"
	"strings"
	"fmt"
	"sync"
	"time"
)

type PostHandler struct {
	sync.Mutex
}

func (ph *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.HandleGet(w, r)
	case "POST":
		ph.HandlePost(w, r)
	case "PUT", "PATCH":
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	case "DELETE":
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}
}

func (ph *PostHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var post models.Post
	post.ID = primitive.NewObjectID()
	err = json.Unmarshal(body, &post)
	// getHash(user.Password)
	utils.CreatePost(post)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, post)

}

func (ph *PostHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, err := utils.IdFromUrl(r)
	if err != nil {
	}
	if strings.Contains(r.URL.String(), "users") {
		query := r.URL.Query()
		filters := query["limit"] 
		fmt.Println(filters)
		if len(filters) == 0 {
			fmt.Println("filters not present")
    }
		objectID, _ := primitive.ObjectIDFromHex(id)
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}
		defer client.Disconnect(ctx)
		quickstartDatabase := client.Database("instagramapi")
		postCollection := quickstartDatabase.Collection("posts")
		findOptions := options.Find()
		limit,_ := strconv.ParseInt(filters[0],10,64)
		findOptions.SetLimit(limit)
		cursor, err := postCollection.Find(ctx, bson.M{"user": objectID},findOptions)

		if err != nil {
			log.Fatal(err)
		}
		var posts []bson.M
		if err = cursor.All(ctx, &posts); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		utils.RespondWithJSON(w, http.StatusCreated, posts)

		// hi
	} else {

		objectID, _ := primitive.ObjectIDFromHex(id)
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		defer client.Disconnect(ctx)
		quickstartDatabase := client.Database("instagramapi")
		postCollection := quickstartDatabase.Collection("posts")
		var post bson.M
		if err = postCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post); err != nil {
			log.Fatal(err)
		}
		utils.RespondWithJSON(w, http.StatusCreated, post)

	}

}
