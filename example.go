package main

import (
	"fmt"
	"labix.org/v2/mgo"
	//"github.com/maidan-co/mgo-paginator/paginator"
	"labix.org/v2/mgo/bson"
)

type Community struct {
	Id                interface{} `bson:"_id" json:"_id,omitempty"`
	Owner             mgo.DBRef   `bson:"owner" json:"owner,omitempty"`
	Title             string      `valid:"Required;MaxSize(60)" bson:"title" json:"title,omitempty"`
	Cover             string      `bson:"cover" json:"cover,omitempty"`
	Description       string      `valid:"Required;MaxSize(5000)" bson:"description" json:"description,omitempty"`
	RulesOfEngagement string      `bson:"rulesOfEngagement" json:"rulesOfEngagement,omitempty"`
	StaffPicked       bool        `bson:"staffPicked" json:"staffPicked,omitempty"`
	Popular           bool        `bson:"popular" json:"popular,omitempty"`
	Sections          []string    `bson:"sections" json:"sections,omitempty"`
	Members           []mgo.DBRef `bson:"members" json:"members,omitempty"`
	MembersCount      int         `bson:"-" json:"members_count,omitempty"`
	Type              string      `json:"type,omitempty"`
	DiscussionCount   int         `bson:"discussion_count" json:"discussionCount,omitempty"`
	Discussions       interface{} `json:"discussions"`
	IsMember          bool        `bson:"-" json:"isMember,omitempty"`
}

func main() {
	session, _ := mgo.Dial("localhost:27017")
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB("maidan").C("communities")
/*
	examplePaginator := paginator.Paginator{collection, Community{}}
	resultSet, beforeId, afterId, err := examplePaginator.Paginate("53f338bc62c77407d0000001", "").Limit(10).Execute()
	fmt.Println(resultSet, beforeId, afterId, err)
*/
	tags :=[2]string{"fdsf","samsung"}
	result := []Community{}
	collection.Find(bson.M{"_id": bson.ObjectIdHex("53f3224b83f01b0550000001"), "tags":bson.M{"$in": tags}}).All(&result)
	fmt.Println(result)
}
