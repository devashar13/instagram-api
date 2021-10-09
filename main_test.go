package main

import (
	"net/http"
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
	"bytes"
	"fmt"
	"time"
)

func GetUser(t *testing.T) {
	uh := newUserHandler()
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
	req, _ := http.NewRequest("GET", "/users"+id.Hex(), nil)
	r := httptest.NewRecorder()
	uh.ServeHTTP(r,req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}

}

func GetPost(t *testing.T) {
	ph := newPostHandler()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	userid := primitive.NewObjectID()
	prodid := primitive.NewObjectID()
	user := models.User{ID:userid,Name:"testprod",Email:"testprod@gmail.com",Password:"password"}
	post := models.Post{ID:prodid,User:userid,Caption:"caption",Imageurl:"https://imgur.com/gallery/no3t9ib",CreatedDate:time.Now()}
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
	result, err := userCollection.InsertOne(ctx, post)
	fmt.Println(resultInsertionNumber,insertErr,result,err)
	req, _ := http.NewRequest("GET", "/posts/"+prodid.Hex(),nil)
	r := httptest.NewRecorder()
	ph.ServeHTTP(r,req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	

}



func CreateUser(t *testing.T) {
	var jsonStr = []byte(`{"Name":"TestUser","Email":"test@admin.com","password":"admin"}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

}

func CreatePost(t *testing.T) {
	ph := newPostHandler()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	userid := primitive.NewObjectID()
	prodid := primitive.NewObjectID()
	user := models.User{ID:userid,Name:"testprod",Email:"testprod@gmail.com",Password:"password"}
	post := models.Post{ID:prodid,User:userid,Caption:"caption",Imageurl:"https://imgur.com/gallery/no3t9ib",CreatedDate:time.Now()}

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
	result, err := userCollection.InsertOne(ctx, post)
	fmt.Println(resultInsertionNumber,insertErr,result,err)
	req, _ := http.NewRequest("GET", "/users/post"+userid.Hex(),nil)
	r := httptest.NewRecorder()
	ph.ServeHTTP(r,req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{ID:`+prodid.Hex()+`,User:`+userid.Hex()+`,Caption:"caption",Imageurl:"https://imgur.com/gallery/no3t9ib",CreatedDate:time.Now()}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}

}

func GetUserPost(t *testing.T) {
	ph := newPostHandler()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	userid := primitive.NewObjectID()
	prodid := primitive.NewObjectID()
	user := models.User{ID:userid,Name:"testprod",Email:"testprod@gmail.com",Password:"password"}
	post := []byte(`{ID:`+prodid.Hex()+`,User:`+userid.Hex()+`,Caption:"caption",Imageurl:"https://imgur.com/gallery/no3t9ib",CreatedDate:time.Now()}`)

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
	result, err := userCollection.InsertOne(ctx, post)
	fmt.Println(resultInsertionNumber,insertErr,result,err)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(post))
	r := httptest.NewRecorder()
	ph.ServeHTTP(r,req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{ID:`+prodid.Hex()+`,User:`+userid.Hex()+`,Caption:"caption",Imageurl:"https://imgur.com/gallery/no3t9ib",CreatedDate:time.Now()}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}

}


