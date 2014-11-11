package paginator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
	//"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	//"reflect"
	"labix.org/v2/mgo/bson"
	base "gitlab.salamworld.dev/framework/base"
	"gitlab.salamworld.dev/platform/discussionapi/models"
)

var Collection *mgo.Collection
var fakedExecuteAndBind = func(query *mgo.Query ,resultSet interface {}) (error) {
	return nil
}
var originalExecuteAndBind = executeAndBind

func init() {
	base.CreateSession("db.salamworld.dev:27017,db.salamworld.dev:27018,db.salamworld.dev:27019")
	Collection = base.GetCollection("maidan", "discussions")
	fmt.Print("Collection: " + Collection.FullName)

}

func TestPaginator(t *testing.T) {
	Convey("Subject: GivePaginatedResult Unit Test", t, func() {

			Convey("When called with correct parameters", func() {

					Convey("When called without sort and find criterias", func() {
							//resultModel := []models.Discussion{}
							paginator := Paginator{Collection, models.Discussion{}}
							//executeAndBind = fakedExecuteAndBind
							//objectId := bson.ObjectIdHex("544a24724484120012000007")
							//err := paginator.GivePaginatedResult(true, objectId, 3, &resultModel)
							//So(err, ShouldEqual, nil)
							arr,_ , _, _:= paginator.Paginate("544a24724484120012000007", "").Filter(bson.M{"type":"discussion","staffPicked":false}).Limit(10).Execute()
							fmt.Println("\nResult: ",arr.([]models.Discussion)[9])
							executeAndBind = originalExecuteAndBind
						})

					Convey("When called with sort and find criterias", nil)

				})

			Convey("When called", nil)

		})
}
