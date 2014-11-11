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
paginator.paginate("after=505050345").filter(bson.M).limit("50").sort("name").execute()


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
var SortOrder int
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
		SortOrder = -1
		rangeFilter = bson.M{"_id": bson.M{"$lt":findId}}
	} else {
		SortOrder = 1
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
		fmt.Println("Paginate didn't called correctly. Only call it with after or before ID.")
		return paginatedQuery
	}
	if afterId != "" {
		SortOrder = -1
		findId := bson.ObjectIdHex(afterId)
		rangeFilter = bson.M{"_id": bson.M{"$lt":findId}}
	} else if beforeId != "" {
		SortOrder = 1
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

func(paginatedQuery *paginatedQuery) Execute() (resultSet interface {}, beforeId string, afterId string, err error) {
	modelType := reflect.TypeOf(paginatedQuery.resultModel)
	sliceTypeOfModal := reflect.SliceOf(modelType)
	//modelSlice := reflect.MakeSlice(sliceTypeOfModal, 0, 0).Interface()

	v := reflect.New(sliceTypeOfModal)
	//v.Elem().Set(reflect.ValueOf(modelSlice))
	resultModelArray := v.Interface()// we resolved sth. like: *[]modelObject

	paginatedQuery.query = paginatedQuery.collection.Find(paginatedQuery.rangeFilter)

	if paginatedQuery.limitValue > 0 {
		paginatedQuery.query = paginatedQuery.query.Limit(paginatedQuery.limitValue)
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

	err = executeAndBind(paginatedQuery.query, resultModelArray)
	if err != nil {
		return nil, "", "", err
	}
	actualReflectArray := reflect.ValueOf(resultModelArray).Elem()

	if actualReflectArray.Len() > 1 {
		lastElem := actualReflectArray.Index(actualReflectArray.Len() - 1)
		firstElem := actualReflectArray.Index(0)
		beforeId = firstElem.FieldByName(DefaultObjectIdFieldName).Interface().(bson.ObjectId).Hex()
		afterId = lastElem.FieldByName(DefaultObjectIdFieldName).Interface().(bson.ObjectId).Hex()
	}
	//firstElem :=
	return actualReflectArray.Interface(), beforeId, afterId, nil//we resolved []modelObject
}


