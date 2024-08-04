package dao

import "qasociety/model"

func AddAnswer(answer model.Answer) error {
	return DB.Model(&model.Answer{}).Create(&answer).Error
}
