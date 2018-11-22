package handler

import (
	"fmt"
	"net/http"
	"strconv"

	session "github.com/daisuke310vvv/echo-session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/yusk/todo-sample/model"
	validator "gopkg.in/go-playground/validator.v9"
)

type todoParam struct {
	Title   string `json:"email" validate:"required"`
	Content string `json:"password"`
}

func TodoList(c echo.Context) error {
	mapData := map[string]interface{}{}
	failureURL := "/signin"

	s := session.Default(c)
	v := s.Get("session")
	if v == nil {
		return c.Redirect(http.StatusFound, failureURL)
	}
	sess, ok := v.(model.Session)
	if !ok {
		return c.Redirect(http.StatusFound, failureURL)
	}

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, failureURL)
	}

	var todos []*model.Todo
	res := db.Where(model.Todo{UserID: sess.UserID}).Find(&todos)
	if len(res.GetErrors()) > 0 {
		fmt.Println(res.GetErrors())
		return c.Redirect(http.StatusFound, failureURL)
	}

	mapData["UserID"] = sess.UserID
	mapData["Todos"] = todos
	return c.Render(http.StatusOK, "todo/list", mapData)
}

func TodoShow(c echo.Context) error {
	mapData := map[string]interface{}{}
	failureURL := "/"
	idStr := c.Param("id")

	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, failureURL)
	}
	id := uint(id64)

	s := session.Default(c)
	v := s.Get("session")
	if v == nil {
		return c.Redirect(http.StatusFound, failureURL)
	}
	sess, ok := v.(model.Session)
	if !ok {
		return c.Redirect(http.StatusFound, failureURL)
	}

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, failureURL)
	}

	var todos []*model.Todo
	res := db.Where(model.Todo{UserID: sess.UserID, ID: id}).Find(&todos)
	if len(res.GetErrors()) > 0 {
		fmt.Println(res.GetErrors())
		return c.Redirect(http.StatusFound, failureURL)
	}

	if len(todos) == 0 {
		return c.Redirect(http.StatusFound, failureURL)
	}

	mapData["Todo"] = todos[0]
	return c.Render(http.StatusOK, "todo/show", mapData)
}

func TodoNew(c echo.Context) error {
	mapData := map[string]interface{}{}

	s := session.Default(c)
	v := s.Get("session")
	if v == nil {
		return c.Redirect(http.StatusFound, "/signin")
	}
	sess, ok := v.(model.Session)
	if !ok {
		return c.Redirect(http.StatusFound, "/signin")
	}

	mapData["UserID"] = sess.UserID
	mapData["CSRF"] = c.Get("csrf").(string)
	return c.Render(http.StatusOK, "todo/new", mapData)
}

func TodoCreate(c echo.Context) error {
	s := session.Default(c)
	v := s.Get("session")
	if v == nil {
		return c.Redirect(http.StatusFound, "/signin")
	}
	sess, ok := v.(model.Session)
	if !ok {
		return c.Redirect(http.StatusFound, "/signin")
	}

	var p todoParam
	err := c.Bind(&p)
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}
	err = validator.New().Struct(p)
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}
	fmt.Println(p)
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	todo := model.Todo{UserID: sess.UserID, Title: p.Title, Content: p.Content}
	tx := db.Begin()

	res := tx.Create(&todo)
	if len(res.GetErrors()) > 0 {
		tx.Rollback()
		fmt.Println(res.GetErrors())
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	tx.Commit()

	return c.Redirect(http.StatusFound, "/")
}
