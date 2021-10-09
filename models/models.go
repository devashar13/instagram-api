package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	User        primitive.ObjectID `bson:"user"`
	Caption     string             `json:"caption" bson:"caption"`
	Imageurl    string             `json:"imageurl" bson:"imageurl"`
	CreatedDate time.Time          `json:"time" bson:"time"`
}
