package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetFinesByID(id int) (*ds.Fines, error) {
	Fines := &ds.Fines{}

	err := r.db.First(Fines, "id = ?", "1").Error // find Fines with id = 1
	if err != nil {
		return nil, err
	}

	return Fines, nil
}

func (r *Repository) CreateFines(Fines ds.Fines) error {
	return r.db.Create(Fines).Error
}
