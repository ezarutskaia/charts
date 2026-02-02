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
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type App struct {
	Domain *domain.Domain
	Infra *infra.Infra
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

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	app := &App{
		Domain: &domain.Domain{},
		Infra: &infra.Infra{
			Repository: &infra.Repository{
				DB: db,
			},
			Redis: &infra.RedisRepository{
				Client: rdb,
			},
		},
		Interfaces: &interfaces.HttpServer{},
	}

	app.Interfaces.HandleHttp(&controller.Controller{
		Repo: app.Infra.Repository,
		Domain: app.Domain,
		Redis: app.Infra.Redis,
	})
}
