package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"github.com/maidan-co/mgo-paginator/paginator"
	"labix.org/v2/mgo/bson"
)

type Person struct {
	Id interface{} "_id"
	Name string
	Phone string
}

func main() {
	session, _ := mgo.Dial("localhost:27017")
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB("maidan").C("communities")

	examplePaginator := paginator.Paginator{collection, Person{}}
	resultSet, beforeId, afterId, err := examplePaginator.Paginate("544a24724484120012000007", "").Filter(bson.M{"type":"discussion","staffPicked":false}).Limit(10).Execute()
	fmt.Println(resultSet, beforeId, afterId, err)
}
