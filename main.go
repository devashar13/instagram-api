

package main

import (
	"log"
	"context"
	"io"

	"time"
    "net/http"
	"encoding/json"

	"fmt"

	"strings"
	"sync"
	"go.mongodb.org/mongo-driver/bson"
	"crypto/sha256"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

type userHandler struct {
	sync.Mutex
}
type postHandler struct {
	sync.Mutex
}

type User struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    Name string `json:"name" bson:"name"` 
    Email string `json:"email" bson:"email"` 
    Password string `json:"password" bson:"password"`
}

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	User     primitive.ObjectID `bson:"user"`
	Caption       string             `json:"caption" bson:"caption"`
	Imageurl string            `json:"imageurl" bson:"imageurl"`
	CreatedDate time.Time   `json:"time" bson:"time"`
}




func (uh *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.get(w, r)
	case "POST":
		uh.post(w, r)
	case "PUT", "PATCH":
		uh.put(w, r)
	case "DELETE":
		uh.delete(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}
}


func (ph *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.get(w, r)
	case "POST":
		ph.post(w, r)
	case "PUT", "PATCH":
		// ph.put(w, r)
	case "DELETE":
		// ph.delete(w, r)
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
	var post Post
	post.ID = primitive.NewObjectID()
	err = json.Unmarshal(body, &post)
    // getHash(user.Password)
	createPost(post)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}


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
				log.Fatal(err)
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
				log.Fatal(err)
			}
			fmt.Println(posts)
		// hi
	}else{
		
		objectID, _ := primitive.ObjectIDFromHex(id)
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		quickstartDatabase := client.Database("instagramapi")
		postCollection := quickstartDatabase.Collection("posts")
		var post bson.M
		if err = postCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post); err != nil {
			log.Fatal(err)
		}
		fmt.Println(post)

	}
	
	

}



func createPost(post Post){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
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
	quickstartDatabase := client.Database("instagramapi")
	postCollection := quickstartDatabase.Collection("posts")
	userCollection := quickstartDatabase.Collection("users")
	var user bson.M
	if err = userCollection.FindOne(ctx, bson.M{"_id": post.User}).Decode(&user); err != nil {
		log.Fatal(err)
	}
	if user != nil{
		post.CreatedDate = time.Now()
		resultInsertionNumber, insertErr := postCollection.InsertOne(ctx, post)
		fmt.Println(resultInsertionNumber,insertErr)
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


    // var user User  = json.NewDecoder(r.Body).Decode(&user)  
    // collection := client.Database("GODB").Collection("user") 
    // ctx,_ := context.WithTimeout(context.Background(),      
    //          10*time.Second) 
    // result,err := collection.InsertOne(ctx,user)
    // if err!=nil{     
    //     w.WriteHeader(http.StatusInternalServerError)    
    //     w.Write([]byte(`{"message":"`+err.Error()+`"}`))    
    //     return
    // }    
    // json.NewEncoder(w).Encode(result)


}


func getHash(pwd string){        
    h := sha256.New()
	io.WriteString(h, pwd)

}

func (ph *userHandler) put(w http.ResponseWriter, r *http.Request) {
	

}

func (ph *userHandler) delete(w http.ResponseWriter, r *http.Request) {
	

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}




func idFromUrl(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.String(), "/")
	fmt.Println("hello",parts)
	if len(parts) > 3{
		id := parts[3]
		return id, nil
	}
	id := parts[2]
	return id, nil
}
func newUserHandler() *userHandler {
	return &userHandler{}

}

func newPostHandler() *postHandler {
	return &postHandler{}

}



func createUser(user User){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
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
	quickstartDatabase := client.Database("instagramapi")
	userCollection := quickstartDatabase.Collection("users")
	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	fmt.Println(resultInsertionNumber,insertErr)

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