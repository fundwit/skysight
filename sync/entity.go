package sync

import (
	"github.com/fundwit/go-commons/types"
	"github.com/sony/sonyflake"
)

type SyncState int

const (
	SyncStateUnknown SyncState = 0
	SyncStatePending SyncState = 1
	SyncStateRunning SyncState = 2
	SyncStateSuccess SyncState = 6
	SyncStateFail    SyncState = 7
)

var (
	idWorker = sonyflake.NewSonyflake(sonyflake.Settings{})
)

type SyncRecord struct {
	ID         types.ID        `json:"id" gorm:"primary_key;type:BIGINT UNSIGNED NOT NULL"`
	RepoUri    string          `json:"repoUri" gorm:"size:300;type:VARCHAR(300) NOT NULL"`
	State      SyncState       `json:"state" gorm:"size:7;type:TINYINT UNSIGNED NOT NULL"`
	CreateTime types.Timestamp `json:"createTime" gorm:"type:DATETIME(6) NOT NULL"`

	BeginTime types.Timestamp `json:"beginTime" gorm:"type:DATETIME(6) NULL"`
	EndTime   types.Timestamp `json:"endTime" gorm:"type:DATETIME(6) NULL"`
	RootCause string          `json:"rootCause" gorm:"size:1000;type:VARCHAR(1000) NULL DEFAULT ''"`
}

func (r *SyncRecord) TableName() string {
	return "repo_syncs"
}

type SyncRunQuery struct {
	// RepoUri string `form:"repoUri" binding:"omitempty,lte=300"`
}
