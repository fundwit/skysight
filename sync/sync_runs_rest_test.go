package sync

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"skysight/infra/fail"
	"skysight/infra/sessions"
	"skysight/testinfra"
	"strings"
	"testing"
	"time"

	"github.com/fundwit/go-commons/types"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

func TestQuerySyncRunsAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(fail.ErrorHandling())
	RegisterSyncRunsRestAPI(router)

	// t.Run("should be able to validate parameters", func(t *testing.T) {
	// 	req := httptest.NewRequest(http.MethodGet, repository.PathRepositories, nil)
	// 	status, body, _ := testinfra.ExecuteRequest(req, router)
	// 	Expect(status).To(Equal(http.StatusBadRequest))
	// 	Expect(body).To(MatchJSON(`{"code":"common.bad_param",
	// 		"message":"Key: 'RepositoryQuery.ProjectID' Error:Field validation for 'ProjectID' failed on the 'required' tag",
	// 		"data":null}`))
	// })

	t.Run("should be able to handle error", func(t *testing.T) {
		QuerySyncRunsFunc = func(q SyncRunQuery, s *sessions.Session) ([]SyncRecord, error) {
			return nil, errors.New("some error")
		}
		req := httptest.NewRequest(http.MethodGet, PathSyncRuns, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "message":"some error", "data":null}`))
	})

	t.Run("should be able to handle query request successfully", func(t *testing.T) {
		demoTime := types.TimestampOfDate(2020, 1, 1, 1, 0, 0, 0, time.Now().Location())
		timeBytes, err := demoTime.Time().MarshalJSON()
		Expect(err).To(BeNil())
		timeString := strings.Trim(string(timeBytes), `"`)

		var q1 SyncRunQuery
		QuerySyncRunsFunc = func(q SyncRunQuery, s *sessions.Session) ([]SyncRecord, error) {
			q1 = q
			return []SyncRecord{{ID: 123, RepoUri: "http://example.com/foo.git", State: SyncStatePending, CreateTime: demoTime,
				BeginTime: demoTime, EndTime: demoTime, RootCause: "some error"}}, nil
		}
		req := httptest.NewRequest(http.MethodGet, PathSyncRuns, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(MatchJSON(`[{"id": "123", "repoUri": "http://example.com/foo.git", "createTime": "` + timeString + `",` +
			`"state": 1, "beginTime": "` + timeString + `", "endTime": "` + timeString + `", "rootCause": "some error"` +
			`}]`))
		Expect(q1).To(Equal(SyncRunQuery{}))
	})
}

func TestCreateScheduleRequestAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(fail.ErrorHandling())
	RegisterSyncRunsRestAPI(router)

	t.Run("should be able to handle error", func(t *testing.T) {
		SyncScheduleFunc = func(s *sessions.Session) error {
			return errors.New("some error")
		}
		req := httptest.NewRequest(http.MethodPost, PathSyncSchedules, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "message":"some error", "data":null}`))
	})

	t.Run("should be able to handle create schedule request successfully", func(t *testing.T) {
		SyncScheduleFunc = func(s *sessions.Session) error {
			return nil
		}
		req := httptest.NewRequest(http.MethodPost, PathSyncSchedules, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(BeEmpty())
	})
}
