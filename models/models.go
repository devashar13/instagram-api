
package models

import (

    "go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
    ID            primitive.ObjectID `bson:"_id"`
    Name string `json:"name" bson:"name"` 
    Email string `json:"email" bson:"email"` 
    Password string `json:"password" bson:"password"`
}
