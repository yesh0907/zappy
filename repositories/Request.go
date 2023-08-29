package repositories

import (
	"gorm.io/gorm"
	"zappy.sh/models"
)

type RequestRepository struct {
	db *gorm.DB
}

type RequestsWithCount struct {
	Requests []models.Request `json:"requests"`
	Count    int64            `json:"count"`
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *RequestRepository) GetAllRequests(aliasName string) (RequestsWithCount, error) {
	var requests []models.Request
	var count int64
	tx := r.db.Find(&requests, "alias_name = ?", aliasName).Count(&count)
	return RequestsWithCount{
		Requests: requests,
		Count:    count,
	}, tx.Error
}
