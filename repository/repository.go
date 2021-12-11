package repository

import (
	"skysight/infra/idgen"
	"skysight/infra/persistence"
	"skysight/infra/sessions"

	"github.com/fundwit/go-commons/types"
	"github.com/sony/sonyflake"
)

var (
	CreateRepositoryFunc  = CreateRepository
	QueryRepositoriesFunc = QueryRepositories
	DeleteRepositoryFunc  = DeleteRepository

	idWorker = sonyflake.NewSonyflake(sonyflake.Settings{})
)

type Repository struct {
	Uri string `json:"uri" binding:"required,lte=250"`
}

type RepositoryRecord struct {
	ID types.ID `json:"id" gorm:"primary_key"`

	Repository

	CreateTime types.Timestamp `json:"createTime"`
}

type RepositoryQuery struct {
}

func (r *RepositoryRecord) TableName() string {
	return "repositories"
}

func CreateRepository(r Repository, s *sessions.Session) (types.ID, error) {
	repo := RepositoryRecord{
		ID:         idgen.NextID(idWorker),
		Repository: r,
		CreateTime: types.CurrentTimestamp(),
	}

	if err := persistence.ActiveGormDB.WithContext(s.Context).Create(&repo).Error; err != nil {
		return 0, err
	}

	return repo.ID, nil
}

func QueryRepositories(q RepositoryQuery, s *sessions.Session) ([]RepositoryRecord, error) {
	repos := []RepositoryRecord{}
	if err := persistence.ActiveGormDB.WithContext(s.Context).Model(&RepositoryRecord{}).Scan(&repos).Error; err != nil {
		return repos, err
	}
	return repos, nil
}

func DeleteRepository(id types.ID, s *sessions.Session) error {
	if err := persistence.ActiveGormDB.WithContext(s.Context).Delete(&RepositoryRecord{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
