package controller

import (
	"charts/domain"
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"charts/infra"
	"encoding/json"
	"time"
)

type Controller struct {
	Repo *infra.Repository
	Domain *domain.Domain
}

type LinePoint struct {
	 Label string
	 Data []int
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

func (controller *Controller) CreateDiff(issueID uint, jsonBody map[string]interface{}, oldIssue *issue.Issue) (id uint, err error) {

	comment, err := json.Marshal(jsonBody)
	if err != nil {
	    return 0, err
	}

	newIssue, err := controller.Repo.GetIssue(oldIssue.ID)
		if err != nil {
			return 0, err
		}

	newJson := map[string]interface{}{}
	if _, ok := jsonBody["title"].(string); ok {
		oldNewTitle := map[string]interface{}{
	        "old": oldIssue.Title,
	        "new": newIssue.Title,
    	}
    	newJson["title"] = oldNewTitle
	}
	if _, ok := jsonBody["priority"].(float64); ok {
		oldNewPriority := map[string]interface{}{
	        "old": oldIssue.Priority,
	        "new": newIssue.Priority,
	    }
	    newJson["priority"] = oldNewPriority
	}
	if _, ok := jsonBody["status"].(string); ok {
		oldNewStatus := map[string]interface{}{
	        "old": oldIssue.Status,
	        "new": newIssue.Status,
	    }
	    newJson["status"] = oldNewStatus
	}
	if _, ok := jsonBody["watchers"]; ok {
		var oldWatchers []uint
		var newWatchers []uint

		for _,watcher := range oldIssue.Watchers {
			oldWatchers = append(oldWatchers, watcher.ID)
		}
		for _,watcher := range newIssue.Watchers {
			newWatchers = append(newWatchers, watcher.ID)
		}

		oldNewWatchers := map[string]interface{}{
	        "old": oldWatchers,
	        "new": newWatchers,
	    }
	    newJson["watchers"] = oldNewWatchers
	}

	resultComment, err := json.Marshal(newJson)
	if err != nil {
		return 0, err
	}

	newComment := controller.Domain.CreateDiff(comment, issueID, resultComment)
	id, err = controller.Repo.CreateDiff(newComment)
	return
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

func (controller *Controller) LineIssues() ([]LinePoint, error) {
	var points []LinePoint
	filters := []string{"open", "closed", "in_progress", "canceled"}

	for _, filter := range filters {
		count := make([]int, 10)
		for i := range 9 {
		    count[i] = 1
		}
		num, err := controller.Repo.CountIssuesLine(filter)
		if err != nil {
			return nil, err
		}

		count[9] = num
		points = append(points, LinePoint{Label: filter, Data: count})
	}

	return points, nil
}
