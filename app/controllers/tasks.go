package controllers

import (
	"encoding/json"
	"github.com/jmcvetta/neo4j"
	"github.com/robfig/revel"
	"hammertasks/app/models"
	"hammertasks/db"
)

type Tasks struct {
	Application
}

func (c Tasks) Index() rev.Result {
	return c.Render()
}

func (c Tasks) CreateTask() rev.Result {

	requestDecoder := json.NewDecoder(c.Request.Body)
	var task models.Task
	err := requestDecoder.Decode(&task)

	if err != nil {
		rev.ERROR.Printf("Decode error: %s", err)
		panic(err)
	}

	summary := task.Summary
	description := task.Description

	neo := db.Connect("http://localhost:7474/db/data")
	nodes := neo.Connection.Nodes
	rootNode := neo.GetRootNode()
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
		"summary":     summary,
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
	task.Id = taskNode.Id()
	return c.RenderJson(task)
}

func (c Tasks) DetailPage() rev.Result {
	return c.Render()
}

func (c Tasks) GetTask(id int) rev.Result {
	neo := db.Connect("http://localhost:7474/db/data")
	task, err := neo.GetTask(id)
	if err != nil {
		panic(err)
	}
	return c.RenderJson(task)
}

func (c Tasks) List() rev.Result {
	neo := db.Connect("http://localhost:7474/db/data")
	tasks := neo.GetTasksList()
	return c.RenderJson(tasks)
}

func (c Tasks) ListPage() rev.Result {
	return c.Render()
}
