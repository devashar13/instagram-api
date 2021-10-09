package main

import (

	"net/http/httptest"
	"testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/devashar13/instagram-api/models"
	"github.com/devashar13/instagram-api/utils"
	"context"
	"log"
	"net/http"
	"fmt"
	"time"
	"bytes"

)

func CreateUser(t *testing.T) {
	var jsonStr = []byte(`{"Name":"TestUser","Email":"test@admin.com","password":"admin"}`)
	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

}




func GetUser(t *testing.T) {
	// uh := newUserHandler()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	id := primitive.NewObjectID()
	user := models.User{ID:id,Name:"Dev",Email:"devashar13@gmail.com",Password:"password"}
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


func checkStatusCode(code int, want int, t *testing.T) {
	if code != want {
		t.Errorf("Wrong status code: got %v want %v", code, want)
	}
}

func checkContentType(r *httptest.ResponseRecorder, t *testing.T) {
	ct := r.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Wrong Content Type: got %v want application/json", ct)
	}
}


// func checkProperties(st student, t *testing.T) {
// 	if st.Name != "John Doe" {
// 		t.Errorf("Name should match: got %v want %v", st.Name, "Peter Doe")
// 	}
// 	if st.Age != 20 {
// 		t.Errorf("Age should match: got %v want %v", st.Age, 20)
// 	}
// }

// func checkBody(body *bytes.Buffer, st student, t *testing.T) {
// 	var students []student
// 	_ = json.Unmarshal(body.Bytes(), &students)
// 	if len(students) != 1 {
// 		t.Errorf("Wrong lenght: got %v want 1", len(students))
// 	}
// 	if students[0] != st {
// 		t.Errorf("Wrong body: got %v want %v", students[0], st)
// 	}
// }

