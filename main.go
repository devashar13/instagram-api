

package main

import (
	"log"
	"context"
	"time"
    "net/http"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

type userHandler struct {
	sync.Mutex
}

type User struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    Name string `json:"name" bson:"name"` 
    Email string `json:"email" bson:"email"` 
    Password string `json:"password" bson:"password"`
}


func (ph *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.get(w, r)
	case "POST":
		ph.post(w, r)
	case "PUT", "PATCH":
		ph.put(w, r)
	case "DELETE":
		ph.delete(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}
}


func (ph *userHandler) get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Func")

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
	createUser(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}


    // var user User  = json.NewDecoder(r.Body).Decode(&user)  
    // user.Password = getHash([]byte(user.Password)) 
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


// func getHash(pwd []byte) string {        
//     hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)          
//     if err != nil {
//        log.Println(err)
//     }
//     return string(hash)
// }

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




func idFromUrl(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		return 0, errors.New("not found")
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, errors.New("not found")
	}
	return id, nil
}

func newProductHandler() *userHandler {
	return &userHandler{}

}

func createUser(user User){
	fmt.Println(user)
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
	

	ph := newProductHandler()
	port := ":5000"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.Handle("/products", ph)
	http.Handle("/user", ph)
	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))



}