package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devashar13/instagram-api/models"
	"github.com/devashar13/instagram-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type UserHandler struct {
	sync.Mutex
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.GetUsers(w, r)
	case "POST":
		uh.PostUsers(w, r)
	case "PUT", "PATCH":
	case "DELETE":
	default:
	}
}

// Post Controlls

// User Controls

func (ph *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, err := utils.IdFromUrl(r)
	if err != nil {
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	quickstartDatabase := client.Database("instagramapi")
	userCollection := quickstartDatabase.Collection("users")
	var user bson.M
	if err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		log.Fatal(err)
	}
	utils.RespondWithJSON(w, http.StatusOK, user)

}

func (ph *UserHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi")
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
	var user models.User
	user.ID = primitive.NewObjectID()
	err = json.Unmarshal(body, &user)
	CreateUser(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, user)

}

func CreateUser(user models.User){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	user.Password = utils.GetHash(user.Password)
	quickstartDatabase := client.Database("instagramapi")
	userCollection := quickstartDatabase.Collection("users")
	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	fmt.Println(resultInsertionNumber,insertErr)

}