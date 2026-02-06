package repository

import (
	"technical-test/src/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestRepository interface {
	Create(request *model.Request) error
	FindByID(id int) (model.Request, error)
	FindByIDWithLock(tx *gorm.DB, id int) (model.Request, error)
	FindPendingByWorkflowID(workflowID int) (model.Request, error)
	FindAllWithPagination(offset, limit int, search, status string) ([]model.Request, int64, error)
	Update(request *model.Request) error
	UpdateTx(tx *gorm.DB, request *model.Request) error
	BeginTransaction() *gorm.DB
}

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(request *model.Request) error {
	return r.db.Create(request).Error
}

func (r *requestRepository) FindByID(id int) (model.Request, error) {
	var request model.Request
	err := r.db.First(&request, id).Error
	return request, err
}

func (r *requestRepository) FindByIDWithLock(tx *gorm.DB, id int) (model.Request, error) {
	var request model.Request
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&request, id).Error
	return request, err
}

func (r *requestRepository) FindPendingByWorkflowID(workflowID int) (model.Request, error) {
	var request model.Request
	err := r.db.Where("workflow_id = ? AND status = ?", workflowID, "PENDING").First(&request).Error
	return request, err
}

func (r *requestRepository) FindAllWithPagination(offset, limit int, search, status string) ([]model.Request, int64, error) {
	var requests []model.Request
	var total int64

	query := r.db.Model(&model.Request{})
	if search != "" {
		query = query.Where("workflow_id = ?", search)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return requests, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error

	return requests, total, err
}

func (r *requestRepository) Update(request *model.Request) error {
	return r.db.Save(request).Error
}

func (r *requestRepository) UpdateTx(tx *gorm.DB, request *model.Request) error {
	return tx.Save(request).Error
}

func (r *requestRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
