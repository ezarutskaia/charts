package controller

import (
	"charts/domain"
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"charts/infra"
	"time"
)

type Controller struct {
	Repo *infra.Repository
	Domain *domain.Domain
}

func (controller *Controller) CreateIssue(title string, user user.User, project project.Project, priority int, status string, deadline time.Time, watchers []user.User) (id uint, err error) {
	newIssue := controller.Domain.CreateIssue(title, user, project, priority, status, deadline, watchers)
	id, err = controller.Repo.CreateIssue(newIssue)
	return
}

func (controller *Controller) CreateIssues(issues []issue.Issue) error {
	err := controller.Repo.CreateIssues(issues)
	return err
}

func (controller *Controller) CreateProject(name string, blocked ...bool) (id uint, err error) {
	blockedValue := false
	if len(blocked) > 0 {
		blockedValue = blocked[0]
	}
	newProject := controller.Domain.CreateProject(name, blockedValue)
	id, err = controller.Repo.CreateProject(newProject)
	return
}

func (controller *Controller) CreateProjects(projects []project.Project) error {
	err := controller.Repo.CreateProjects(projects)
	return err
}

func (controller *Controller) CreateUser(email string) (id uint, err error) {
	newUser := controller.Domain.CreateUser(email)
	id, err = controller.Repo.CreateUser(newUser)
	return
}

func (controller *Controller) CreateUsers(users []user.User) error {
	err := controller.Repo.CreateUsers(users)
	return err
}

func (controller *Controller) DeleteIssue(id uint) error {
	err := controller.Repo.DeleteIssue(id)
	return err
}

func (controller *Controller) DeleteProject(id uint) error {
	err := controller.Repo.DeleteProject(id)
	return err
}

func (controller *Controller) DeleteUser(id uint) error {
	err := controller.Repo.DeleteUser(id)
	return err
}