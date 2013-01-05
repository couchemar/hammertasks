package controllers

import (
	"errors"
	"github.com/davemeehan/Neo4j-GO"
	"github.com/robfig/revel"
	"hammertasks/app/models"
)

type Application struct {
	*rev.Controller
}

func (c Application) AddUser() rev.Result {
	if user, _ := c.connected(); user != nil {
		c.RenderArgs["user"] = user
	}
	return nil
}

func (c Application) connected() (*models.User, error) {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User), nil
	}
	if login, ok := c.Session["login"]; ok {
		return c.getUser(login)
	}
	return nil, nil
}

func (c Application) Index() rev.Result {
	if user, _ := c.connected(); user != nil {
		return c.Redirect(Tasks.Index)
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c Application) getUser(login string) (*models.User, error) {
	// TODO избавиться от этой зависимости.
	neo, err := neo4j.NewNeo4j("http://localhost:7474/db/data", "", "")

	if err != nil {
		rev.WARN.Println("Could not connect to Neo4j: ", err)
		return nil, err
	}

	users, err := neo.SearchIdx("login", login, "", "users", "node")
	if err != nil {
		rev.WARN.Println("Error when search: ", err)
		return nil, err
	}
	userData := users[0]

	if userData == nil || userData.Data == nil {
		rev.INFO.Println("Not found")
		return nil, errors.New("Not found")
	}

	user := models.User{
		Id:       userData.ID,
		Login:    userData.Data["login"].(string),
		Password: userData.Data["password"].(string),
	}

	return &user, nil
}

func (c Application) Login(login, password string) rev.Result {

	user, err := c.getUser(login)

	if err != nil || user.Password != password {
		c.Flash.Error("Wrong login and password")
		return c.Redirect(Application.Index)
	}

	c.Session["id"] = string(user.Id)
	c.Session["login"] = login
	c.Flash.Success("Welcome, " + login)
	return c.Redirect(Tasks.Index)
}

func (c Application) Logout() rev.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}
