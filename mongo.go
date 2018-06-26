package utils

import(
    "github.com/araddon/dateparse"
    "strconv"
    "strings"
    "reflect"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type MongoField struct {
    Name string
    Type string
}

func MongoFields(fields []string)(map[string]int){

    //Fields
    mongoFields := make(map[string]int)
    if len(fields) > 0 {
        for _, element := range fields{
    		mongoFields[element] = 1
    	}
    }

    return mongoFields
}

func MongoParamsToQuery(coll *mgo.Collection, params Params, model interface{}) (*mgo.Query){

    //Formatted Query
    queries := MongoFormatQuery(params.Query, model)

    //Find All
    query := coll.
    Find(queries)

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
    }

    return query
}

func MongoFormatQuery(query map[string]interface{}, model interface{}) (map[string][]map[string]interface{}){

    //Var
    var queries []map[string]interface{}
    var searchFields []string

    //BSON Tag Map
    bsonStructName := make(map[string]MongoField)
    t := reflect.TypeOf(model)
    for i := 0; i < t.NumField(); i++ {

        //BSON Tag
        tag, _ := t.Field(i).Tag.Lookup("bson")

        //Search
        search, _ := t.Field(i).Tag.Lookup("search")
        if search == "true" {
            searchFields = append(searchFields, tag)
        }
        bsonStructName[tag] = MongoField{
            Name : t.Field(i).Name,
            Type : t.Field(i).Type.Elem().Name(),
        }
    }

    //Loop Through Query
    for key, val := range query {

        //Check Comparison
        var comparison string
        comparisonKey := key
        if strings.Contains(comparisonKey, "_greater_than") {
            comparisonKey = strings.Replace(comparisonKey, "_greater_than", "", 1)
            comparison = "$gt"
        }
        if strings.Contains(comparisonKey, "_greater_than_or_equal_to") {
            comparisonKey = strings.Replace(comparisonKey, "_greater_than_or_equal_to", "", 1)
            comparison = "$gte"
        }
        if strings.Contains(comparisonKey, "_less_than") {
            comparisonKey = strings.Replace(comparisonKey, "_less_than", "", 1)
            comparison = "$lt"
        }
        if strings.Contains(comparisonKey, "_less_than_or_equal_to") {
            comparisonKey = strings.Replace(comparisonKey, "_less_than_or_equal_to", "", 1)
            comparison = "$lte"
        }

        //Check Query
        if key == "q" {

            if len(searchFields) > 0 {

                var orQueries []map[string]interface{}

                for _, searchField := range searchFields {
                    orQueries = append(orQueries, map[string]interface{}{
                        searchField : bson.RegEx{val.(string), "i"},
                    })
                }

                queries = append(queries, map[string]interface{}{
                    "$or" : orQueries,
                })

            }
            delete(query, key)
        }

        if field, ok := bsonStructName[comparisonKey]; ok {

            //Get Val Type
            valType := reflect.TypeOf(val)

            //Time
            if field.Type == "Time" {
                var convertedVal time.Time

                //String
                if valType.String() == "string" {
                    timeVal, err := dateparse.ParseAny(val.(string))
                    if err == nil {
                        convertedVal = timeVal
                    }
                }

                //Not Zero
                if !convertedVal.IsZero(){
                    //Comparison
                    if comparison != "" {
                        queries = append(queries, map[string]interface{}{
                            comparisonKey : map[string]interface{}{
                                comparison : convertedVal,
                            },
                        })
                    } else {
                        queries = append(queries, map[string]interface{}{
                            comparisonKey : convertedVal,
                        })
                    }

                    //Delete
                    delete(query, key)

                }

            }

            //Int64
            if field.Type == "int64" {

                var convertedVal int64

                //String
                if valType.String() == "string" {
                    int64Val, err := strconv.ParseInt(val.(string), 10, 64)
                    if err == nil {
                        convertedVal = int64Val
                    }
                }

                //Int64
                if valType.String() == "int64" {
                    convertedVal = val.(int64)
                }

                //Comparison
                if comparison != "" {
                    queries = append(queries, map[string]interface{}{
                        comparisonKey : map[string]interface{}{
                            comparison : convertedVal,
                        },
                    })
                } else {
                    queries = append(queries, map[string]interface{}{
                        comparisonKey : convertedVal,
                    })
                }

                //Delete
                delete(query, key)
            }




            //Bool
            if field.Type == "bool" {

                //String
                if valType.String() == "string" {
                    boolVal, err := strconv.ParseBool(val.(string))
                    if err == nil {
                        queries = append(queries, map[string]interface{}{
                            comparisonKey : boolVal,
                        })
                    }
                }

                //Bool
                if valType.String() == "bool" {
                    queries = append(queries, map[string]interface{}{
                        comparisonKey : val.(bool),
                    })
                }

                //Delete
                delete(query, key)

            }

        }

    }

    //All Others
    for key, val := range query {
        queries = append(queries, map[string]interface{}{
            key : val,
        })
    }

    var and map[string][]map[string]interface{}
    if len(queries) > 0 {
        and = map[string][]map[string]interface{}{
            "$and" : queries,
        }
    }

    return and
}
