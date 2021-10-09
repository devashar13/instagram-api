package utils


import (
	"log"
	"context"

	"time"
    "net/http"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"encoding/hex"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/devashar13/instagram-api/models"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

func CreatePost(post models.Post){
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

func GetHash(pwd string) string{        
	sha256Bytes := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(sha256Bytes[:])

}





func IdFromUrl(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.String(), "/")
	// fmt.Println("hello",parts)
	if len(parts) > 3{

		if strings.Contains(parts[3], "?"){
			id := strings.Split(parts[3], "?")
			fmt.Println("hi",id[0])
			return id[0], nil
		}
		id := parts[3]
		return id, nil
	}
	id := parts[2]
	return id, nil
}

func CreateUser(user models.User){
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
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


