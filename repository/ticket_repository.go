package repository

import (
	"errors"
	"time"

	"ticket.narindo.com/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	Create(newTicket *model.Ticket) error
	FindBy(by map[string]interface{}) (model.Ticket, error)
	FindAllBy(by map[string]interface{}) ([]model.Ticket, error)
	FindAll() ([]model.Ticket, error)
	UpdateBy(by map[string]interface{}, value map[string]interface{}) error
	Delete(ticket *model.Ticket) error

	FindWithQuery(by *gorm.DB) ([]model.Ticket, error)
	SortBy(filterBy, orderBy map[string]interface{}) ([]model.Ticket, error)
	SelectBy(selectBy string, filterBy, orderBy map[string]interface{}) ([]model.Ticket, error)
	WhereClauseByDate(csId *string, statusId, picId *int, startDate, endDate time.Time) *gorm.DB
}

type ticketRepository struct {
	db *gorm.DB
}

func (m *ticketRepository) Create(newTicket *model.Ticket) error {
	result := m.db.Create(newTicket).Error
	return result
}

func (t *ticketRepository) FindBy(by map[string]interface{}) (model.Ticket, error) {
	var ticket model.Ticket
	res := t.db.Preload(clause.Associations).Where(by).First(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) FindAllBy(by map[string]interface{}) ([]model.Ticket, error) {
	var ticket []model.Ticket
	res := t.db.Preload(clause.Associations).Where(by).Find(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) FindAll() ([]model.Ticket, error) {
	var ticket []model.Ticket
	res := t.db.Preload(clause.Associations).Find(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) UpdateBy(by map[string]interface{}, value map[string]interface{}) error {
	return t.db.Model(model.Ticket{}).Where(by).Updates(value).Error
}

func (t *ticketRepository) Delete(ticket *model.Ticket) error {
	res := t.db.Delete(&model.Ticket{}, ticket).Error
	return res
}

func (t *ticketRepository) FindWithQuery(by *gorm.DB) ([]model.Ticket, error) {
	var ticket []model.Ticket
	res := t.db.Preload(clause.Associations).Where(by).Find(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) SortBy(filterBy, orderBy map[string]interface{}) ([]model.Ticket, error) {
	var ticket []model.Ticket
	res := t.db.Preload(clause.Associations).Where(filterBy).Order(orderBy).Find(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) SelectBy(selectBy string, filterBy, orderBy map[string]interface{}) ([]model.Ticket, error) {
	var ticket []model.Ticket
	res := t.db.Preload(clause.Associations).Select(selectBy).Where(filterBy).Order(orderBy).Find(&ticket)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, nil
		} else {
			return ticket, err
		}
	}
	return ticket, nil
}

func (t *ticketRepository) WhereClauseByDate(csId *string, statusId, picId *int, startDate, endDate time.Time) *gorm.DB {
	if csId != nil && statusId != nil && picId != nil {
		return t.db.Where("status_id = ? AND created_at >= ? AND created_at <= ? AND (?::varchar IS NULL or cs_id = ?) AND (?:: int IS NULL or pic_id = ?)", statusId, startDate, endDate, csId, csId, picId, picId)
	} else if csId != nil && statusId != nil {
		return t.db.Where("status_id = ? AND created_at >= ? AND created_at <= ? AND (?::varchar IS NULL or cs_id = ?)", statusId, startDate, endDate, csId, csId, nil, nil)
	} else if statusId != nil {
		return t.db.Where("status_id = ? AND created_at >= ? AND created_at <= ? AND (?::varchar IS NULL or cs_id = ?)", statusId, startDate, endDate, nil, nil, nil, nil)
	}
	return nil
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	repo := new(ticketRepository)
	repo.db = db
	return repo
}
