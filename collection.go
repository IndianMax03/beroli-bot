package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	collection *mongo.Collection
}

func NewCollection(repo MongoRepository) *Collection {
	collection := repo.CreateCollection(MONGO_DB_NAME, MONGO_COLLECTION_NAME)
	return &Collection{
		collection: collection,
	}
}

func (c *Collection) ExistsByUsername(ctx context.Context, username string) error {
	return c.collection.FindOne(
		ctx, bson.D{
			{Key: USERNAME_FIELD, Value: username},
		},
	).Err()
}
