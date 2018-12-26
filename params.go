package utils

import(
    "strconv"
    "strings"
    "github.com/araddon/dateparse"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/jinzhu/gorm"
)

type UpdatePart struct {
    Query []QueryPart
    PullQuery []QueryPart
    Operator string
    Label string
    Field string
    PreviousValue interface{}
    NewValue interface{}
}

type QueryPart struct {
    Field string
    Operator string
    Value interface{}
}

type Params struct {
    Query []QueryPart
    Fields []string
    Sort string
    Page int
    Limit int
    Skip int
}

func QueryValuesToParams(urlMap map[string][]string) (Params, error){

    //Params
    var params Params
    var err error

    //Token
    if(len(urlMap["token"]) > 0){
        delete(urlMap, "token")
    }

    //Fields
    if(len(urlMap["fields"]) > 0){
        params.Fields = strings.Split(urlMap["fields"][0], ",")
        delete(urlMap, "fields")
    }

    //Sort
    if(len(urlMap["sort"]) > 0){
        params.Sort = urlMap["sort"][0]
        delete(urlMap, "sort")
    }

    //Page
	params.Page = 1
    if(len(urlMap["page"]) > 0){
        page := urlMap["page"][0]
        i, _ := strconv.Atoi(page)
        params.Page = i
        delete(urlMap, "page")
    }

    //Limit
    if(len(urlMap["limit"]) > 0){
        limit := urlMap["limit"][0]
        x, _ := strconv.Atoi(limit)
        params.Limit = x
        delete(urlMap, "limit")
    }

    //Skip
    if(len(urlMap["skip"]) > 0){
        skip := urlMap["skip"][0]
        x, _ := strconv.Atoi(skip)
        params.Skip = x
        delete(urlMap, "skip")
    }

    //Query
    for field, value := range urlMap {

        var valFormatted interface{}
        operator := "$eq"

        val := value[0]

        if strings.Contains(field, "_greater_than") {
            field = strings.Replace(field, "_greater_than", "", 1)
            operator = "$gt"
        } else if strings.Contains(field, "_greater_than_or_equal_to") {
            field = strings.Replace(field, "_greater_than_or_equal_to", "", 1)
            operator = "$gte"
        } else if strings.Contains(field, "_less_than") {
            field = strings.Replace(field, "_less_than", "", 1)
            operator = "$lt"
        } else if strings.Contains(field, "_less_than_or_equal_to") {
            field = strings.Replace(field, "_less_than_or_equal_to", "", 1)
            operator = "$lte"
        }

        intVal, intValErr := strconv.Atoi(val)
        timeVal, timeValErr := dateparse.ParseAny(val)

        if intValErr == nil {
            valFormatted = intVal
        } else if timeValErr == nil {
            valFormatted = timeVal
        } else if val == "true" {
            valFormatted = true
        } else if val == "false" {
            valFormatted = false
        } else if strings.Index(val, ",") > 0 {
			valFormatted = strings.Split(val, ",")
            operator = "$in"
		} else {
			valFormatted = val
		}

        params.Query = append(params.Query, QueryPart{
            Field : field,
            Operator : operator,
            Value : valFormatted,
        })
    }

    return params, err
}

func QueryStringParametersToParams(queryStringParameters map[string]string) (Params, error){

    //Params
    var params Params
    var err error

    //queryMap := make(map[string]interface{})

    //Token
    if _, ok := queryStringParameters["token"]; ok {
        delete(queryStringParameters, "token")
    }

    //Fields
    if val, ok := queryStringParameters["fields"]; ok {
        params.Fields = strings.Split(val, ",")
        delete(queryStringParameters, "fields")
    }

    //Sort
    if val, ok := queryStringParameters["sort"]; ok {
        params.Sort = val
        delete(queryStringParameters, "sort")
    }

    //Page
    params.Page = 1
    if val, ok := queryStringParameters["page"]; ok {
	    x, err := strconv.Atoi(val)
        if err != nil {
            return params, err
        }
        params.Page = x
        delete(queryStringParameters, "page")
    }

    //Limit
    if val, ok := queryStringParameters["limit"]; ok {
        x, err := strconv.Atoi(val)
        if err != nil {
            return params, err
        }
        params.Limit = x
        delete(queryStringParameters, "limit")
    }

    //Skip
    if val, ok := queryStringParameters["skip"]; ok {
        x, err := strconv.Atoi(val)
        if err != nil {
            return params, err
        }
        params.Skip = x
        delete(queryStringParameters, "skip")
    }

    //Query
    for field, val := range queryStringParameters {

        var valFormatted interface{}
        operator := "$eq"


        if strings.Contains(field, "_greater_than") {
            field = strings.Replace(field, "_greater_than", "", 1)
            operator = "$gt"
        } else if strings.Contains(field, "_greater_than_or_equal_to") {
            field = strings.Replace(field, "_greater_than_or_equal_to", "", 1)
            operator = "$gte"
        } else if strings.Contains(field, "_less_than") {
            field = strings.Replace(field, "_less_than", "", 1)
            operator = "$lt"
        } else if strings.Contains(field, "_less_than_or_equal_to") {
            field = strings.Replace(field, "_less_than_or_equal_to", "", 1)
            operator = "$lte"
        }


        intVal, intValErr := strconv.Atoi(val)
        timeVal, timeValErr := dateparse.ParseAny(val)

        if intValErr == nil {
            valFormatted = intVal
        } else if timeValErr == nil {
            valFormatted = timeVal
        } else if val == "true" {
            valFormatted = true
        } else if val == "false" {
            valFormatted = false
        } else if strings.Index(val, ",") > 0 {
			valFormatted = strings.Split(val, ",")
            operator = "$in"
		} else {
			valFormatted = val
		}

        params.Query = append(params.Query, QueryPart{
            Field : field,
            Operator : operator,
            Value : valFormatted,
        })

    }

    return params, err
}

func ParseUpdateMongo(updates []UpdatePart)(mgo.Change){

    //Vars
    update := make(map[string]interface{})
    set := make(map[string]interface{})
    push := make(map[string]interface{})
    pull := make(map[string]interface{})
    for _, update := range updates {

        //Set
        if update.Operator == "$set" || update.Operator == "" {
            set[update.Field] = update.NewValue
        }

        //Push
        if update.Operator == "$push"{
            push[update.Field] = update.NewValue
        }

        //Pull
        if update.Operator == "$pull"{

            //BSON Query
            if len(update.PullQuery) > 0 {
                q := make(map[string]interface{})
                for _, queryPart := range update.PullQuery {
                    q[queryPart.Field] = bson.M{
                        queryPart.Operator : queryPart.Value,
                    }
                }
                pull[update.Field] = q
            }

        }

    }

    //Update
    if len(set) > 0 {
        update["$set"] = set
    }
    if len(push) > 0 {
        update["$push"] = push
    }
    if len(pull) > 0 {
        update["$pull"] = pull
    }
    change := mgo.Change{
        Update: update,
        ReturnNew: true,
    }

    return change
}

func ParseParamsGorm(db *gorm.DB, params Params)(*gorm.DB){

    //Query
    var queryFields []string
    var queryVals []interface{}
    for _, queryPart := range params.Query {

        //$eq
        if queryPart.Operator == "$eq" {
            queryFields = append(queryFields, queryPart.Field + " = ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$ne
        if queryPart.Operator == "$ne" {
            queryFields = append(queryFields, queryPart.Field + " <> ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$in
        if queryPart.Operator == "$in" {
            queryFields = append(queryFields, queryPart.Field + " IN(?)")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$lt
        if queryPart.Operator == "$lt" {
            queryFields = append(queryFields, queryPart.Field + " < ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$lte
        if queryPart.Operator == "$lte" {
            queryFields = append(queryFields, queryPart.Field + " <= ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$gt
        if queryPart.Operator == "$gt" {
            queryFields = append(queryFields, queryPart.Field + " > ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$gte
        if queryPart.Operator == "$gte" {
            queryFields = append(queryFields, queryPart.Field + " >= ?")
            queryVals = append(queryVals, queryPart.Value)
        }

        //$regex
        if queryPart.Operator == "$regex" {
            queryFields = append(queryFields, queryPart.Field + " LIKE ?")
            queryVals = append(queryVals, "%" + queryPart.Value.(string) + "%")
        }
    }

    //Build
    if len(queryFields) > 0 && len(queryVals) > 0 && len(queryFields) == len(queryVals){
        db = db.Where(strings.Join(queryFields, " AND "), queryVals...)
    }

    //Fields
    if len(params.Fields) > 0 {
        db = db.Select(params.Fields)
    }

    //Sort
    if params.Sort != "" {
        if strings.Contains(params.Sort, "-"){
            db = db.Order(strings.Replace(params.Sort, "-", "", 1) + " DESC")
        } else {
            db = db.Order(params.Sort + " ASC")
        }
    }

    //Limit
    if params.Limit != 0 {
        db = db.Limit(params.Limit)
    }

    //Skip
    if params.Skip != 0 {
        db = db.Offset(params.Skip)
    } else {

        //Page
        if params.Page != 0 && params.Limit != 0 {
            db = db.Offset(params.Limit * (params.Page - 1))
        }

    }

    return db
}

func ParseParamsMongo(coll *mgo.Collection, params Params)(*mgo.Query){

    //BSON Query
    q := make(map[string]interface{})

    for _, queryPart := range params.Query {
        q[queryPart.Field] = bson.M{
            queryPart.Operator : queryPart.Value,
        }
    }

    //Query
    query := coll.
    Find(q)

    //Fields
    fields := make(map[string]int)
    if len(params.Fields) > 0 {
        for _, element := range params.Fields{
    		fields[element] = 1
    	}
    }
    query.Select(fields)

    //Sort
    if params.Sort != "" {
        query.Sort(params.Sort)
    }

    //Limit
    if params.Limit != 0 {
        query.Limit(params.Limit)
    }

    //Skip
    if params.Skip != 0 {
        query.Skip(params.Skip)
    } else {

        //Page
        if params.Page != 0 && params.Limit != 0 {
            query.Skip(params.Limit * (params.Page - 1))
        }

    }



    return query
}
