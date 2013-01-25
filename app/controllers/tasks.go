package controllers

import (
	"encoding/json"
	"github.com/jmcvetta/neo4j"
	"github.com/robfig/revel"
	"hammertasks/app/models"
	"hammertasks/db"
	"net/http"
)

type Tasks struct {
	Application
}

func (c Tasks) Index() rev.Result {
	return c.Render()
}

func (c Tasks) decodeTask() (*models.Task, error) {
	requestDecoder := json.NewDecoder(c.Request.Body)
	var task models.Task
	err := requestDecoder.Decode(&task)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (c Tasks) CreateTask() rev.Result {
	task, err := c.decodeTask()

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
		"type":        "task",
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

func (c Tasks) UpdateTask(id int) rev.Result {

	var err error
	task, err := c.decodeTask()

	if err != nil {
		rev.ERROR.Printf("Decode error: %s", err)
		panic(err)
	}

	summary := task.Summary
	description := task.Description

	neo := db.Connect("http://localhost:7474/db/data")
	nodes := neo.Connection.Nodes

	taskNode, err := nodes.Get(id)
	if err != nil {
		rev.ERROR.Printf("Could not get task: %s", err)
		panic(err)
	}

	err = taskNode.SetProperty("summary", summary)
	if err != nil {
		rev.ERROR.Printf("Could not set propery 'summary' : %s", err)
		panic(err)
	}
	err = taskNode.SetProperty("description", description)
	if err != nil {
		rev.ERROR.Printf("Could not set propery 'description' : %s", err)
		panic(err)
	}

	return c.RenderJson(task)
}

func (c Tasks) EditPage() rev.Result {
	return c.Render()
}

type errorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c Tasks) GetTask(id int) rev.Result {
	neo := db.Connect("http://localhost:7474/db/data")
	task, err := neo.GetTask(id)
	if err != nil {
		if err == db.NotFound {
			c.Response.Status = http.StatusNotFound
			return c.RenderJson(errorJSON{
				Code:    404,
				Message: "Not found",
			})
		} else {
			panic(err)
		}
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

func (c Tasks) DeleteTask(id int) rev.Result {
	neo := db.Connect("http://localhost:7474/db/data")
	nodes := neo.Connection.Nodes

	taskNode, err := nodes.Get(id)
	if err != nil {
		rev.ERROR.Printf("Could not get task: %s", err)
		panic(err)
	}

	rels, err := taskNode.Relationships()
	if err != nil {
		rev.ERROR.Printf("Could not get task relationships: %s", err)
		panic(err)
	}
	for _, rel := range rels {
		err = rel.Delete()
		if err != nil {
			rev.ERROR.Printf("Could not delete task relation: %s", err)
			panic(err)
		}
	}

	err = taskNode.Delete()
	if err != nil {
		rev.ERROR.Printf("Could not delete task: %s", err)
		panic(err)
	}
	return c.RenderText("ok")
}
