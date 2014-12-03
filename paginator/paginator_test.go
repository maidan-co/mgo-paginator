package paginator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	//"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	//"reflect"
	"math/rand"
	"time"
	"fmt"
	"reflect"
	"labix.org/v2/mgo/bson"
)

var fakedExecuteAndBind = func(query *mgo.Query ,resultSet interface {}) (error) {
	quarterOfSociety := people[offset : size + offset]
	resultv := reflect.Indirect(reflect.ValueOf(resultSet))
	resultv.Set(reflect.ValueOf(quarterOfSociety))
	fmt.Println("Ahanda: ", len(quarterOfSociety))
	return nil
}

var originalExecuteAndBind = executeAndBind
var firstObjectId = bson.NewObjectId()
var paginator *Paginator
var people = []Person{}
var offset,size int

func init() {
	/*to be able to use collection.Find() inside paginator we need a dummy collection object*/
	session, _ := mgo.Dial("localhost")
	collection := session.DB("test").C("dummy")
	paginator = &Paginator{collection, Person{}}
	genesis()
	offset = 0
	size = 23
}

type Person struct {
	Name string
	Age int
	Id bson.ObjectId
}

func genesis() {
	names := []string{"zeus", "eve", "adam"}
	for i := 0; i < 92; i++ {
		randName := names[i % len(names)]
		objectId := firstObjectId
		rand.Seed(time.Now().UTC().UnixNano())
		newBorn := Person{randName, rand.Intn(100), objectId}
		people = append(people, newBorn)
		objectId = bson.NewObjectId()
	}
}

func TestPaginator(t *testing.T) {
	Convey("Subject: Paginator Unit Tests", t, func() {
			Convey("When called with correct parameters", func() {
					Convey("When called with only Paginate without after or before", func() {
							originalExecuteAndBind = executeAndBind
							executeAndBind = fakedExecuteAndBind
							result := []Person{}
							beforeId, afterId, err := paginator.Paginate("","").Execute(&result)
							if err != nil {
								fmt.Println("Error on test: ", err)
							}

							Convey("So it should not provide before link but after", func() {
									So(beforeId, ShouldEqual, "")
									So(afterId, ShouldEqual, result[len(result) - 1].Id.Hex())
								})

							Convey("So number of elements inside of result set should be less than or equal to maximum limit", func() {
									So(len(result), ShouldBeLessThanOrEqualTo, 23)//23 is the default limit
								})

							Convey("Result set should include data from the beginning of data set", func() {
									firstHuman := result[0]
									So(firstHuman.Id, ShouldEqual, firstObjectId)
								})
							executeAndBind = originalExecuteAndBind
						})

					Convey("When called with only Paginate within after", func() {
							originalExecuteAndBind = executeAndBind
							executeAndBind = fakedExecuteAndBind
							offset = 1
							size = 23
							result := []Person{}
							afterIdStr := firstObjectId.Hex()
							beforeId, afterId, err := paginator.Paginate(afterIdStr,"").Execute(&result)
							if err != nil {
								fmt.Println("Error on test: ", err)
							}

							Convey("So it should provide before and after links", func() {
									So(beforeId, ShouldEqual, result[0].Id.Hex())
									So(afterId, ShouldEqual, result[len(result) - 1].Id.Hex())
								})

							Convey("So number of elements inside of result set should be less than or equal to maximum limit", func() {
									So(len(result), ShouldBeLessThanOrEqualTo, 23)//23 is the default limit
								})

							Convey("Result set should include data that comes after than the provided object ID (current cursor)", func() {
									firstElem := result[0]
									So(firstElem.Id, ShouldEqual, people[1].Id)//because we wanted to retrieve people right after the 0th person
								})

						})

					Convey("When called with only Paginate within before", func() {
							originalExecuteAndBind = executeAndBind
							executeAndBind = fakedExecuteAndBind
							offset = 23
							size = 23
							result := []Person{}
							beforeIdStr := people[46].Id.Hex()
							beforeId, afterId, err := paginator.Paginate("",beforeIdStr).Execute(&result)
							if err != nil {
								fmt.Println("Error on test: ", err)
							}

							Convey("So it should provide before and after links", func() {
									So(beforeId, ShouldEqual, result[0].Id.Hex())
									So(afterId, ShouldEqual, result[len(result) - 1].Id.Hex())
								})

							Convey("So number of elements inside of result set should be less than or equal to maximum limit", func() {
									So(len(result), ShouldBeLessThanOrEqualTo, 23)//23 is the default limit
								})

							Convey("Result set should include data that comes before than the provided object ID (current cursor)", func() {
									firstElem := result[0]
									So(firstElem.Id, ShouldEqual, people[23].Id)//because we wanted to retrieve people right before the 1th person
								})

						})

					Convey("When called with Limit also", func() {
							originalExecuteAndBind = executeAndBind
							executeAndBind = fakedExecuteAndBind
							offset = 5
							size = 5
							result := []Person{}
							beforeId, afterId, err := paginator.Paginate("","").Limit(5).Execute(&result)
							if err != nil {
								fmt.Println("Error on test: ", err)
							}

							Convey("So it should provide before and after links", func() {
									So(beforeId, ShouldEqual, result[0].Id.Hex())
									So(afterId, ShouldEqual, result[len(result) - 1].Id.Hex())
								})

							Convey("So number of elements inside of result set should be less than or equal to given limit number", func() {
									So(len(result), ShouldBeLessThanOrEqualTo, 5)
								})

						})

					Convey("When called with Sort criteria", func() {

							Convey("So it should provide before and after links", nil)

							Convey("So result set should be sorted correctly", nil)

						})

					Convey("When called with Filter criteria", func() {

							Convey("So it should provide before and after links", nil)

							Convey("So related items should be retrieved regarding criteria", nil)

							Convey("So others should be filtered regarding criteria", nil)

						})

				})

			Convey("When called with wrong params", func() {

					Convey("When called with after and before same time", func() {

							Convey("So it should return an error", nil)

						})

				})

		})}
