package sync

import (
	"context"
	"database/sql"
	"regexp"
	"skysight/infra/sessions"
	"skysight/testinfra"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fundwit/go-commons/types"
	. "github.com/onsi/gomega"
)

func TestQuerySyncRuns(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to query sync runs", func(t *testing.T) {
		_, mock := testinfra.SetUpMockSql()

		r := SyncRecord{
			ID:         100,
			RepoUri:    "https://exmaple/foo.git",
			State:      SyncStatePending,
			CreateTime: types.CurrentTimestamp(),

			BeginTime: types.CurrentTimestamp(),
			EndTime:   types.CurrentTimestamp(),
			RootCause: "some error",
		}

		rows := sqlmock.NewRows([]string{"id", "repo_uri", "state", "create_time", "begin_time", "end_time", "root_cause"}).
			AddRow(r.ID, r.RepoUri, r.State, r.CreateTime, r.BeginTime, r.EndTime, r.RootCause)

		const sqlExpr = "SELECT * FROM `repo_syncs`"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WillReturnError(sql.ErrConnDone)

		result, err := QuerySyncRuns(SyncRunQuery{}, &sessions.Session{Context: context.TODO()})
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal([]SyncRecord{r}))

		result, err = QuerySyncRuns(SyncRunQuery{}, &sessions.Session{Context: context.TODO()})
		Expect(err).To(Equal(sql.ErrConnDone))
		Expect(result).To(BeEmpty())

		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}
