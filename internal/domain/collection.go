package domain

import (
	"context"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	global "github.com/IndianMax03/beroli-bot/internal/global"
	infra "github.com/IndianMax03/beroli-bot/internal/infra"
)

type Collection struct {
	Collection MongoCollection
}

func NewCollection(repo infra.MongoRepository) *Collection {
	collection := repo.CreateCollection(global.MONGO_DB_NAME, global.MONGO_COLLECTION_NAME)
	return &Collection{
		Collection: collection,
	}
}

func (c *Collection) CreateUser(ctx context.Context, u *User) {
	c.Collection.InsertOne(ctx, u)
}

func (c *Collection) UpdateUser(ctx context.Context, u *User) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": u.Username,
		},
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, u).Err()
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

	err := c.Collection.FindOne(ctx, filter).Decode(&user)
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

	return c.Collection.FindOne(ctx, filter).Err()
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
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
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
	err := c.Collection.FindOne(ctx, filter, &opt).Decode(&u)
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
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
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
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
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
	err := c.Collection.FindOne(ctx, filter, &opt).Decode(&u)
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
	err := c.Collection.FindOne(ctx, filter, &opt).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u.Issue, nil
}

func (c *Collection) UpdateSummaryIssue(ctx context.Context, username, text string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$set": bson.M{
			ISSUE_SUMMARY_FIELD: text,
		},
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) UpdateDescriptionIssue(ctx context.Context, username, text string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$set": bson.M{
			ISSUE_DESCRIPTION_FIELD: text,
		},
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) AppendAttachmentIssue(ctx context.Context, username string, attachmentID string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$push": bson.M{
			ISSUE_ATTACHMENTS_FIELD: attachmentID,
		},
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) AppendDescriptionAttachmentIssue(ctx context.Context, username string, descriptionAttachmentID string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$push": bson.M{
			ISSUE_DESCRIPTION_ATTACHMENTS_FIELD: descriptionAttachmentID,
		},
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, update).Err()
}

func (c *Collection) AppendTagIssue(ctx context.Context, username string, tags []string) error {
	filter := bson.M{
		USERNAME_FIELD: bson.M{
			"$eq": username,
		},
	}
	update := bson.M{
		"$push": bson.M{
			ISSUE_TAGS_FIELD: bson.M{
				"$each": tags,
			},
		},
	}

	_, err := c.Collection.UpdateOne(ctx, filter, update)
	return err
}
