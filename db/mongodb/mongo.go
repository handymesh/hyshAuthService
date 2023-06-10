package mongodb

import (
	"context"
	"github.com/handymesh/handy_authService/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	Session *mongo.Client
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToMongo() {
	// Get configuration
	MONGO_URL := utils.Getenv("MONGO_URL", "mongodb://localhost/auth")
	clientOptions := options.Client().ApplyURI(MONGO_URL)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic("Failed connect to Mongo", err)
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
	}

	log.Info("Success connect to MongoDB")
	Session = client
}
