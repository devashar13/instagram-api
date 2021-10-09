
package models

import (

    "go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
    ID            primitive.ObjectID `bson:"_id"`
    First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
    Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
    Password      *string            `json:"Password" validate:"required,min=6""`
    Email         *string            `json:"email" validate:"email,required"`
    Phone         *string            `json:"phone" validate:"required"`
    User_id       string             `json:"user_id"`
}
