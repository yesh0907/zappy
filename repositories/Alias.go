package repositories

import (
	"gorm.io/gorm"
	"zappy.sh/models"
)

type AliasRepository struct {
	db *gorm.DB
}

func NewAliasRepository(db *gorm.DB) *AliasRepository {
	return &AliasRepository{db}
}

func (r *AliasRepository) CreateAlias(alias *models.Alias) error {
	return r.db.Create(&alias).Error
}

func (r *AliasRepository) GetAlias(name string) (models.Alias, error) {
	var alias models.Alias
	tx := r.db.First(&alias, "name = ?", name)
	return alias, tx.Error
}
