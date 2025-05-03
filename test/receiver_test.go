package test

import (
	"context"
	"testing"

	domain "github.com/IndianMax03/beroli-bot/internal/domain"
	mocks "github.com/IndianMax03/beroli-bot/mocks"
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	issueCreateRequestStub = &model.IssueCreateRequest{}
	stubCommand            = &stubCommandStruct{}
	commandMapStub         = map[string]domain.Command{
		domain.MY_ISSUES_COMMAND:    stubCommand,
		domain.DONE_COMMAND:         stubCommand,
		domain.CANCEL_COMMAND:       stubCommand,
		domain.NIL_COMMAND:          stubCommand,
		domain.CREATE_ISSUE_COMMAND: stubCommand,
		domain.STATE_COMMAND:        stubCommand,
		domain.HELP_COMMAND:         stubCommand,
	}
	emptyCommandMapStub = map[string]domain.Command{}
	noCommandMessage    = "Ни одной команды не зарегестрировано"
	nonEmptyIssuesStub  = []domain.Issue{
		*issueStub,
		*issueStub,
		*issueStub,
	}
	emptyIssuesStub = []domain.Issue{}
	fileIDStub      = tgbotapi.FileID("stub")
	emptyStringStub = ""
)

type stubCommandStruct struct{}

func (*stubCommandStruct) Execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return "stubResult", nil
}

func (*stubCommandStruct) GetName() string {
	return "stubName"
}

func (*stubCommandStruct) GetDescription() string {
	return "stubDescription"
}

func TestReceiverMyIssuesPositiveNonEmptyIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetIssues(gomock.Any(), usernameStub).
		Return(nonEmptyIssuesStub, nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.MyIssues(usernameStub)
	assert.NoError(t, err)
}

func TestReceiverMyIssuesPositiveEmptyIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetIssues(gomock.Any(), usernameStub).
		Return(emptyIssuesStub, nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.MyIssues(usernameStub)
	assert.NoError(t, err)
}

func TestReceiverMyIssuesNegativeGetIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetIssues(gomock.Any(), usernameStub).
		Return(nil, errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.MyIssues(usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverCreateIssuePositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CREATING_STATE).
		Return(nil)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(nil)

	mockCollectionService.
		EXPECT().
		GetIssue(gomock.Any(), usernameStub).
		Return(issueCreateRequestStub, nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.CreateIssue(usernameStub, textStub, commandMapStub)
	assert.NoError(t, err)
}

func TestReceiverCreateIssueNegativeUpdateStateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	// mockYandexTrackerService := mocks.NewMockYandexTrackerService(ctrl)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CREATING_STATE).
		Return(errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.CreateIssue(usernameStub, textStub, commandMapStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverCreateIssueNegativeClearIssue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	// mockYandexTrackerService := mocks.NewMockYandexTrackerService(ctrl)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CREATING_STATE).
		Return(nil)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.CreateIssue(usernameStub, textStub, commandMapStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverCreateIssueNegativeGetIssue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	// mockYandexTrackerService := mocks.NewMockYandexTrackerService(ctrl)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CREATING_STATE).
		Return(nil)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(nil)

	mockCollectionService.
		EXPECT().
		GetIssue(gomock.Any(), usernameStub).
		Return(nil, errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.CreateIssue(usernameStub, textStub, commandMapStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverCancelPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(nil)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CANCELED_STATE).
		Return(nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.Cancel(usernameStub)
	assert.NoError(t, err)
}

func TestReceiverCancelNegativeUpdateStateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(nil)

	mockCollectionService.
		EXPECT().
		UpdateStateUser(gomock.Any(), usernameStub, domain.CANCELED_STATE).
		Return(errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.Cancel(usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverCancelNegativeClearIssue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		ClearIssue(gomock.Any(), usernameStub).
		Return(errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}

	_, err := receiver.Cancel(usernameStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverNoCommandPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	for tag := range domain.LocalizedTagsDescriptionMap {
		switch tag {
		case domain.ISSUE_SUMMARY_TAG:
			mockCollectionService.
				EXPECT().
				UpdateSummaryIssue(gomock.Any(), usernameStub, textStub).
				Return(nil)
		case domain.ISSUE_DESCRIPTION_TAG:
			mockCollectionService.
				EXPECT().
				UpdateDescriptionIssue(gomock.Any(), usernameStub, textStub).
				Return(nil)
		case domain.ISSUE_ATTACHMENT_TAG:
			mockCollectionService.
				EXPECT().
				AppendAttachmentIssue(gomock.Any(), usernameStub, string(fileIDStub)).
				Return(nil)
		case domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG:
			mockCollectionService.
				EXPECT().
				AppendDescriptionAttachmentIssue(gomock.Any(), usernameStub, string(fileIDStub)).
				Return(nil)
		case domain.ISSUE_TAGS_TAG:
			mockCollectionService.
				EXPECT().
				AppendTagIssue(gomock.Any(), usernameStub, []string{textStub}).
				Return(nil)
		}
		receiver := &domain.Handler{
			Collection: mockCollectionService,
		}
		_, err := receiver.NoCommand(usernameStub, textStub, tag, fileIDStub)
		assert.NoError(t, err)
	}
}

func TestReceiverNoCommandNegativeEmptyTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	_, err := receiver.NoCommand(usernameStub, textStub, emptyStringStub, fileIDStub)
	assert.ErrorIs(t, err, domain.ErrEmptyTag)
}

func TestReceiverHelpCommandPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	result, err := receiver.HelpCommand(commandMapStub)
	assert.NoError(t, err)

	assert.NotContains(t, result, noCommandMessage)
}

func TestReceiverHelpCommandNegativeNoCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	result, err := receiver.HelpCommand(emptyCommandMapStub)
	assert.NoError(t, err)

	assert.Contains(t, result, noCommandMessage)
}

func TestReceiverStateCommandPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetStateUser(gomock.Any(), usernameStub).
		Return(domain.NIL_STATE, nil).
		MinTimes(1)

	mockCollectionService.
		EXPECT().
		GetIssues(gomock.Any(), usernameStub).
		Return(emptyIssuesStub, nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	_, err := receiver.StateCommand(usernameStub, commandMapStub)
	assert.NoError(t, err)
}

func TestReceiverStateCommandNegativeGetIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetStateUser(gomock.Any(), usernameStub).
		Return(domain.NIL_STATE, nil).
		MinTimes(1)

	mockCollectionService.
		EXPECT().
		GetIssues(gomock.Any(), usernameStub).
		Return(nil, errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	_, err := receiver.StateCommand(usernameStub, commandMapStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverStateCommandNegativeGetStateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		GetStateUser(gomock.Any(), usernameStub).
		Return(domain.NIL_STATE, errStub)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	_, err := receiver.StateCommand(usernameStub, commandMapStub)
	assert.ErrorIs(t, err, errStub)
}

func TestReceiverValidateAndInitUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		ExistsUser(gomock.Any(), usernameStub).
		Return(nil)

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	err := receiver.ValidateAndInitUser(usernameStub)
	assert.NoError(t, err)
}

func TestReceiverValidateAndInitUserNegativeExistsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	mockCollectionService.
		EXPECT().
		ExistsUser(gomock.Any(), usernameStub).
		Return(errStub)

	mockCollectionService.
		EXPECT().
		CreateUser(gomock.Any(), gomock.Any())

	receiver := &domain.Handler{
		Collection: mockCollectionService,
	}
	err := receiver.ValidateAndInitUser(usernameStub)
	assert.NoError(t, err)
}

func TestReceiverValidateStatePositiveFromNilDoneCanceledState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	states := []string{domain.NIL_STATE, domain.DONE_STATE, domain.CANCELED_STATE}
	commands := []string{domain.MY_ISSUES_COMMAND, domain.CREATE_ISSUE_COMMAND, domain.HELP_COMMAND, domain.STATE_COMMAND}
	for _, st := range states {
		for _, cmd := range commands {
			mockCollectionService.
				EXPECT().
				GetStateUser(gomock.Any(), usernameStub).
				Return(st, nil)
			receiver := &domain.Handler{
				Collection: mockCollectionService,
			}
			err := receiver.ValidateState(usernameStub, cmd)
			assert.NoError(t, err)
		}
	}
}

func TestReceiverValidateStatePositiveFromCreatingState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	state := domain.CREATING_STATE
	commands := []string{domain.MY_ISSUES_COMMAND, domain.DONE_COMMAND, domain.CANCEL_COMMAND, domain.NIL_COMMAND, domain.HELP_COMMAND, domain.STATE_COMMAND}
	for _, cmd := range commands {
		mockCollectionService.
			EXPECT().
			GetStateUser(gomock.Any(), usernameStub).
			Return(state, nil)
		receiver := &domain.Handler{
			Collection: mockCollectionService,
		}
		err := receiver.ValidateState(usernameStub, cmd)
		assert.NoError(t, err)
	}
}

func TestReceiverValidateStateNegativeFromNilDoneCanceledStateErrCreatingNotStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	states := []string{domain.NIL_STATE, domain.DONE_STATE, domain.CANCELED_STATE}
	commands := []string{domain.DONE_COMMAND, domain.CANCEL_COMMAND, domain.NIL_COMMAND}
	for _, st := range states {
		for _, cmd := range commands {
			mockCollectionService.
				EXPECT().
				GetStateUser(gomock.Any(), usernameStub).
				Return(st, nil)
			receiver := &domain.Handler{
				Collection: mockCollectionService,
			}
			err := receiver.ValidateState(usernameStub, cmd)
			assert.ErrorIs(t, err, domain.ErrCreatingNotStarted)
		}
	}
}

func TestReceiverValidateStateNegativeFromCreatingStateErrCreatingStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)

	state := domain.CREATING_STATE
	commands := []string{domain.CREATE_ISSUE_COMMAND}
	for _, cmd := range commands {
		mockCollectionService.
			EXPECT().
			GetStateUser(gomock.Any(), usernameStub).
			Return(state, nil)
		receiver := &domain.Handler{
			Collection: mockCollectionService,
		}
		err := receiver.ValidateState(usernameStub, cmd)
		assert.ErrorIs(t, err, domain.ErrCreatingStarted)
	}
}
