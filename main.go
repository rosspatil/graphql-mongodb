package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {

	nestedObject := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Address",
		Description: "This respresents an Author",
		Fields: graphql.Fields{
			"street": &graphql.Field{
				Type: graphql.String,
			},
			"pincode": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	fields := graphql.Fields{
		"address": &graphql.Field{
			Type: nestedObject,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					Description:  "Input name",
					DefaultValue: "abc",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				return map[string]interface{}{"street": "sample street"}, nil
			},
		},
		"hello": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					Description:  "Input name",
					DefaultValue: "abc",
				},
			},
			Description: "Hello field",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, ok := p.Args["name"]
				if ok {
					return name, nil
				}
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "queryStarter", Fields: fields}

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"addAddress": &graphql.Field{
				Type: nestedObject, // the return type for this field
				Args: graphql.FieldConfigArgument{
					"address": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					fmt.Println(p.Args["address"])
					return map[string]interface{}{"street": "sample street"}, nil
				},
			},
		},
	})
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: rootMutation}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8081", nil)
}
