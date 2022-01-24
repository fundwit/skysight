package fail

import (
	"errors"
	"net/http"
)

var ErrUnexpected = errors.New("common.internal_server_error")
var ErrInvalidArguments = errors.New("invalid arguments")
var ErrUnauthenticated = errors.New("unauthenticated")
var ErrForbidden = errors.New("forbidden")

var ErrWorkflowIsReferenced = errors.New("workflow is referenced")
var ErrUnknownState = errors.New("unknown workflow state")
var ErrStateExisted = errors.New("state existed")
var ErrStateInvalid = errors.New("state is invalid")
var ErrStateCategoryInvalid = errors.New("state category is invalid")
var ErrArchiveStatusInvalid = errors.New("archive status is invalid")
var ErrWorkProcessStepStateInvalid = errors.New("state of work process step is invalid")

var ErrLabelNotFound = errors.New("label not found")
var ErrLabelIsReferenced = errors.New("label is referenced")

var ErrPropertyDefinitionInvalid = errors.New("invalid property definition")
var ErrPropertyDefinitionNotFound = errors.New("property definition not found")
var ErrPropertyDefinitionIsReferenced = errors.New("property definition is referenced")

var ErrNotFound = errors.New("not found")
var ErrNoContent = errors.New("no content")
var ErrInvalidPassword = errors.New("invalid password")

var ErrLastProjectManagerDelete = errors.New("last project manager delete")
var ErrProjectMemberSelfGrant = errors.New("project member self grant")

type BizError interface {
	Respond() *BizErrorDetail
}

type BizErrorDetail struct {
	Status  int
	Code    string
	Message string

	Data interface{}
}

type ErrBadParam struct {
	Param        string
	InvalidValue string

	Cause error
}

func (e *ErrBadParam) Unwrap() error {
	return e.Cause
}

func (e *ErrBadParam) Error() string {
	message := "bad param"
	if e.Param != "" {
		message = "invalid " + e.Param + " '" + e.InvalidValue + "'"
	} else if e.Cause != nil {
		message = e.Cause.Error()
	}
	return message
}

func (e *ErrBadParam) Respond() *BizErrorDetail {
	return &BizErrorDetail{
		Status:  http.StatusBadRequest,
		Code:    "common.bad_param",
		Message: e.Error(),
		Data:    nil,
	}
}
