package handler

import (
	"fmt"
	"net/http"

	session "github.com/daisuke310vvv/echo-session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/yusk/todo-sample/model"
	"gopkg.in/go-playground/validator.v9"
)

type sessionParam struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

func GetSignIn(c echo.Context) error {
	mapData := map[string]interface{}{}
	mapData["CSRF"] = c.Get("csrf").(string)
	return c.Render(http.StatusOK, "session/signin", mapData)
}

func PostSignIn(c echo.Context) error {
	var p sessionParam
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

	var user model.User

	res := db.Where(model.User{Email: p.Email, Password: p.Password}).First(&user)
	if len(res.GetErrors()) > 0 {
		fmt.Println(res.GetErrors())
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	fmt.Println(user)

	sess := model.Session{UserID: user.ID}
	s := session.Default(c)
	s.Set("session", &sess)
	err = s.Save()
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	return c.Redirect(http.StatusFound, "/")
}

func GetSignUp(c echo.Context) error {
	mapData := map[string]interface{}{}
	mapData["CSRF"] = c.Get("csrf").(string)
	return c.Render(http.StatusFound, "session/signup", mapData)
}

func PostSignUp(c echo.Context) error {
	var p sessionParam
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
	user := model.User{Email: p.Email, Password: p.Password}
	tx := db.Begin()

	res := tx.Create(&user)
	if len(res.GetErrors()) > 0 {
		tx.Rollback()
		fmt.Println(res.GetErrors())
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	tx.Commit()

	fmt.Println(user)

	sess := model.Session{UserID: user.ID}
	s := session.Default(c)
	s.Set("session", &sess)
	err = s.Save()
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusFound, c.Request().URL.String())
	}

	return c.Redirect(http.StatusFound, "/")
}

func GetSignOut(c echo.Context) error {
	s := session.Default(c)
	s.Clear()
	s.Save()
	return c.Redirect(http.StatusFound, "/")
}
