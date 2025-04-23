package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/IndianMax03/beroli-bot/mocks"
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

var (
	userStub             = &User{}
	usernameStub         = "stub"
	textStub             = "stub"
	stringArrayStub      = []string{}
	userWithUsernameStub = &User{
		Username: usernameStub,
	}
	userWithUsernameAndNilStateStub = &User{
		Username: usernameStub,
		State:    NIL_STATE,
	}
	userWithIssuesStub = &User{
		Issues: []Issue{{}, {}},
	}
	userWithIssueCreateRequestStub = &User{
		Issue: &model.IssueCreateRequest{},
	}
	issueStub                              = &Issue{}
	errStub                                = errors.New("stub error")
	userWithUsernameAndNilStateMatcherStub = userMatcher{
		want: userWithUsernameAndNilStateStub,
	}
)

type userMatcher struct {
	want *User
}

type issueMatcher struct {
	want *Issue
}

func (m userMatcher) Matches(other interface{}) bool {
	usr, ok := other.(*User)
	if !ok {
		return false
	}
	return usr.Username == m.want.Username && usr.State == m.want.State
}

func (m userMatcher) String() string {
	return fmt.Sprintf("(Username=%s; State=%s)", m.want.Username, m.want.State)
}

func TestCollectionCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)

	mockCollection.
		EXPECT().
		InsertOne(gomock.Any(), userStub).
		Return(&mongo.InsertOneResult{}, nil)

	collection := &Collection{collection: mockCollection}
	collection.CreateUser(context.TODO(), userStub)
}

func TestCollectionUpdateUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithUsernameStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), userWithUsernameStub).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateUser(context.TODO(), userWithUsernameStub)
	assert.NoError(t, err)
}

func TestCollectionUpdateUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), userWithUsernameStub).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateUser(context.TODO(), userWithUsernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionUpdateUserNegativeErrNilDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(nil, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), userWithUsernameStub).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateUser(context.TODO(), userWithUsernameStub)
	assert.ErrorIs(t, err, mongo.ErrNilDocument)
}

func TestCollectionResetUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), userWithUsernameAndNilStateMatcherStub).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ResetUser(context.TODO(), usernameStub)
	assert.NoError(t, err)
}

func TestCollectionResetUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), userWithUsernameAndNilStateMatcherStub).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ResetUser(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionGetUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithUsernameStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	resultUser, err := collection.GetUser(context.TODO(), usernameStub)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, userWithUsernameStub, resultUser)
}

func TestCollectionGetUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResultWithErr := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any()).
		Return(mockResultWithErr)

	collection := &Collection{collection: mockCollection}
	_, err := collection.GetUser(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionExistsUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithUsernameStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ExistsUser(context.TODO(), usernameStub)
	assert.NoError(t, err)
}

func TestCollectionExistsUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ExistsUser(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionUpdateStateUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateStateUser(context.TODO(), usernameStub, CREATING_STATE)
	assert.NoError(t, err)
}

func TestCollectionUpdateStateUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateStateUser(context.TODO(), usernameStub, CREATING_STATE)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionGetStateUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithUsernameAndNilStateStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	state, err := collection.GetStateUser(context.TODO(), usernameStub)
	assert.NoError(t, err)
	assert.Equal(t, userWithUsernameAndNilStateStub.State, state)
}

func TestCollectionGetStateUserNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	_, err := collection.GetStateUser(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionClearIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ClearIssue(context.TODO(), usernameStub)
	assert.NoError(t, err)
}

func TestCollectionClearIssueNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.ClearIssue(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionAppendDataIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendDataIssue(context.TODO(), usernameStub, issueStub)
	assert.NoError(t, err)
}

func TestCollectionAppendDataIssueNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendDataIssue(context.TODO(), usernameStub, issueStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionGetIssuesPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithIssuesStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	iss, err := collection.GetIssues(context.TODO(), usernameStub)
	assert.NoError(t, err)
	assert.Equal(t, userWithIssuesStub.Issues, iss)
}

func TestCollectionGetIssuesNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	_, err := collection.GetIssues(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionGetIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&userWithIssueCreateRequestStub, nil, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	is, err := collection.GetIssue(context.TODO(), usernameStub)
	assert.NoError(t, err)
	assert.Equal(t, userWithIssueCreateRequestStub.Issue, is)
}

func TestCollectionGetIssueNegativeStubError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	_, err := collection.GetIssue(context.TODO(), usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionUpdateSummaryIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateSummaryIssue(context.TODO(), usernameStub, textStub)
	assert.NoError(t, err)
}

func TestCollectionUpdateSummaryIssueNegativeErrorStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateSummaryIssue(context.TODO(), usernameStub, textStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionUpdateDescriptionIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateDescriptionIssue(context.TODO(), usernameStub, textStub)
	assert.NoError(t, err)
}

func TestCollectionUpdateDescriptionIssueNegativeErrorStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.UpdateDescriptionIssue(context.TODO(), usernameStub, textStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionAppendAttachmentIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendAttachmentIssue(context.TODO(), usernameStub, textStub)
	assert.NoError(t, err)
}

func TestCollectionAppendAttachmentIssueNegativeErrorStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendAttachmentIssue(context.TODO(), usernameStub, textStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionAppendDescriptionAttachmentIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, nil, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendDescriptionAttachmentIssue(context.TODO(), usernameStub, textStub)
	assert.NoError(t, err)
}

func TestCollectionAppendDescriptionAttachmentIssueNegativeErrorStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := mongo.NewSingleResultFromDocument(&bson.M{}, errStub, nil)
	mockCollection.
		EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendDescriptionAttachmentIssue(context.TODO(), usernameStub, textStub)
	assert.ErrorIs(t, err, errStub)
}

func TestCollectionAppendTagIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockResult := &mongo.UpdateResult{}
	mockCollection.
		EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockResult, nil)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendTagIssue(context.TODO(), usernameStub, stringArrayStub)
	assert.NoError(t, err)
}

func TestCollectionAppendTagIssueNegativeErrorStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockCollection.
		EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errStub)

	collection := &Collection{collection: mockCollection}
	err := collection.AppendTagIssue(context.TODO(), usernameStub, stringArrayStub)
	assert.ErrorIs(t, err, errStub)
}
