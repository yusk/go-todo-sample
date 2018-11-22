package main

import (
	"fmt"
	"html/template"
	"io"

	session "github.com/daisuke310vvv/echo-session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yusk/todo-sample/handler"
	"github.com/yusk/todo-sample/model"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	tpl := &Template{
		templates: template.Must(template.ParseGlob("views/**/*.html")),
	}

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err.Error())
	}
	if res := db.AutoMigrate(
		&model.User{},
		&model.Todo{},
	); len(res.GetErrors()) > 0 {
		panic(res.GetErrors()[0].Error())
	}
	if res := db.Model(&model.Todo{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE"); len(res.GetErrors()) > 0 {
		fmt.Println(res.GetErrors()[0].Error())
	}

	store, err := session.NewRedisStore(32, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		panic(err.Error())
	}

	e := echo.New()
	e.Renderer = tpl

	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))
	e.Use(session.Sessions("session", store))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_csrf",
	}))

	e.GET("/sample/string", handler.SampleString)
	e.GET("/sample/json", handler.SampleJSON)
	e.GET("/sample/html", handler.SampleHTML)

	e.GET("/signin", handler.GetSignIn)
	e.GET("/signup", handler.GetSignUp)
	e.GET("/signout", handler.GetSignOut)
	e.POST("/signin", handler.PostSignIn)
	e.POST("/signup", handler.PostSignUp)

	e.GET("/", handler.TodoList)
	e.GET("/:id", handler.TodoShow)
	e.POST("/", handler.TodoCreate)
	e.GET("/new", handler.TodoNew)

	e.Start(":9090")
}
