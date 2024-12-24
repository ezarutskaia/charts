package domain

import (
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"time"
)

type Domain struct {
}

func (domain *Domain) CreateIssue(title string, user user.User, project project.Project, priority int, status string, deadline time.Time, watchers []user.User) *issue.Issue{
	return &issue.Issue{Title: title, User: user, Project: project, Priority: priority, Status: status, Deadline: deadline, Watchers: watchers}
}

func (domain *Domain) CreateProject(name string, blocked bool) *project.Project {
	return &project.Project{Name: name, Blocked: blocked}
}

func (domain *Domain) CreateUser(email string) *user.User {
	return &user.User{Email: email}
}