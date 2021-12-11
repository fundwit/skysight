package bizerror_test

import (
	"net/http"
	"skysight/bizerror"
	"testing"

	. "github.com/onsi/gomega"
)

func TestErrBadParam(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return detailed message if param is nil", func(t *testing.T) {
		err := bizerror.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: bizerror.ErrForbidden}
		Expect(err.Error()).To(Equal("invalid id 'aaa'"))
	})

	t.Run("should invoke the Error() function of cause property if cause is not nil", func(t *testing.T) {
		err := bizerror.ErrBadParam{Cause: bizerror.ErrForbidden}
		Expect(err.Error()).To(Equal("forbidden"))
	})

	t.Run("should return default message if cause is nil", func(t *testing.T) {
		err := bizerror.ErrBadParam{}
		Expect(err.Error()).To(Equal("bad param"))
	})
}

func TestErrBadParam_Respond(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return response data as expected", func(t *testing.T) {
		err := bizerror.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: bizerror.ErrForbidden}
		Expect(*err.Respond()).To(Equal(bizerror.BizErrorDetail{
			Status:  http.StatusBadRequest,
			Code:    "common.bad_param",
			Message: "invalid id 'aaa'",
			Data:    nil,
		}))
	})
}

func TestErrBadParam_Unwrap(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return response data as expected", func(t *testing.T) {
		err := bizerror.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: bizerror.ErrForbidden}
		Expect(err.Unwrap()).To(Equal(bizerror.ErrForbidden))

		err = bizerror.ErrBadParam{}
		Expect(err.Unwrap()).To(BeNil())
	})
}
