package paginator

import (
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	"fmt"
	"reflect"
)

/* This is the sample usage of paginator:

Paginator {
	Collecton
	TypeOfDomain
}
paginator.giveResultsFromRequest(URI)
paginator.paginate("","").filter(bson.M).limit("50").sort("name").execute(&sampleArray)


after=5050505-limit=50
this.Get("after")

*/

type Paginator struct {
	Collection *mgo.Collection
	ResultModel interface{}
}

type paginatedQuery struct {
	query *mgo.Query
	rangeFilter bson.M
	limitValue int
	sortValue string
	selectValue bson.M
	collection *mgo.Collection
	resultModel interface{}
}


var FindCriterias []bson.M
var SortCriteria bson.M
var Delimiter = "="
var DelimiterForFilter = ","
var DefaultLimit = 23
var DefaultObjectIdFieldName = "Id"

func(this *Paginator) GivePaginatedResult(isAfter bool, findId interface{}, limit int, resultModel interface{}) (error) {
	if limit > 1000 {//to avoid memory leak
		limit = 999
	}
	rangeFilter := bson.M{}
	if isAfter {
		rangeFilter = bson.M{"_id": bson.M{"$lt":findId}}
	} else {
		rangeFilter = bson.M{"_id": bson.M{"$gt":findId}}
	}
	query := this.Collection.Find(rangeFilter).Limit(limit)
	err := executeAndBind(query, resultModel)
	if err != nil {
		return err
	}
	return nil
}

var executeAndBind = func(query *mgo.Query ,resultSet interface {}) (error) {
	err := query.All(resultSet)
	return err
}

func(this *Paginator) Paginate(afterId string, beforeId string) (*paginatedQuery) {
	paginatedQuery := &paginatedQuery{}
	paginatedQuery.resultModel = this.ResultModel
	paginatedQuery.collection = this.Collection
	rangeFilter := bson.M{}
	if afterId != "" && beforeId != "" {
		fmt.Println("Paginate didn't called correctly. Call it only with after ID or before ID.")
		return paginatedQuery
	}
	if afterId != "" {
		findId := bson.ObjectIdHex(afterId)
		rangeFilter = bson.M{"_id": bson.M{"$lt":findId}}
	} else if beforeId != "" {
		findId := bson.ObjectIdHex(beforeId)
		rangeFilter = bson.M{"_id": bson.M{"$gt":findId}}
	}
	//rangeFilter = append(rangeFilter,bson.DocElem{"staffPicked",false})
	//rangeFilter = append(rangeFilter,bson.M{"type":"discussion"})
	//paginatedQuery.Query = this.Collection.Find(rangeFilter)
	paginatedQuery.rangeFilter = rangeFilter
	return paginatedQuery
}

func(paginatedQuery *paginatedQuery) Limit(limit int) (*paginatedQuery) {
	if limit > 1000 {//to avoid memory leak
		limit = 999
	}
	paginatedQuery.limitValue = limit
	return paginatedQuery
}

func(paginatedQuery *paginatedQuery) Sort(sort string) (*paginatedQuery) {
	if sort == "" {
		return paginatedQuery
	}
	paginatedQuery.sortValue = sort
	return paginatedQuery
}

func(paginatedQuery *paginatedQuery) Filter(filter bson.M) (*paginatedQuery) {
	if filter == nil {
		return paginatedQuery
	}
	for k, v := range filter {//only way to append two maps:http://stackoverflow.com/questions/7436864/go-copying-all-elements-of-a-map-into-another
		paginatedQuery.rangeFilter[k] = v
	}
	return paginatedQuery
}

func(paginatedQuery *paginatedQuery) Select(selectMap bson.M) (*paginatedQuery) {
	if paginatedQuery == nil {
		return paginatedQuery
	}
	paginatedQuery.selectValue = selectMap
	return paginatedQuery
}

func(paginatedQuery *paginatedQuery) Execute(resultSet interface {}) (beforeId string, afterId string, err error) {
	paginatedQuery.query = paginatedQuery.collection.Find(paginatedQuery.rangeFilter)
	limit := DefaultLimit
	if paginatedQuery.limitValue > 0 {
		paginatedQuery.query = paginatedQuery.query.Limit(paginatedQuery.limitValue)
		limit = paginatedQuery.limitValue
	} else {
		paginatedQuery.query = paginatedQuery.query.Limit(DefaultLimit)
	}

	if paginatedQuery.sortValue != "" {
		paginatedQuery.query = paginatedQuery.query.Sort(paginatedQuery.sortValue)
	} else {
		paginatedQuery.query = paginatedQuery.query.Sort("-_id")//default order is desc
	}

	if paginatedQuery.selectValue != nil {
		paginatedQuery.query = paginatedQuery.query.Select(paginatedQuery.selectValue)
	}

	err = executeAndBind(paginatedQuery.query, resultSet)
	if err != nil {
		return "", "", err
	}
	actualReflectArray := reflect.ValueOf(resultSet).Elem()

	if actualReflectArray.Len() > 1 {
		lastElem := actualReflectArray.Index(actualReflectArray.Len() - 1)
		firstElem := actualReflectArray.Index(0)
		if paginatedQuery.rangeFilter["_id"] != nil { //if it's not the beginning of the paginated data
			beforeId = firstElem.FieldByName(DefaultObjectIdFieldName).Interface().(bson.ObjectId).Hex()
		}
		if actualReflectArray.Len() == limit { //if it's not the end of the paginated data
			afterId = lastElem.FieldByName(DefaultObjectIdFieldName).Interface().(bson.ObjectId).Hex()
		}
	}
	return beforeId, afterId, nil//we resolved []modelObject
}


