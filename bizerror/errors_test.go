package bizerror_test

import (
	"skysight/bizerror"
	"testing"

	. "github.com/onsi/gomega"
)

func TestErrBadParam(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return default message if cause is nil", func(t *testing.T) {
		err := bizerror.ErrBadParam{}
		Expect(err.Error()).To(Equal("common.bad_param"))
	})

	t.Run("should invoke the Error() function of cause property if cause is not nil", func(t *testing.T) {
		err := bizerror.ErrBadParam{Cause: bizerror.ErrForbidden}
		Expect(err.Error()).To(Equal("forbidden"))
	})
}
