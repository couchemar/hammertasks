package db

import (
	"github.com/jmcvetta/neo4j"
	"hammertasks/app/models"
)

type DataBase struct {
	Connection *neo4j.Database
}

func Connect(host string) *DataBase {
	connection, err := neo4j.Connect(host)
	if err != nil {
		panic(err)
	}
	db := DataBase{connection}

	return &db
}

func (db DataBase) GetRootNode() *neo4j.Node {
	nodes := db.Connection.Nodes
	rootNode, err := nodes.Get(0)
	if err != nil {
		panic(err)
	}
	return rootNode
}

func (db DataBase) GetTasksList() *models.TaskList {
	rootNode := db.GetRootNode()

	tasksRels, err := rootNode.Outgoing("TASKS")
	if err != nil {
		panic(err)
	}

	var tasksNode *neo4j.Node
	for _, taskRel := range tasksRels {
		tasksNode, err = taskRel.End()
		if err != nil {
			panic(err)
		}
	}

	relsToTasks, err := tasksNode.Incoming("IS_TASK")
	if err != nil {
		panic(err)
	}

	tasks := make(models.TaskList, 1)

	for _, rel := range relsToTasks {
		taskNode, err := rel.Start()
		if err != nil {
			panic(err)
		}
		sum, err := taskNode.GetProperty("summary")
		if err != nil {
			panic(err)
		}
		desc, err := taskNode.GetProperty("description")
		if err != nil {
			panic(err)
		}

		task := models.Task{
			Id:          taskNode.Id(),
			Summary:     sum,
			Description: desc,
		}
		tasks = append(tasks, task)
	}

	return &tasks
}
