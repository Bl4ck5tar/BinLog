package service

import (
	"BinLog/server/global"
	"BinLog/server/model/database"

	"gorm.io/gorm"
)

//LoadChildren 加载该评论下的所有子评论
func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment

	//查找该评论的所有子评论
	if err := global.DB.
	Where("p_id = ?", comment.ID).
	Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).
	Find(&children).Error; err != nil {
		return err
	}

	//如果有子评论，递归加载
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}
	//将子评论绑定到当前评论
	comment.Children = children
	return nil
}

//DeleteCommentAndChildren 根据id删除该评论及其所有子评论
func (commentService *CommentService) DeleteCommentAndChildren(tx *gorm.DB, commentID uint) error {
	var children []database.Comment
	if err := tx.Where("p_id = ?", commentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := commentService.DeleteCommentAndChildren(tx, child.ID); err != nil {
			return err
		}
	}

	if err := tx.Delete(&database.Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}

func (commentService *CommentService) FindChildCommentsIDByRootCommentUserUUID(comments []database.Comment) map[uint]struct{} {
	result := make(map[uint]struct{})

	for _, rootComment := range comments {
		var findChildren func([]database.Comment)

		findChildren = func(children []database.Comment) {
			for _, child := range children {
				if child.UserUUID == rootComment.UserUUID{
					result[child.ID] = struct{}{}
				}
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}
		}
		findChildren(rootComment.Children)
	}
	return result
}