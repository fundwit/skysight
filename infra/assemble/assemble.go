package assemble

import (
	"skysight/infra/doc"
	"skysight/infra/meta"
	"skysight/repository"
	"skysight/sync"

	"github.com/gin-gonic/gin"
)

/*
* registry endpoint for:
*
*   1. database auto migrations
*   2. rest api routes
*   3. error serialize
*   4. metric collectors
 */

type RestAPIRegister func(*gin.Engine, ...gin.HandlerFunc)

var AutoMigrations = []interface{}{}
var RestAPIRegistry = []RestAPIRegister{}

func init() {
	AutoMigrations = []interface{}{&repository.RepositoryRecord{}, &sync.SyncRecord{}}
	RestAPIRegistry = []RestAPIRegister{
		meta.RegisterMetaRestAPI,
		doc.RegisterDocsAPI,
		repository.RegisterRepositoriesRestAPI,
		sync.RegisterSyncRunsRestAPI,
	}
}
