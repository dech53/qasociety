package dao

import "qasociety/model"

func AddComment(comment model.Comment) error {
	return DB.Create(&comment).Error
}
