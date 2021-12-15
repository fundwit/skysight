package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"skysight/infra/sessions"
	"skysight/repository"
	"skysight/testinfra"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fundwit/go-commons/types"
	. "github.com/onsi/gomega"
)

func TestRepositoryTableName(t *testing.T) {
	RegisterTestingT(t)

	t.Run("table name of Repository should correct", func(t *testing.T) {
		r := repository.RepositoryRecord{}
		Expect(r.TableName()).To(Equal("repositories"))
	})
}

func TestRepositoryTableCreation(t *testing.T) {
	RegisterTestingT(t)

	t.Run("table create sql should be correct", func(t *testing.T) {
		gormDB, mock := testinfra.SetUpMockSql()
		const sqlExpr = "CREATE TABLE `repositories` (" +
			"`id` BIGINT UNSIGNED NOT NULL," +
			"`uri` VARCHAR(512) NOT NULL," +
			"`create_time` DATETIME(6) NOT NULL," +
			"PRIMARY KEY (`id`)" +
			")"

		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		Expect(gormDB.AutoMigrate(&repository.RepositoryRecord{})).ShouldNot(HaveOccurred())

		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}

func TestCreateRepository(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to create repository", func(t *testing.T) {
		_, mock := testinfra.SetUpMockSql()

		const sqlExpr = "INSERT INTO `repositories` (`uri`,`create_time`,`id`) VALUES (?,?,?)"
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs("http://aaa", testinfra.AnyArgument{}, testinfra.AnyArgument{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs("bad-value", testinfra.AnyArgument{}, testinfra.AnyArgument{}).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		id, err := repository.CreateRepository(repository.Repository{Uri: "http://aaa"}, &sessions.Session{Context: context.TODO()})
		Expect(err).ToNot(HaveOccurred())
		Expect(id).ToNot(BeZero())

		id, err = repository.CreateRepository(repository.Repository{Uri: "bad-value"}, &sessions.Session{Context: context.TODO()})
		Expect(err).To(Equal(sql.ErrConnDone))
		Expect(id).To(BeZero())

		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}

func TestQueryRepositories(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to query repositories", func(t *testing.T) {
		_, mock := testinfra.SetUpMockSql()

		repo := repository.RepositoryRecord{ID: 100, CreateTime: types.CurrentTimestamp(),
			Repository: repository.Repository{Uri: "http://test"}}

		rows := sqlmock.NewRows([]string{"id", "uri", "create_time"}).
			AddRow(repo.ID, repo.Uri, repo.CreateTime)

		const sqlExpr = "SELECT * FROM `repositories`"
		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta(sqlExpr)).
			WillReturnError(sql.ErrConnDone)

		result, err := repository.QueryRepositories(repository.RepositoryQuery{}, &sessions.Session{Context: context.TODO()})
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal([]repository.RepositoryRecord{repo}))

		result, err = repository.QueryRepositories(repository.RepositoryQuery{}, &sessions.Session{Context: context.TODO()})
		Expect(err).To(Equal(sql.ErrConnDone))
		Expect(result).To(BeEmpty())

		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}

func TestDeleteRepository(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to delete repository", func(t *testing.T) {
		_, mock := testinfra.SetUpMockSql()

		const sqlExpr = "DELETE FROM `repositories` WHERE id = ?"
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs(200).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs(400).
			WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlExpr)).
			WithArgs(500).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repository.DeleteRepository(200, &sessions.Session{Context: context.TODO()})
		Expect(err).ToNot(HaveOccurred())

		err = repository.DeleteRepository(400, &sessions.Session{Context: context.TODO()})
		Expect(err).ToNot(HaveOccurred())

		err = repository.DeleteRepository(500, &sessions.Session{Context: context.TODO()})
		Expect(err).To(Equal(sql.ErrConnDone))

		Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
	})
}
