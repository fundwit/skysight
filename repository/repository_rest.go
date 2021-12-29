package repository

import (
	"net/http"
	"skysight/infra/fail"
	"skysight/infra/sessions"
	"skysight/misc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	PathRepositories = "/v1/repositories"
)

func RegisterRepositoriesRestAPI(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group(PathRepositories, middleWares...)
	g.POST("", handleCreateRepository)
	g.GET("", handleQueryRepositories)
	g.DELETE(":id", handleDeleteRepository)
}

// @ID repository-add
// @Param _ body repository.Repository true "request body"
// @Success 201 {object} misc.IdObject
// @Failure default {object} fail.ErrorBody "error"
// @Router /v1/repositories [post]
func handleCreateRepository(c *gin.Context) {
	creation := Repository{}
	err := c.ShouldBindBodyWith(&creation, binding.JSON)
	if err != nil {
		panic(&fail.ErrBadParam{Cause: err})
	}
	// session.ExtractSessionFromGinContext(c)
	id, err := CreateRepositoryFunc(creation, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusCreated, misc.NewIdObject(id))
}

// @ID repository-list
// @Param keyword query string false "query keyword"
// @Success 200 {array} repository.RepositoryRecord
// @Failure default {object} fail.ErrorBody "error"
// @Router /v1/repositories [get]
func handleQueryRepositories(c *gin.Context) {
	query := RepositoryQuery{}
	// err := c.MustBindWith(&query, binding.Query)
	// if err != nil {
	// 	panic(&fail.ErrBadParam{Cause: err})
	// }
	record, err := QueryRepositoriesFunc(query, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, record)
}

// @ID repository-delete
// @Param id path uint64 true "id of repository"
// @Success 204 {object} string "response body is empty"
// @Failure default {object} fail.ErrorBody "error"
// @Router /v1/repositories/{id} [delete]
func handleDeleteRepository(c *gin.Context) {
	id, err := misc.BindingPathID(c)
	if err != nil {
		panic(err)
	}

	err = DeleteRepositoryFunc(id, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.Status(http.StatusNoContent)
}
