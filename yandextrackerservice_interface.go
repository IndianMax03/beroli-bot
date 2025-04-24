package main

import (
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"resty.dev/v3"
)

type YandexTrackerService interface {
	CreateIssue(req *model.IssueCreateRequest) (*model.IssueResponse, error)
	UploadTemporaryAttachment(multipartReq *resty.MultipartField) (*model.AttachmentFileResponse, error)
}
