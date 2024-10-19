package repository

import (
	"gorm.io/gorm"
	"mnc-test/core/entity"
	"mnc-test/data_source/postgres_datasource"
)

// Repository struct is used to store every data source connections.
type Repository struct {
	PostgresDatasource *postgres_datasource.PostgresDatasource
}

// NewRepository is to initiate new repository.
func NewRepository(postgresDatasource *postgres_datasource.PostgresDatasource) *Repository {
	return &Repository{
		PostgresDatasource: postgresDatasource,
	}
}

func (r *Repository) AutoMigration(tx *gorm.DB) (err error) {
	db := r.GetDBPostgres(nil)
	err = db.AutoMigrate(
		&entity.User{},
		&entity.UserBalance{},
		&entity.Transaction{},
	)
	return
}

func (r *Repository) GetDBPostgres(tx *gorm.DB) (db *gorm.DB) {
	return r.PostgresDatasource.GetDB(tx)
}
