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

func (c Application) Index() rev.Result {
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
	user_data := users[0]

	if user_data == nil {
		rev.INFO.Println("Not found")
		return nil, errors.New("Not found")
	}

	user := models.User{
		Id:       user_data.ID,
		Login:    user_data.Data["login"].(string),
		Password: user_data.Data["password"].(string),
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
