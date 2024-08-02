package dao

import "qasociety/model"

func AddQuestion(question model.Question) error {
	return DB.Model(&model.Question{}).Create(&question).Error
}
