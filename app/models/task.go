package models

type Task struct {
	Id          int `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type TaskList []Task
