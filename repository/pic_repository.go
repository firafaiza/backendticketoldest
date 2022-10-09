package repository

import (
	"errors"

	"ticket.narindo.com/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type picRepository struct {
	repo *gorm.DB
}

type PicRepository interface {
	Insert(pic *model.Pic) error
	FindBy(by map[string]interface{}) (model.Pic, error)
	FindAllBy(by map[string]interface{}) ([]model.Pic, error)
	FindAllWithPreload(by map[string]interface{}, preload ...string) ([]model.Pic, error)
	UpdateBy(by map[string]interface{}, value map[string]interface{}) error
}

func (p *picRepository) Insert(pic *model.Pic) error {
	res := p.repo.Create(pic).Error
	return res
}

func (t *picRepository) FindBy(by map[string]interface{}) (model.Pic, error) {
	var pic model.Pic
	res := t.repo.Preload(clause.Associations).Where(by).First(&pic)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pic, nil
		} else {
			return pic, err
		}
	}
	return pic, nil
}

func (p *picRepository) FindAllBy(by map[string]interface{}) ([]model.Pic, error) {
	var pics []model.Pic
	res := p.repo.Where(by).Find(&pics)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pics, nil
		} else {
			return pics, err
		}
	}
	return pics, nil
}

func (p *picRepository) FindAllWithPreload(by map[string]interface{}, preload ...string) ([]model.Pic, error) {
	var pics []model.Pic
	res := p.repo.Preload(preload[0]).Preload(preload[1]).Preload(preload[2]).Where(by).Find(&pics)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pics, nil
		} else {
			return pics, err
		}
	}
	return pics, nil
}

func (p *picRepository) UpdateBy(by map[string]interface{}, value map[string]interface{}) error {
	return p.repo.Model(model.Pic{}).Where(by).Updates(value).Error
}

func InitPicRepository(db *gorm.DB) PicRepository {
	picRepo := new(picRepository)
	picRepo.repo = db
	return picRepo
}
