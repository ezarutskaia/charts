package interfaces

import (
	"charts/controller"
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"charts/helpers"
	"context"
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
	"time"
)

const batchSize int = 1000

type HttpServer struct{}

type ChartsRequest struct {
	GroupBy string `json:"groupBy"`
	ChartType string `json:"chartType"`
	Filters []Filter
}

type Filter struct {
	FilterType string `json:"type"`
	Value string      `json:"value"`
}
type Options struct {
	Message string
	Data    interface{}
}

func (server HttpServer) Response(c echo.Context, options Options) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": options.Message,
		"data":    options.Data,
	})
}

func (server HttpServer) HandleHttp(controller *controller.Controller) {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))

	userGroup := e.Group("/user")
	projectGroup := e.Group("/project")
	issueGroup := e.Group("/issue")

	// ***
	// USER

	userGroup.GET("/list", func(c echo.Context) error {
		users, err := controller.Repo.ListUser()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "cann't finde users",
			})
		}
		return server.Response(c, Options{
			Data: users,
		})
	})

	userGroup.POST("/add", func(c echo.Context) error {
		newUser := new(user.User)
		if err := c.Bind(newUser); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		id, err := controller.CreateUser(newUser.Email)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data: map[string]interface{}{"id": id},
		})
	})

	userGroup.POST("/batch", func(c echo.Context) error {
		var users []user.User
		if err := c.Bind(&users); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "invalid JSON payload",
			})
		}

		err := controller.CreateUsers(users)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "Failed to insert users",
			})
		}

		return server.Response(c, Options{
			Message: "Users inserted successfully",
			Data:    map[string]interface{}{"count": len(users)},
		})
	})

	userGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.Logger().Error("Parse error:", err)
			return server.Response(c, Options{
				Message: "invalid ID",
			})
		}
		id := uint(idInt)

		err = controller.DeleteUser(id)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "user not found",
			})
		}

		return server.Response(c, Options{
			Message: "user was deleted",
		})
	})

	// ***
	// PROJECT

	projectGroup.GET("/list", func(c echo.Context) error {
		projects, err := controller.Repo.ListProject()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "cann't finde projects",
			})
		}
		return server.Response(c, Options{
			Data: projects,
		})
	})

	projectGroup.POST("/add", func(c echo.Context) error {
		newProject := new(project.Project)
		if err := c.Bind(newProject); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		id, err := controller.CreateProject(newProject.Name, newProject.Blocked)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data: map[string]interface{}{"id": id},
		})
	})

	projectGroup.POST("/batch", func(c echo.Context) error {
		var projects []project.Project
		if err := c.Bind(&projects); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "invalid JSON payload",
			})
		}

		err := controller.CreateProjects(projects)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "Failed to insert projects",
			})
		}

		return server.Response(c, Options{
			Message: "Projects inserted successfully",
			Data:    map[string]interface{}{"count": len(projects)},
		})
	})

	projectGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.Logger().Error("Parse error:", err)
			return server.Response(c, Options{
				Message: "invalid ID",
			})
		}
		id := uint(idInt)

		err = controller.DeleteProject(id)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "project not found",
			})
		}

		return server.Response(c, Options{
			Message: "project was deleted",
		})
	})

	// ***
	// ISSUE

	issueGroup.GET("/list", func(c echo.Context) error {
		issues, err := controller.Repo.ListIssue()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "cann't finde issues",
			})
		}
		return server.Response(c, Options{
			Data: issues,
		})
	})

	issueGroup.POST("/add", func(c echo.Context) error {
		dto := new(issue.DTOissue)
		if err := c.Bind(dto); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		deadline, err := time.Parse("02-01-2006", dto.Deadline)
		if err != nil {
			c.Logger().Error("Parse error:", err)
			return server.Response(c, Options{
				Message: "invalid deadline format",
			})
		}

		newProject, err := controller.Repo.GetProject(dto.ProjectID)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "project search error",
			})
		}

		newUser, err := controller.Repo.GetUser(dto.UserID)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "user search error",
			})
		}

		users, err := controller.Repo.UsersByID(dto.Watchers)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "users search error",
			})
		}

		id, err := controller.CreateIssue(dto.Title, *newUser, *newProject, dto.Priority, dto.Status, deadline, users)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data: map[string]interface{}{"id": id},
		})
	})

	issueGroup.POST("/batch", func(c echo.Context) error {
		var payloads []issue.DTOissue
		if err := c.Bind(&payloads); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "invalid JSON payload",
			})
		}

		var issues []issue.Issue
		for _, p := range payloads {
			users, _ := controller.Repo.UsersByID(p.Watchers)
			deadline, _ := time.Parse("02-01-2006", p.Deadline)

			issues = append(issues, issue.Issue{
				Title:     p.Title,
				UserID:    p.UserID,
				ProjectID: p.ProjectID,
				Priority:  p.Priority,
				Status:    p.Status,
				Deadline:  deadline,
				Watchers:  users,
			})
		}

		for i := 0; i < len(issues); i += batchSize {
			end := i + batchSize
			if end > len(issues) {
				end = len(issues)
			}

			err := controller.CreateIssues(issues[i:end])
			if err != nil {
				c.Logger().Error("Insert error:", err)
				return server.Response(c, Options{
					Message: "Failed to insert issues",
				})
			}
		}

		return server.Response(c, Options{
			Message: "Issues inserted successfully",
			Data:    map[string]interface{}{"count": len(issues)},
		})
	})

	issueGroup.PATCH("/update", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.Logger().Error("Parse error:", err)
			return server.Response(c, Options{
				Message: "invalid ID",
			})
		}
		id := uint(idInt)

		var jsonBody map[string]interface{}
		if err := c.Bind(&jsonBody); err != nil {
			c.Logger().Error("Bind error:", err)
			return server.Response(c, Options{
				Message: "invalid JSON payload",
			})
		}

		oldIssue, err := controller.Repo.GetIssue(id)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "issue search error",
			})
		}

		updatedIssue := *oldIssue
		err = controller.Repo.UpdateIssue(&updatedIssue, jsonBody)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "issue update error",
			})
		}

		diffID, err := controller.CreateDiff(id, jsonBody, oldIssue)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data: map[string]interface{}{"id": diffID},
		})
	})

	issueGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.Logger().Error("Parse error:", err)
			return server.Response(c, Options{
				Message: "invalid ID",
			})
		}
		id := uint(idInt)

		err = controller.DeleteIssue(id)
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "issue not found",
			})
		}

		return server.Response(c, Options{
			Message: "issue was deleted",
		})
	})

	e.GET("/stat", func(c echo.Context) error {
		userCount, err := controller.Repo.CountUsers()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "row counting error for user",
			})
		}

		projectCount, err := controller.Repo.CountProjects()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "row counting error for project",
			})
		}

		issueCount, err := controller.Repo.CountIssues()
		if err != nil {
			c.Logger().Error("SQL error:", err)
			return server.Response(c, Options{
				Message: "row counting error for issue",
			})
		}

		return server.Response(c, Options{
			Data: map[string]interface{}{
				"count_of_issues":   issueCount,
				"count_of_projects": projectCount,
				"count_of_users":    userCount,
			},
		})
	})

	e.POST("/charts", func(c echo.Context) error {
		ctx := context.Background()
		var req ChartsRequest
		var fields interface{}
		filters := map[string]string{}

		if err := c.Bind(&req); err != nil {
			c.Logger().Error("Bind groupby error:", err)
			return server.Response(c, Options{
				Message: "invalid JSON payload",
			})
		}

		if len(req.Filters) != 0 {
			for _,item := range req.Filters {
				filters[item.FilterType] = item.Value
			}
		}

		for req.ChartType == "bar" || req.ChartType == "" {
			result, err := controller.Repo.CountIssuesGroup(req.GroupBy, filters)
			if err != nil {
				c.Logger().Error("SQL error:", err)
				return server.Response(c, Options{
					Message: "row counting error for issue",
					})
			}

			switch req.GroupBy {
			case "user":
				users, err := controller.Repo.ListUser()
				if err != nil {
					c.Logger().Error("SQL error:", err)
					return server.Response(c, Options{
						Message: "can't found users",
						})
				}
				fields = users

			case "project":
				projects, err := controller.Repo.ListProject()
				if err != nil {
					c.Logger().Error("SQL error:", err)
					return server.Response(c, Options{
						Message: "can't found projects",
						})
				}
				fields = projects

			case "priority":
				fields = nil

			case "status":
				fields = nil
			}

			return server.Response(c, Options{
				Data: map[string]interface{}{
					"groupBy": req.GroupBy,
					"result":  result,
					"fields":  fields,
					},
					})
		}

		for req.ChartType == "line" {
			jsonData, err := json.Marshal(req)
			if err != nil {
				c.Logger().Error("JSON serialization error:", err)
				return server.Response(c, Options{
					Message: "cannot get string to generate CacheKey",
				})
			}
			cacheKey := helpers.GenerateCacheKey(jsonData)
			cachedData, err := controller.Redis.Get(ctx, cacheKey)
			cachedResult := map[time.Time]map[string]int{}
			if err == nil {
				err = json.Unmarshal([]byte(cachedData), &cachedResult)
				if err != nil {
        			c.Logger().Error("Cache JSON decode error:", err)
        			return server.Response(c, Options{
            			Message: "error decoding cached data",
        			})
    			}
				c.Logger().Info("Cache hit for key:", cacheKey)
				return server.Response(c, Options{
					Data: map[string]interface{}{
					"groupBy": req.GroupBy,
					"result":  cachedResult,
					},
				})
			}

			result, err := controller.LineIssues()
			if err != nil {
				c.Logger().Error("SQL error:", err)
				return server.Response(c, Options{
					Message: "row counting error for issue",
				})
			}
			resultJSON, err := json.Marshal(result)
			if err == nil {
				_ = controller.Redis.Set(ctx, cacheKey, string(resultJSON))
			}

			return server.Response(c, Options{
				Data: map[string]interface{}{
					"groupBy": req.GroupBy,
					"result":  result,
					},
				})
		}

		return server.Response(c, Options{
			Message: "Unknown request",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
