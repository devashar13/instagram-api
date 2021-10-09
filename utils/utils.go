package utils


import (
	"log"
	"context"
	
	"time"
    "net/http"
	"encoding/json"

	"fmt"

	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

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

// func getHash(pwd string){        
//     h := sha256.New()
// 	io.WriteString(h, pwd)

// }

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