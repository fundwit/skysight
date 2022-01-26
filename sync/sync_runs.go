package sync

import (
	"skysight/infra/persistence"
	"skysight/infra/sessions"
)

var (
	QuerySyncRunsFunc = QuerySyncRuns
)

func QuerySyncRuns(q SyncRunQuery, s *sessions.Session) ([]SyncRecord, error) {
	repos := []SyncRecord{}
	if err := persistence.ActiveGormDB.WithContext(s.Context).Model(&SyncRecord{}).Scan(&repos).Error; err != nil {
		return repos, err
	}
	return repos, nil
}
