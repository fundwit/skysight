package sync

import (
	"skysight/infra/idgen"
	"skysight/infra/persistence"
	"skysight/infra/sessions"
	"skysight/repository"
	"time"

	"github.com/fundwit/go-commons/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	nowGenerator = time.Now
	idGenerator  = idgen.NextID

	SyncScheduleFunc = SyncSchedule
)

// SyncSchedule trigger all repositories to poll all changeset from remote.
// If the sync action of a repository is just triggered recently or the sync task is still undergoing,
// no new sync action should be started.
// Scan all repositories which last sync time already 3mins to now
func SyncSchedule(s *sessions.Session) error {
	repoUris := []string{}

	d := persistence.ActiveGormDB.WithContext(s.Context).
		Model(&repository.RepositoryRecord{}).
		Where("(last_sync_time IS NULL OR last_sync_time < ?)", types.Timestamp(nowGenerator().Add(-3*time.Minute))).Distinct()

	if err := d.Pluck("uri", &repoUris).Error; err != nil {
		return err
	}

	for _, repoUri := range repoUris {
		// min tx scope
		dbErr := persistence.ActiveGormDB.Transaction(func(tx *gorm.DB) error {
			// no active sync record (PENDING or RUNNING)
			r := []types.ID{}
			if err := tx.Model(&SyncRecord{}).
				Where("repo_uri = ? AND state IN ?", repoUri, []SyncState{SyncStatePending, SyncStateRunning}).
				Limit(1).Pluck("id", &r).Error; err != nil {
				return err
			}
			if len(r) != 0 {
				logrus.Warnf("pending or running sync run exist for repo: %s\n", repoUri)
				return nil
			}

			// create new sync record
			n := SyncRecord{
				ID:         idGenerator(idWorker),
				RepoUri:    repoUri,
				State:      SyncStatePending,
				CreateTime: types.Timestamp(nowGenerator()),
			}
			if err := tx.Create(&n).Error; err != nil {
				return err
			}
			return nil
		})

		if dbErr != nil {
			logrus.Warnf("error occurred on check and add new sync record: %v\n", dbErr)
			return dbErr
		}
	}

	return nil
}
