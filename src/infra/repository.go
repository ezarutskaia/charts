package infra

import (
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type UserTaskCount struct {
		UserID     uint  `json:"user_id"`
        IssueCount int64 `json:"issue_count"`
	}

func (repo *Repository) CreateIssue(issue *issue.Issue) (uint, error) {
	result := (*repo.DB).Create(issue)
	return issue.ID, result.Error
}

func (repo *Repository) CreateIssues(issues []issue.Issue) error {
	result := (*repo.DB).Create(&issues)
	return result.Error
}

func (repo *Repository) CreateUser(user *user.User) (uint, error) {
	result := (*repo.DB).Create(user)
	return user.ID, result.Error
}

func (repo *Repository) CreateUsers(users []user.User) error {
	result := (*repo.DB).Create(&users)
	return result.Error
}

func (repo *Repository) CreateProject(project *project.Project) (uint, error) {
	result := (*repo.DB).Create(project)
	return project.ID, result.Error
}

func (repo *Repository) CreateProjects(projects []project.Project) error {
	result := (*repo.DB).Create(&projects)
	return result.Error
}

func (repo *Repository) DeleteIssue (id uint) error {
	result := (*repo.DB).Delete(&issue.Issue{}, id)
	return result.Error
}

func (repo *Repository) DeleteUser (id uint) error {
	result := (*repo.DB).Delete(&user.User{}, id)
	return result.Error
}

func (repo *Repository) DeleteProject (id uint) error {

	result := (*repo.DB).Delete(&project.Project{}, id)
	return result.Error
}

func (repo *Repository) ListIssue () (issues []*issue.Issue, err error) {
	result := (*repo.DB).Preload("Watchers").Find(&issues)
	return issues, result.Error
}

func (repo *Repository) ListUser () (users []*user.User, err error) {
	result := (*repo.DB).Find(&users)
	return users, result.Error
}

func (repo *Repository) ListProject () (projects []*project.Project, err error) {
	result := (*repo.DB).Find(&projects)
	return projects, result.Error
}

func (repo *Repository) GetIssue (id uint) (issue *issue.Issue, err error) {
	result := (*repo.DB).Where("id = ?", id).First(&issue)
	return issue, result.Error
}

func (repo *Repository) GetUser (id uint) (user *user.User, err error) {
	result := (*repo.DB).Where("id = ?", id).First(&user)
	return user, result.Error
}

func (repo *Repository) GetProject (id uint) (project *project.Project, err error) {
	result := (*repo.DB).Where("id = ?", id).First(&project)
	return project, result.Error
}

func (repo *Repository) UsersByID (ids []uint) (users []user.User, err error) {
	result := (*repo.DB).Find(&users, ids)
	return users, result.Error
}

func (repo *Repository) CountIssues() (count int64, err error) {
	result := (*repo.DB).Model(&issue.Issue{}).Count(&count)
	return count, result.Error
}

func (repo *Repository) CountProjects() (count int64, err error) {
	result := (*repo.DB).Model(&project.Project{}).Count(&count)
	return count, result.Error
}

func (repo *Repository) CountUsers() (count int64, err error) {
	result := (*repo.DB).Model(&user.User{}).Count(&count)
	return count, result.Error
}

func (repo *Repository) CountIssuesUsers() ([]UserTaskCount, error){
	var results []UserTaskCount

	result := (*repo.DB).Model(&issue.Issue{}).
		Select("user_id, count(id) as issue_count").
		Group("user_id").
		Order("issue_count desc").
		Scan(&results)

	return results, result.Error
}