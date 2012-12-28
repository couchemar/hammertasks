package controllers

import (
	"github.com/robfig/revel"
)

type Tasks struct {
	Application
}

func (c Tasks) Index() rev.Result {
	return c.Render()
}
