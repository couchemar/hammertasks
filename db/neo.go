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

func (db *DataBase) GetTask(id int) (*models.Task, error) {
	nodes := db.Connection.Nodes
	taskNode, err := nodes.Get(id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get id: %s (%s)", id, err))
	}
	sum, err := taskNode.GetProperty("summary")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get 'summary' for id: %s (%s)", id, err))
	}
	desc, err := taskNode.GetProperty("description")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get 'description' for id: %s (%s)", id, err))
	}
	task := models.Task{
		Id:          taskNode.Id(),
		Summary:     sum,
		Description: desc,
	}
	return &task, nil
}

func (db *DataBase) GetTasksList() *models.TaskList {

	var result cypherResponse
	var nerr interface{}

	query := cypherRequest{
		Query:  "START r=node(0) MATCH r-[:TASKS]->t<-[:IS_TASK]-tasks RETURN tasks",
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

		task, err := db.GetTask(id)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, *task)
	}

	return &tasks

}
