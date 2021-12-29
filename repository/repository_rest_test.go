package repository_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"skysight/infra/fail"
	"skysight/infra/sessions"
	"skysight/repository"
	"skysight/testinfra"
	"strings"
	"testing"
	"time"

	"github.com/fundwit/go-commons/types"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

func TestQueryRepositoriesAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(fail.ErrorHandling())
	repository.RegisterRepositoriesRestAPI(router)

	// t.Run("should be able to validate parameters", func(t *testing.T) {
	// 	req := httptest.NewRequest(http.MethodGet, repository.PathRepositories, nil)
	// 	status, body, _ := testinfra.ExecuteRequest(req, router)
	// 	Expect(status).To(Equal(http.StatusBadRequest))
	// 	Expect(body).To(MatchJSON(`{"code":"common.bad_param",
	// 		"message":"Key: 'RepositoryQuery.ProjectID' Error:Field validation for 'ProjectID' failed on the 'required' tag",
	// 		"data":null}`))
	// })

	t.Run("should be able to handle error", func(t *testing.T) {
		repository.QueryRepositoriesFunc = func(q repository.RepositoryQuery, s *sessions.Session) ([]repository.RepositoryRecord, error) {
			return nil, errors.New("some error")
		}
		req := httptest.NewRequest(http.MethodGet, repository.PathRepositories, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "message":"some error", "data":null}`))
	})

	t.Run("should be able to handle query request successfully", func(t *testing.T) {
		demoTime := types.TimestampOfDate(2020, 1, 1, 1, 0, 0, 0, time.Now().Location())
		timeBytes, err := demoTime.Time().MarshalJSON()
		Expect(err).To(BeNil())
		timeString := strings.Trim(string(timeBytes), `"`)

		var q1 repository.RepositoryQuery
		repository.QueryRepositoriesFunc = func(q repository.RepositoryQuery, s *sessions.Session) ([]repository.RepositoryRecord, error) {
			q1 = q
			return []repository.RepositoryRecord{{ID: 123, CreateTime: demoTime, Repository: repository.Repository{Uri: "http://xxxx"}}}, nil
		}
		req := httptest.NewRequest(http.MethodGet, repository.PathRepositories, nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(MatchJSON(`[{"id": "123", "uri": "http://xxxx", "createTime": "` + timeString + `"}]`))
		Expect(q1).To(Equal(repository.RepositoryQuery{}))
	})
}

func TestCreateRepositoryAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(fail.ErrorHandling())
	repository.RegisterRepositoriesRestAPI(router)

	t.Run("should be able to validate parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, repository.PathRepositories, strings.NewReader("{}"))
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(body).To(MatchJSON(`{"code":"common.bad_param", "data":null,
			"message": "Key: 'Repository.Uri' Error:Field validation for 'Uri' failed on the 'required' tag"}`))

		req = httptest.NewRequest(http.MethodPost, repository.PathRepositories, nil)
		status, body, _ = testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(body).To(MatchJSON(`{"code": "common.bad_param", "message": "EOF", "data": null}`))

		req = httptest.NewRequest(http.MethodPost, repository.PathRepositories, strings.NewReader(" \t "))
		status, body, _ = testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(body).To(MatchJSON(`{"code": "common.bad_param", "message": "EOF", "data": null}`))

		req = httptest.NewRequest(http.MethodPost, repository.PathRepositories, strings.NewReader(" xx "))
		status, body, _ = testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(body).To(MatchJSON(`{"code": "common.bad_param", "message": "invalid character 'x' looking for beginning of value", "data": null}`))
	})

	t.Run("should be able to handle error", func(t *testing.T) {
		repository.CreateRepositoryFunc = func(r repository.Repository, s *sessions.Session) (types.ID, error) {
			return 100, errors.New("some error")
		}
		reqBody := `{"uri": "http://some-repo"}`
		req := httptest.NewRequest(http.MethodPost, repository.PathRepositories, strings.NewReader(reqBody))
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "message":"some error", "data":null}`))
	})

	t.Run("should be able to create repository successfully", func(t *testing.T) {
		repository.CreateRepositoryFunc = func(r repository.Repository, s *sessions.Session) (types.ID, error) {
			return 100, nil
		}
		reqBody := `{"uri": "http://some-repo"}`
		req := httptest.NewRequest(http.MethodPost, repository.PathRepositories, strings.NewReader(reqBody))
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusCreated))
		Expect(body).To(MatchJSON(`{"id": "100"}`))
	})
}

func TestDeleteRepositoryAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(fail.ErrorHandling())
	repository.RegisterRepositoriesRestAPI(router)

	t.Run("should be able to validate parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, repository.PathRepositories+"/aaa", nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(body).To(MatchJSON(`{"code":"common.bad_param",
		"message": "invalid id 'aaa'",
		"data":null}`))
	})

	t.Run("should be able to delete repository", func(t *testing.T) {
		var reqId types.ID
		repository.DeleteRepositoryFunc = func(id types.ID, s *sessions.Session) error {
			reqId = id
			return nil
		}
		req := httptest.NewRequest(http.MethodDelete, repository.PathRepositories+"/100", nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusNoContent))
		Expect(body).To(BeZero())

		Expect(reqId).To(Equal(types.ID(100)))
	})

	t.Run("should be able to handle error", func(t *testing.T) {
		repository.DeleteRepositoryFunc = func(id types.ID, s *sessions.Session) error {
			return errors.New("some error")
		}
		req := httptest.NewRequest(http.MethodDelete, repository.PathRepositories+"/100", nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "message":"some error", "data":null}`))
	})
}
