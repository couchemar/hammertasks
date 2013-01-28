package db

import (
	"errors"
	"fmt"
	"github.com/jmcvetta/neo4j"
	"github.com/jmcvetta/restclient"
	"hammertasks/app/models"
	"net/url"
	"strconv"
	"strings"
)

var NotFound = errors.New("Not Found")

type DataBase struct {
	url        *url.URL
	Connection *neo4j.Database
}

func Connect(uri string) *DataBase {
	url, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	connection, err := neo4j.Connect(uri)
	if err != nil {
		panic(err)
	}
	db := DataBase{url, connection}

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

type queryParams map[string]interface{}
type cypherRequest struct {
	Query  string      `json:"query"`
	Params queryParams `json:"params"`
}

type nodeResponse struct {
	HrefSelf string `json:"self"`
}

type cypherResponse struct {
	Columns []string         `json:"columns"`
	Data    [][]nodeResponse `json:"data"`
}

type myNode struct {
	neo4j.Node
}

func (taskNode *myNode) getInDeps() (*models.TaskList, error) {
	inRels, err := taskNode.Incoming("DEPENDS_ON")
	if err != nil {
		return nil, err
	}
	nodeList := make(models.TaskList, 0)
	for _, rel := range inRels {
		node, err := rel.Start()
		newNode := myNode{*node}
		task, err := newNode.toModel()
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, task)
	}
	return &nodeList, nil
}

func (taskNode *myNode) getOutDeps() (*models.TaskList, error) {
	outRels, err := taskNode.Outgoing("DEPENDS_ON")
	if err != nil {
		return nil, err
	}
	nodeList := make(models.TaskList, 0)
	for _, rel := range outRels {
		node, err := rel.End()
		newNode := myNode{*node}
		task, err := newNode.toModel()
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, task)
	}
	return &nodeList, nil
}

func (taskNode *myNode) toModel() (*models.Task, error) {
	nodeType, err := taskNode.GetProperty("type")
	if err != nil || nodeType != "task" {
		return nil, NotFound
	}
	sum, err := taskNode.GetProperty("summary")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get '%s' for id: %s (%s)", "summary", taskNode.Id(), err))
	}
	desc, err := taskNode.GetProperty("description")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get '%s' for id: %s (%s)", "description", taskNode.Id(), err))
	}
	task := models.Task{
		Id:          taskNode.Id(),
		Summary:     sum,
		Description: desc,
	}
	return &task, nil
}

func (db *DataBase) GetTask(id int, dependencies bool) (*models.Task, error) {
	nodes := db.Connection.Nodes
	taskNode, err := nodes.Get(id)
	if err != nil {
		return nil, NotFound
	}
	newTaskNode := myNode{*taskNode}

	task, err := newTaskNode.toModel()
	if err != nil {
		return nil, err
	}

	if dependencies == true {
		inDeps, err := newTaskNode.getInDeps()
		if err != nil {
			return nil, err
		}
		task.InDeps = inDeps
		outDeps, err := newTaskNode.getOutDeps()
		if err != nil {
			return nil, err
		}
		task.OutDeps = outDeps
	}

	return task, nil
}

func (db *DataBase) GetTasksList() *models.TaskList {

	var result cypherResponse
	var nerr interface{}

	query := cypherRequest{
		Query:  "START r=node(0) MATCH r-[:TASKS]->t<-[:IS_TASK]-task RETURN task",
		Params: queryParams{},
	}

	url := db.url.String() + "/cypher"
	r := restclient.RestRequest{
		Url:    url,
		Method: restclient.POST,
		Data:   &query,
		Result: &result,
		Error:  &nerr,
	}

	client := restclient.New()
	_, err := client.Do(&r)
	if err != nil {
		panic(err)
	}

	tasks := make(models.TaskList, 0)

	responseData := result.Data

	for _, row := range responseData {
		nodeInfo := row[0]
		self := nodeInfo.HrefSelf
		selfM := strings.Split(self, "/")
		id, err := strconv.Atoi(selfM[len(selfM)-1])
		if err != nil {
			panic(err)
		}

		task, err := db.GetTask(id, false)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Could not get task (id: %s) (%s)", id, err)))
		}
		tasks = append(tasks, task)
	}

	return &tasks

}
