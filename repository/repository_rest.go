package repository

import (
	"net/http"
	"skysight/bizerror"
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

func handleCreateRepository(c *gin.Context) {
	creation := Repository{}
	err := c.ShouldBindBodyWith(&creation, binding.JSON)
	if err != nil {
		panic(&bizerror.ErrBadParam{Cause: err})
	}
	// session.ExtractSessionFromGinContext(c)
	id, err := CreateRepositoryFunc(creation, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusCreated, misc.IdObject(id))
}

func handleQueryRepositories(c *gin.Context) {
	query := RepositoryQuery{}
	// err := c.MustBindWith(&query, binding.Query)
	// if err != nil {
	// 	panic(&bizerror.ErrBadParam{Cause: err})
	// }
	record, err := QueryRepositoriesFunc(query, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, record)
}

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
