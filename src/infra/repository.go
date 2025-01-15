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

type IdCount struct {
	Reason string  `gorm:"column:Reason"`
	Count  int     `gorm:"column:Count"`
}

type StatusCount struct {
	IssueStatus string
	Count  int
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

func (repo *Repository) ListUser () (users []*user.DTOUser, err error) {
	result := (*repo.DB).Model(&user.User{}).Select("id", "email").Find(&users)
	return users, result.Error
}

func (repo *Repository) ListProject () (projects []*project.DTOProject, err error) {
	result := (*repo.DB).Model(&project.Project{}).Select("id", "name").Find(&projects)
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

func (repo *Repository) CountIssuesGroup(groupby string) (map[string]int, error){
	var results []IdCount
	idCountMap := map[string]int{}

	result := (*repo.DB).Model(&issue.Issue{})

	switch groupby {
	case "user":
		result = result.Select("CAST(user_id AS CHAR) as Reason, count(id) as Count").Group("user_id")
	case "project":
		result = result.Select("CAST(project_id AS CHAR) as Reason, count(id) as Count").Group("project_id")
	case "priority":
		result = result.Select("CAST(priority AS CHAR) as Reason, count(id) as Count").Group("priority")
	case "status":
		result = result.Select("status as Reason, count(id) as Count").Group("status")
	}

	result = result.Scan(&results)

	for _, item := range results {
		idCountMap[item.Reason] = item.Count
	}

	return idCountMap, result.Error
}

func (repo *Repository) CountIssuesProjects() (map[interface{}]int, error){
	var results []IdCount
	idCountMap := map[interface{}]int{}

	result := (*repo.DB).Model(&issue.Issue{}).
		Select("project_id as ItemID, count(id) as Count").
		Group("ItemID").
		Scan(&results)

	for _, item := range results {
		idCountMap[item.Reason] = item.Count
	}

	return idCountMap, result.Error
}

func (repo *Repository) CountIssuesPriority() (map[interface{}]int, error){
	var results []IdCount
	idCountMap := map[interface{}]int{}

	result := (*repo.DB).Model(&issue.Issue{}).
		Select("priority as ItemID, count(id) as Count").
		Group("ItemID").
		Scan(&results)

	for _, item := range results {
		idCountMap[item.Reason] = item.Count
	}

	return idCountMap, result.Error
}

func (repo *Repository) CountIssuesStatus() (map[string]int, error){
	var results []StatusCount
	idCountMap := map[string]int{}

	result := (*repo.DB).Model(&issue.Issue{}).
		Select("status as IssueStatus, count(id) as Count").
		Group("IssueStatus").
		Scan(&results)

	for _, item := range results {
		idCountMap[item.IssueStatus] = item.Count
	}

	return idCountMap, result.Error
}
