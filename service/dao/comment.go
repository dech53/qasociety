package dao

import "qasociety/model"

// AddComment 添加回复
func AddComment(comment model.Comment) error {
	return DB.Create(&comment).Error
}
func FindCommentsByID(answerID, offset, pageSize int, order string) ([]model.Comment, error) {
	var comments []model.Comment
	err := DB.Where("answer_id = ?", answerID).Order("created_time " + order).Offset(offset).Limit(pageSize).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// FindCommentByID 通过评论ID查找评论
func FindCommentByID(commentID int) (*model.Comment, error) {
	var comment model.Comment
	err := DB.Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// RemoveCommentByID 通过评论ID删除评论
func RemoveCommentByID(commentID int) error {
	err := DB.Where("id = ?", commentID).Delete(&model.Comment{}).Error
	if err != nil {
		return err
	}
	return nil
}
func GetAllComments(AnswerID int) ([]model.Comment, error) {
	var comments []model.Comment
	err := DB.Where("answer_id = ?", AnswerID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// RemoveComments 删除多条评论
func RemoveComments(comments []model.Comment) error {
	return DB.Delete(&comments).Error
}
