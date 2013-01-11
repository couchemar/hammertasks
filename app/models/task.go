package models

type Task struct {
	Id          int
	Summary     string
	Description string
}

type TaskList []Task
