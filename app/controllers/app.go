package controllers

import "github.com/robfig/revel"

type Application struct {
	*rev.Controller
}

func (c Application) Index() rev.Result {
	return c.Render()
}

func (c Application) Login() rev.Result {
	return c.Render()
}