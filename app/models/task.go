package models

type Task struct {
	Id          int       `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	InDeps      *TaskList `json:"in_dependencies"`
	OutDeps     *TaskList `json:"out_dependencies"`
}

type TaskList []*Task
