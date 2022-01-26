package sync

import (
	"context"
	"errors"
	"regexp"
	"skysight/infra/sessions"
	"skysight/testinfra"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fundwit/go-commons/types"
	. "github.com/onsi/gomega"
	"github.com/sony/sonyflake"
)

func TestSyncSchedule(t *testing.T) {
	RegisterTestingT(t)
	mockErr := errors.New("mock error")

	t.Run("should be able to trigger sync", func(t *testing.T) {
		nowGenerator = func() time.Time {
			return time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
		}
		nextId := 200
		idGenerator = func(idWorker *sonyflake.Sonyflake) types.ID {
			id := nextId
			nextId = nextId + 1
			return types.ID(id)
		}

		_, mock := testinfra.SetUpMockSql()
		sqlExpr := "SELECT DISTINCT `uri` FROM `repositories` WHERE (last_sync_time IS NULL OR last_sync_time < ?)"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs(types.Timestamp(nowGenerator().Add(-3 * time.Minute))).
			WillReturnRows(sqlmock.NewRows([]string{"uri"}).AddRow("https://example/foo.git"))
		mock.ExpectBegin()
		sqlExpr = "SELECT `id` FROM `repo_syncs` WHERE repo_uri = ? AND state IN (?,?) LIMIT 1"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, SyncStateRunning).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		sqlExpr = "INSERT INTO `repo_syncs` (`repo_uri`,`state`,`create_time`,`begin_time`,`end_time`,`root_cause`,`id`) VALUES (?,?,?,?,?,?,?)"
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, types.Timestamp(nowGenerator()), types.Timestamp{}, types.Timestamp{}, "", types.ID(200)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		Expect(SyncSchedule(&sessions.Session{Context: context.TODO()})).ShouldNot(HaveOccurred())
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})

	t.Run("should skip trigger sync when pending or running sync run exist", func(t *testing.T) {
		nowGenerator = func() time.Time {
			return time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
		}

		_, mock := testinfra.SetUpMockSql()
		sqlExpr := "SELECT DISTINCT `uri` FROM `repositories` WHERE (last_sync_time IS NULL OR last_sync_time < ?)"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs(types.Timestamp(nowGenerator().Add(-3 * time.Minute))).
			WillReturnRows(sqlmock.NewRows([]string{"uri"}).AddRow("https://example/foo.git"))
		mock.ExpectBegin()
		sqlExpr = "SELECT `id` FROM `repo_syncs` WHERE repo_uri = ? AND state IN (?,?) LIMIT 1"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, SyncStateRunning).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
		mock.ExpectCommit()

		Expect(SyncSchedule(&sessions.Session{Context: context.TODO()})).ShouldNot(HaveOccurred())
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})

	t.Run("should return error when create sync record failed", func(t *testing.T) {
		nowGenerator = func() time.Time {
			return time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
		}
		nextId := 200
		idGenerator = func(idWorker *sonyflake.Sonyflake) types.ID {
			id := nextId
			nextId = nextId + 1
			return types.ID(id)
		}

		_, mock := testinfra.SetUpMockSql()
		sqlExpr := "SELECT DISTINCT `uri` FROM `repositories` WHERE (last_sync_time IS NULL OR last_sync_time < ?)"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs(types.Timestamp(nowGenerator().Add(-3 * time.Minute))).
			WillReturnRows(sqlmock.NewRows([]string{"uri"}).AddRow("https://example/foo.git"))
		mock.ExpectBegin()
		sqlExpr = "SELECT `id` FROM `repo_syncs` WHERE repo_uri = ? AND state IN (?,?) LIMIT 1"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, SyncStateRunning).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		sqlExpr = "INSERT INTO `repo_syncs` (`repo_uri`,`state`,`create_time`,`begin_time`,`end_time`,`root_cause`,`id`) VALUES (?,?,?,?,?,?,?)"
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, types.Timestamp(nowGenerator()), types.Timestamp{}, types.Timestamp{}, "", types.ID(200)).
			WillReturnError(mockErr)
		mock.ExpectRollback()

		Expect(SyncSchedule(&sessions.Session{Context: context.TODO()})).To(Equal(mockErr))
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})

	t.Run("should return error when query running or pending sync records failed", func(t *testing.T) {
		nowGenerator = func() time.Time {
			return time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
		}
		_, mock := testinfra.SetUpMockSql()
		sqlExpr := "SELECT DISTINCT `uri` FROM `repositories` WHERE (last_sync_time IS NULL OR last_sync_time < ?)"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs(types.Timestamp(nowGenerator().Add(-3 * time.Minute))).
			WillReturnRows(sqlmock.NewRows([]string{"uri"}).AddRow("https://example/foo.git"))
		mock.ExpectBegin()
		sqlExpr = "SELECT `id` FROM `repo_syncs` WHERE repo_uri = ? AND state IN (?,?) LIMIT 1"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs("https://example/foo.git", SyncStatePending, SyncStateRunning).
			WillReturnError(mockErr)
		mock.ExpectRollback()

		Expect(SyncSchedule(&sessions.Session{Context: context.TODO()})).To(Equal(mockErr))
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})

	t.Run("should return error when query repositories failed", func(t *testing.T) {
		nowGenerator = func() time.Time {
			return time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
		}

		_, mock := testinfra.SetUpMockSql()
		sqlExpr := "SELECT DISTINCT `uri` FROM `repositories` WHERE (last_sync_time IS NULL OR last_sync_time < ?)"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WithArgs(types.Timestamp(nowGenerator().Add(-3 * time.Minute))).
			WillReturnError(mockErr)

		Expect(SyncSchedule(&sessions.Session{Context: context.TODO()})).To(Equal(mockErr))
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}
