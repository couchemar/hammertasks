package controllers

import (
	"github.com/robfig/revel"
	"github.com/davemeehan/Neo4j-GO"
)

type Application struct {
	*rev.Controller
}

func (c Application) Index() rev.Result {
	return c.Render()
}

func (c Application) Login(login, password string) rev.Result {
	// TODO избавиться от этой зависимости.
	neo, err := neo4j.NewNeo4j("http://localhost:7474/db/data", "", "")

	if err != nil {
		rev.WARN.Println("Could not connect to Neo4j: ", err)
		return c.Redirect(Application.Index)
	}

	users, err := neo.SearchIdx("login", login, "", "users", "node")
	if err != nil {
		rev.WARN.Println("Error when search: ", err)
	}

	user := users[0]

	if user == nil {
		rev.INFO.Println("Not found")
		c.Flash.Error("Login failed")
		return c.Redirect(Application.Index)
	}

	rev.INFO.Println("Data: ", user.Data)

	if user.Data["password"] != password {
		c.Flash.Error("Wrong login and password")
		return c.Redirect(Application.Index)
	}

	c.Session["login"] = login
	c.Flash.Success("Welcome, "+ login)
	return c.Redirect(Application.Index)
}