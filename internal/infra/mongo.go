package infra

import (
	"context"
	"time"

	global "github.com/IndianMax03/beroli-bot/internal/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout = 2 * time.Second
	pingTimeout       = 1 * time.Second
)

type MongoRepository struct {
	client *mongo.Client
}

func NewConnection(ctx context.Context) (*MongoRepository, error) {
	connectCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(global.MONGO_URL)
	if global.MONGO_USER != "" && global.MONGO_PASSWORD != "" {
		clientOptions.SetAuth(options.Credential{
			Username: global.MONGO_USER,
			Password: global.MONGO_PASSWORD,
		})
	}
	clientOptions.SetMaxPoolSize(0)

	mongoClient, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err = mongoClient.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoRepository{
		client: mongoClient,
	}, nil
}

func (r *MongoRepository) CloseConnection(ctx context.Context) error {
	connectCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()
	return r.client.Disconnect(connectCtx)
}

func (r *MongoRepository) CreateCollection(dbName, collectionName string) *mongo.Collection {
	return r.client.Database(dbName).Collection(collectionName)
}
