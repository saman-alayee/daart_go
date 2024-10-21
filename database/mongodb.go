package database

import (
    "context"
    "log"
    "os"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Connect initializes the MongoDB connection
func Connect() {
    var err error
    uri := os.Getenv("MONGODB_CONNECTION")
    if uri == "" {
        log.Fatal("MONGODB_CONNECTION environment variable is not set")
    }

    clientOptions := options.Client().ApplyURI(uri)
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the database to verify connection
    if err := client.Ping(context.Background(), nil); err != nil {
        log.Fatal("MongoDB connection failed:", err)
    }

    log.Println("Connected to MongoDB!")
}

// GetCollection returns the collection from the database
func GetCollection(name string) *mongo.Collection {
    dbName := os.Getenv("MONGODB_NAME")
    if dbName == "" {
        log.Fatal("MONGODB_NAME environment variable is not set")
    }

    return client.Database(dbName).Collection(name)
}
