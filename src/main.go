package main

import (
	"charts/controller"
	"charts/domain"
	"charts/domain/diff"
	"charts/domain/issue"
	"charts/domain/project"
	"charts/domain/user"
	"charts/infra"
	"charts/interfaces"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type App struct {
	Domain *domain.Domain
	Infra *infra.Repository
	Interfaces *interfaces.HttpServer
}

func main() {
	dsn := "root:secret@tcp(sql_charts:3306)/charts?parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
        log.Fatal(err)
    }
	err = (*db).AutoMigrate(&issue.Issue{}, &user.User{}, &project.Project{}, &diff.CommentsDiff{})
	if err != nil {
		fmt.Println(err)
	}

	app := &App{
		Domain: &domain.Domain{},
		Infra: &infra.Repository{
			DB: db,
			},
		Interfaces: &interfaces.HttpServer{},
	}

	app.Interfaces.HandleHttp(&controller.Controller{
		Repo: app.Infra,
		Domain: app.Domain,
	})
}
