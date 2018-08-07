package utils

import(
    "strconv"
    "strings"
    "github.com/araddon/dateparse"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type UpdatePart struct {
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


        timeVal, err := dateparse.ParseAny(val)
        if err == nil {
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

    queryMap := make(map[string]interface{})

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
    for key, value := range queryStringParameters {

		//Check for Array
		if strings.Index(value, ",") > 0 {
			queryArray := strings.Split(value, ",")
			queryMap[key] = queryArray
		} else {
			queryMap[key] = value
		}
    }
    //params.Query = queryMap

    return params, err
}

func ParseUpdateMongo(updates []UpdatePart)(mgo.Change){

    set := make(map[string]interface{})

    for _, update := range updates {
        set[update.Field] = update.NewValue
    }
    change := mgo.Change{
        Update: map[string]interface{}{"$set" : set},
        ReturnNew: true,
    }

    return change
}

func ParseParamsMongo(coll *mgo.Collection, params Params)(*mgo.Query){

    //BSON Query
    q := make(map[string]interface{})

    for _, queryPart := range params.Query {
        q[queryPart.Field] = bson.M{
            queryPart.Operator : queryPart.Value,
        }
    }

    //Find All
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
