package utils

import(
    "strconv"
    "strings"
)

type Params struct {
    Query map[string]interface{}
    Fields []string
    Sort string
    Page int
    Limit int
    Skip int
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
    params.Query = queryMap

    return params, err
}
