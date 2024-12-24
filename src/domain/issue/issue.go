package issue

import (
    "charts/domain/project"
    "charts/domain/user"
    "gorm.io/gorm"
	"time"
)

type  Issue struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	Title string `gorm:"size:256"`
	UserID uint
	User user.User `gorm:"foreignKey:UserID"`
	ProjectID uint
	Project project.Project `gorm:"foreignKey:ProjectID"`
	Priority int `gorm:"check:priority IN (1,2,3,4,5)"`
	Status string `gorm:"type:VARCHAR(20);check:status IN ('open', 'in_progress', 'closed', 'canceled')"`
	Deadline time.Time
	Watchers []user.User `gorm:"many2many:issue_watchers;"`
}