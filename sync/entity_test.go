package sync

import (
	"regexp"
	"skysight/testinfra"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/gomega"
)

func TestSyncRecordTableName(t *testing.T) {
	RegisterTestingT(t)

	t.Run("table name of Repository should correct", func(t *testing.T) {
		r := SyncRecord{}
		Expect(r.TableName()).To(Equal("repo_syncs"))
	})
}

func TestSyncRecordTableCreation(t *testing.T) {
	RegisterTestingT(t)

	t.Run("table create sql should be correct", func(t *testing.T) {
		gormDB, mock := testinfra.SetUpMockSql()
		const sqlExpr = "CREATE TABLE `repo_syncs` (" +
			"`id` BIGINT UNSIGNED NOT NULL," +
			"`repo_uri` VARCHAR(300) NOT NULL," +
			"`state` TINYINT UNSIGNED NOT NULL," +
			"`create_time` DATETIME(6) NOT NULL," +
			"`begin_time` DATETIME(6) NULL," +
			"`end_time` DATETIME(6) NULL," +
			"`root_cause` VARCHAR(1000) NULL DEFAULT ''," +
			"PRIMARY KEY (`id`)" +
			")"

		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		Expect(gormDB.AutoMigrate(&SyncRecord{})).ShouldNot(HaveOccurred())
		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}
