package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func getGraphqlType(u interface{}) *graphql.Object {
	reflectedType := reflect.TypeOf(u)
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   reflectedType.Name(),
			Fields: getGraphqlFields(u),
		},
	)
}

func getGraphqlFields(t interface{}) graphql.Fields {
	var fieldMap = make(map[string]*graphql.Field)

	reflectedType := reflect.TypeOf(t)
	reflectedValue := reflect.ValueOf(t)

	for i := 0; i < reflectedType.NumField(); i++ {
		valueField := reflectedValue.Field(i)
		typeField := reflectedType.Field(i)

		switch valueField.Kind() {
		case reflect.String:
			fieldMap[typeField.Name] = &graphql.Field{
				Type:        graphql.String,
				Description: typeField.Tag.Get("doc"),
			}
		case reflect.Int:
			fieldMap[typeField.Name] = &graphql.Field{
				Type:        graphql.Int,
				Description: typeField.Tag.Get("doc"),
			}
		case reflect.Struct:
			fieldMap[typeField.Name] = &graphql.Field{
				Type:        getGraphqlType(valueField.Interface()),
				Description: typeField.Tag.Get("doc"),
			}
		}
	}
	return fieldMap
}

type ContactType struct {
	ID    int    `json:"id" doc:"The id of the Contact"`
	Email string `json:"email" doc:"The email of the Contact"`
}

type UserType struct {
	ID      int         `json:"id" doc:"The id of the User"`
	Name    string      `json:"name" doc:"The Name of the User"`
	Contact ContactType `json:"contact" doc:"The Contact associated with the User"`
}

type QueryType struct {
	User UserType `json:"user" doc:"Query to fetch all users"`
}

var q = QueryType{}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: getGraphqlType(q),
	},
)

func main() {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Println("Starting server at port 8080")
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
