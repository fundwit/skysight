package sync

import (
	"net/http"
	"skysight/infra/sessions"

	"github.com/gin-gonic/gin"
)

var (
	PathSyncRuns      = "/v1/runs"
	PathSyncSchedules = "/v1/schedules"
)

func RegisterSyncRunsRestAPI(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group(PathSyncRuns, middleWares...)
	g.GET("", handleQuerySyncRuns)

	s := r.Group(PathSyncSchedules, middleWares...)
	s.POST("", handleScheduleSyncRuns)
}

// @ID sync-runs-list
// @Param repoUri query string false "repository uri"
// @Success 200 {array} sync.SyncRecord
// @Failure default {object} fail.ErrorBody "error"
// @Router /v1/runs [get]
func handleQuerySyncRuns(c *gin.Context) {
	query := SyncRunQuery{}
	// err := c.MustBindWith(&query, binding.Query)
	// if err != nil {
	// 	panic(&fail.ErrBadParam{Cause: err})
	// }
	records, err := QuerySyncRunsFunc(query, &sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, records)
}

// @ID create-sync-schedule-request
// @Success 200 {string} string
// @Failure default {object} fail.ErrorBody "error"
// @Router /v1/schedules [post]
func handleScheduleSyncRuns(c *gin.Context) {
	err := SyncScheduleFunc(&sessions.Session{Context: c.Request.Context()})
	if err != nil {
		panic(err)
	}
	c.Status(http.StatusOK)
}
