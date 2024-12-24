package interfaces

import (
	"charts/controller"
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {}

type Options struct {
    Message  string
    Data   interface{}
}

func (server HttpServer) Response (c echo.Context, options Options) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
        "message": options.Message,
        "data":    options.Data,
    })
}

func (server HttpServer) HandleHttp(controller *controller.Controller) {
	e := echo.New()

	userGroup := e.Group("/user")
	projectGroup := e.Group("/project")
	issueGroup := e.Group("/issue")

		// ***
	// USER

	userGroup.GET("/list", func(c echo.Context) error {
		users, err := controller.Repo.ListUser()
		if err != nil {
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
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		id, err := controller.CreateUser(newUser.Email)
		if err != nil {
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data:    map[string]interface{}{"id": id},
		})
	})

	userGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
    		return server.Response(c, Options{
        		Message: "invalid ID",
    		})
		}
		id := uint(idInt)

		err = controller.DeleteUser(id)
		if err != nil {
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
			return server.Response(c, Options{
				Message: "cann't finde projects",
			})
		}
		return server.Response(c, Options{
			Data:   projects,
		})
	})

	projectGroup.POST("/add", func(c echo.Context) error {
		newProject := new(project.Project)
		if err := c.Bind(newProject); err != nil {
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		id, err := controller.CreateProject(newProject.Name, newProject.Blocked)
		if err != nil {
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data:    map[string]interface{}{"id": id},
		})
	})

	projectGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
    		return server.Response(c, Options{
        		Message: "invalid ID",
    		})
		}
		id := uint(idInt)

		err = controller.DeleteProject(id)
		if err != nil {
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
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		deadline, err := time.Parse("02-01-2006", dto.Deadline)
		if err != nil {
			return server.Response(c, Options{
				Message: "invalid deadline format",
			})
		}

		newProject, err := controller.Repo.GetProject(dto.ProjectID)
		if err != nil {
			return server.Response(c, Options{
				Message: "project search error",
			})
		}

		newUser, err := controller.Repo.GetUser(dto.UserID)
		if err != nil {
			return server.Response(c, Options{
				Message: "user search error",
			})
		}

		users, err := controller.Repo.UsersByID(dto.Watchers)
		if err != nil {
			return server.Response(c, Options{
				Message: "users search error",
			})
		}

		id, err := controller.CreateIssue(dto.Title, *newUser, *newProject, dto.Priority, dto.Status, deadline, users)
		if err != nil {
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data:    map[string]interface{}{"id": id},
		})
	})

	issueGroup.DELETE("/delete", func(c echo.Context) (err error) {
		idParam := c.QueryParam("id")
		idInt, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
    		return server.Response(c, Options{
        		Message: "invalid ID",
    		})
		}
		id := uint(idInt)

		err = controller.DeleteIssue(id)
		if err != nil {
			return server.Response(c, Options{
				Message: "issue not found",
			})
		}

		return server.Response(c, Options{
			Message: "issue was deleted",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}