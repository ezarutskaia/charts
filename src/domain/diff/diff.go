package diff

import (
	"charts/domain/issue"
	"gorm.io/gorm"
)

type CommentsDiff struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	Diff []byte `gorm:"type:json"`
	IssueID uint
	Issue issue.Issue `gorm:"foreignKey:IssueID"`
	Result []byte `gorm:"type:json"`
}
