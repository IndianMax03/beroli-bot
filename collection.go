package main

import (
	"context"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *Collection) CreateUser(ctx context.Context, u *User) {
	c.collection.InsertOne(ctx, u)
}

func (c *Collection) UpdateUser(ctx context.Context, u *User) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": u.Username,
		},
	}
	return c.collection.FindOneAndUpdate(ctx, filter, u).Err()
}

func (c *Collection) ResetUser(ctx context.Context, username string) error {
	user := User{
		Username: username,
		State:    NIL_STATE,
	}
	return c.UpdateUser(ctx, &user)
}

func (c *Collection) GetUser(ctx context.Context, username string) (*User, error) {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	var user User

	err := c.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Collection) ExistsUser(ctx context.Context, username string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}

	return c.collection.FindOne(ctx, filter).Err()
}

func (c *Collection) UpdateStateUser(ctx context.Context, username, state string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$set": bson.M{STATE_FIELD: state},
	}
	return c.collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) GetStateUser(ctx context.Context, username string) (string, error) {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	opt := options.FindOneOptions{
		Projection: bson.M{
			STATE_FIELD: 1,
		},
	}
	var u User
	err := c.collection.FindOne(ctx, filter, &opt).Decode(&u)
	if err != nil {
		return "", err
	}
	return u.State, nil
}

func (c *Collection) ClearIssue(ctx context.Context, username string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$set": bson.M{
			ISSUE_FIELD: NewDefaultIssue(),
		},
	}
	return c.collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) AppendDataIssue(ctx context.Context, username string, data *Issue) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$push": bson.M{
			ISSUES_FIELD: data,
		},
	}
	return c.collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) GetIssues(ctx context.Context, username string) ([]Issue, error) {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	opt := options.FindOneOptions{
		Projection: bson.M{
			ISSUES_FIELD: 1,
		},
	}
	var u User
	err := c.collection.FindOne(ctx, filter, &opt).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u.Issues, nil
}

func (c *Collection) GetIssue(ctx context.Context, username string) (*model.IssueCreateRequest, error) {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	opt := options.FindOneOptions{
		Projection: bson.M{
			ISSUE_FIELD: 1,
		},
	}
	var u User
	err := c.collection.FindOne(ctx, filter, &opt).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u.Issue, nil
}
