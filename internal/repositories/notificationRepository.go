package repositories

import (
	"github.com/plyovchev/sumup-assignment-notifications/internal/db"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
	"github.com/plyovchev/sumup-assignment-notifications/internal/util"
)

type NotificationRepository interface {
	Create(notification *data.Notification) (*data.Notification, error)
	FindAll() (*[]data.Notification, error)
	FindAllByIds(ids []int) (*[]data.Notification, error)
	FindAllByStatus(status data.NotificationStatus) (*[]data.Notification, error)
	Save(notification *data.Notification) (*data.Notification, error)
}

type noticationRepository struct {
	dbClient db.DbClient
}

func NewNotificationRepository(dbClient db.DbClient) NotificationRepository {
	return &noticationRepository{
		dbClient: dbClient,
	}
}

// Create persists this category data.
func (repository *noticationRepository) Create(notification *data.Notification) (*data.Notification, error) {
	if err := repository.dbClient.Create(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

// FindAll returns all categories of the category table.
func (repository *noticationRepository) FindAll() (*[]data.Notification, error) {
	var notifications []data.Notification
	if err := repository.dbClient.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return &notifications, nil
}

// Returns all notifications in specified status.
func (repository *noticationRepository) FindAllByIds(ids []int) (*[]data.Notification, error) {
	var notifications []data.Notification
	idsArr := util.ArrayToString(ids, ", ")
	if err := repository.dbClient.Where("id in (?)", idsArr).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return &notifications, nil
}

// Returns all notificaitons in specified status.
func (repository *noticationRepository) FindAllByStatus(status data.NotificationStatus) (*[]data.Notification, error) {
	var notifications []data.Notification
	if err := repository.dbClient.Where("id = (?)", status).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return &notifications, nil
}

// Save persists this book data.
func (repository *noticationRepository) Save(notification *data.Notification) (*data.Notification, error) {
	if err := repository.dbClient.Save(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}