package main

import (
	"github.com/jmcvetta/neo4j"
	"log"
)

func main() {
	neo, err := neo4j.Connect("http://localhost:7474/db/data")
	log.Println("Connecting")
	if err != nil {
		log.Fatalln("Could not connect to database")
		return
	}
	log.Println("Checking index: users")

	nodes := neo.Nodes
	indexes := nodes.Indexes

	if _, err := indexes.Get("users"); err == nil {
		log.Println("Index  users already exists")
	} else {
		log.Println("Create index: users")
		usersIndex, err := indexes.Create("users")
		if err != nil {
			log.Fatalln("Could not create index users: ", err)
		}

		rootNode, err := nodes.Get(0)

		if err != nil {
			log.Fatalln("Failed to get root node")
		}

		log.Println("Create users node")
		usersProps := neo4j.Properties{
			"type": "users",
		}
		usersNode, err := nodes.Create(usersProps)

		if err != nil {
			log.Fatalln("Could not create users node")
		}

		log.Println("Create relation from root node to users node")

		_, err = rootNode.Relate("USERS", usersNode.Id(), neo4j.Properties{})
		if err != nil {
			log.Fatalln("Could not relate users node to root node")
		}

		login := "root"
		rootUserProps := neo4j.Properties{
			"login":    login,
			"password": "root",
		}

		log.Println("Create root user")
		rootUser, err := nodes.Create(rootUserProps)
		if err != nil {
			log.Fatalln("Could not relate root user node")
		}

		err = usersIndex.Add(rootUser, "login", login)

		log.Println("Indexing root user")
		if err != nil {
			log.Fatalln("Could not add user root to index users")
		}

		log.Println("Create relation from users node to root user node")
		_, err = rootUser.Relate("IS", usersNode.Id(), neo4j.Properties{})
		if err != nil {
			log.Fatalln("Could not relate root user to users")
		}
	}
}
