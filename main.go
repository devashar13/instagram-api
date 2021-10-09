

package main

import (
	"log"
	"context"

	"time"
    "net/http"
	"encoding/json"

	"fmt"

	"strings"
	"sync"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/devashar13/instagram-api/models"


)

type userHandler struct {
	sync.Mutex
}
type postHandler struct {
	sync.Mutex
}


func (uh *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.get(w, r)
	case "POST":
		uh.post(w, r)
	case "PUT", "PATCH":
	case "DELETE":
	default:
	}
}


func (ph *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.get(w, r)
	case "POST":
		ph.post(w, r)
	case "PUT", "PATCH":
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	case "DELETE":
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}
}

// Post Controlls
func (ph *postHandler) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")   
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var post models.Post
	post.ID = primitive.NewObjectID()
	err = json.Unmarshal(body, &post)
    // getHash(user.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, post)



}

func (ph *postHandler) get(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, err := idFromUrl(r)
	if err != nil {
	}
	if strings.Contains(r.URL.String(),"users"){
			objectID, _ := primitive.ObjectIDFromHex(id)
			client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			err = client.Connect(ctx)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
			}
			defer client.Disconnect(ctx)
			quickstartDatabase := client.Database("instagramapi")
			postCollection := quickstartDatabase.Collection("posts")
			cursor, err := postCollection.Find(ctx, bson.M{"user": objectID})
			if err != nil {
				log.Fatal(err)
			}
			var posts []bson.M
			if err = cursor.All(ctx, &posts); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			respondWithJSON(w, http.StatusCreated, posts)

		// hi
	}else{
		
		objectID, _ := primitive.ObjectIDFromHex(id)
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		defer client.Disconnect(ctx)
		quickstartDatabase := client.Database("instagramapi")
		postCollection := quickstartDatabase.Collection("posts")
		var post bson.M
		if err = postCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post); err != nil {
			log.Fatal(err)
		}
		respondWithJSON(w, http.StatusCreated, post)

	}
	
	

}








// User Controls 

func (ph *userHandler) get(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, err := idFromUrl(r)
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
	respondWithJSON(w,http.StatusOK,user)

}



func (ph *userHandler) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")   
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var user User
	user.ID = primitive.NewObjectID()
	err = json.Unmarshal(body, &user)
    // getHash(user.Password)
	createUser(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)




}








func newUserHandler() *userHandler {
	return &userHandler{}

}

func newPostHandler() *postHandler {
	return &postHandler{}

}






func main(){
	

	uh := newUserHandler()
	ph := newPostHandler()

	port := ":5000"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.Handle("/users", uh)
	http.Handle("/users/", uh)
	http.Handle("/posts", ph)
	http.Handle("/posts/users/", ph)
	http.Handle("/posts/", ph)


	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))



}