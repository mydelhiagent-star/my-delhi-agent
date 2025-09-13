package databases

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ← PRODUCTION-READY CONNECTION POOL OPTIONS
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).                        // Maximum connections in pool
		SetMinPoolSize(5).                          // Minimum connections in pool
		SetMaxConnIdleTime(30 * time.Minute).       // Close idle connections after 30min
		SetServerSelectionTimeout(5 * time.Second). // Timeout for server selection
		SetConnectTimeout(10 * time.Second).        // Connection timeout
		SetSocketTimeout(30 * time.Second).         // Socket timeout
		SetHeartbeatInterval(10 * time.Second).     // Heartbeat interval
		SetRetryWrites(true).                       // Retry failed writes
		SetRetryReads(true)                         // Retry failed reads

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not reachable:", err)
	}

	log.Println("✅ Connected to MongoDB with optimized connection pool")
	return client
}

// InitMongoDB initializes MongoDB connection and returns database instance
func InitMongoDB() (*mongo.Database, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "delhi_agent"
	}

	client := ConnectMongo(uri)
	db := client.Database(dbName)

	log.Printf("MongoDB database '%s' initialized successfully", dbName)
	return db, nil
}
