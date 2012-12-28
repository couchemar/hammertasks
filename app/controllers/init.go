package controllers

import "github.com/robfig/revel"

func init() {
	rev.InterceptMethod(Application.AddUser, rev.BEFORE)
}
