package dao

import "qasociety/model"

// AddComment 添加回复
func AddComment(comment model.Comment) error {
	return DB.Create(&comment).Error
}
func FindCommentsByID(answerID, offset, pageSize int) ([]model.Comment, error) {
	var comments []model.Comment
	err := DB.Where("answer_id = ?", answerID).Offset(offset).Limit(pageSize).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
