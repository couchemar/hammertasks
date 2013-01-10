package controllers

import (
	"github.com/robfig/revel"
	"github.com/jmcvetta/neo4j"
)

type Tasks struct {
	Application
}

func (c Tasks) Index() rev.Result {
	return c.Render()
}

func (c Tasks) CreateTask(summary, description string) rev.Result {
	neo, err := neo4j.Connect("http://localhost:7474/db/data")
	nodes := neo.Nodes
	rootNode, err := nodes.Get(0)
	if err != nil {
		rev.ERROR.Fatalln("Failed to get root node")
		panic(err)
	}

	/*
	 Получим узел задач, с ним будут связана вновь созданая задача.
	 Это map[int] Relation, но фактически там один элемент.
	 */
	tasksRels, err := rootNode.Outgoing("TASKS")
	if err != nil {
		rev.ERROR.Fatalln("Failed to get TASKS relation")
		panic(err)
	}

	// Создадим node задачи.
	taskProps := neo4j.Properties{
		"summary": summary,
		"description": description,
	}
	taskNode, err := nodes.Create(taskProps)
	if err != nil {
		rev.ERROR.Fatalln("Failed to create task")
		panic(err)
	}
	for _, taskRel := range tasksRels {
		tasksNode, err := taskRel.End()
		if err != nil {
			rev.ERROR.Fatalf("Could not get Tasks Node")
			panic(err)
		}
		_, err = taskNode.Relate("IS_TASK", tasksNode.Id(), neo4j.Properties{})
		if err != nil {
			rev.ERROR.Fatalf("Could not relate node %s to tasks node", taskNode)
			panic(err)
		}
	}
	return c.Redirect(Tasks.Index)
}

type Task struct {
	Id int
	Summary string
	Description string
}

type TaskList []Task

func (c Tasks) List() rev.Result {
	// Пока тут будет заглушка.
	return c.RenderJson(TaskList{
		Task{1, "Test", "Wow"},
		Task{2, "Test1", "Wowzaaa"},
    })
}