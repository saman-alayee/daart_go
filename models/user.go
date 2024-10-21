package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"` // ID mapped to MongoDB _id
    Name  string             `json:"name" bson:"name"`
    Email string             `json:"email" bson:"email"`
}
