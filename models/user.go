package models

import (
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// User represents a user in the system
type User struct {
    ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
    PublisherID string           `json:"publisherId" bson:"publisherId"`
    Username string             `json:"username" bson:"username"`
    Password string             `json:"password" bson:"password"`
}

// CheckIfFieldExists checks if a field exists in the collection
func CheckIfFieldExists(collection *mongo.Collection, field string, value string) (bool, error) {
    filter := bson.M{field: value}
    var result User
    err := collection.FindOne(nil, filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return false, nil // No document found, the field doesn't exist
        }
        return false, err // Other errors, e.g., connection issues
    }
    return true, nil // The field exists
}
