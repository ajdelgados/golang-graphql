package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ajdelgados/golang-graphql/models"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

type reqBody struct {
	Query string `json:"query"`
}

func Graphql(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	rawBody, _ := c.GetRawData()
	b := bytes.NewBuffer(rawBody)

	var rBody reqBody
	err := json.NewDecoder(b).Decode(&rBody)
	if err != nil {
		fmt.Println(err)
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"todo": &graphql.Field{
				Type: models.TodoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var find models.Todo
					if id, isOK := params.Args["id"].(int); isOK == true {
						find.ID = id
					}
					if name, isOK := params.Args["name"].(string); isOK == true {
						find.Name = name
					}
					if status, isOK := params.Args["status"].(int); isOK == true {
						find.Status = status
					}

					var todo models.Todo
					if err := db.Where(find).First(&todo).Error; err != nil {
						return nil, err
					}

					return todo, nil
				},
			},
			"todos": &graphql.Field{
				Type: graphql.NewList(models.TodoType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var find models.Todo
					if id, isOK := params.Args["id"].(int); isOK == true {
						find.ID = id
					}
					if name, isOK := params.Args["name"].(string); isOK == true {
						find.Name = name
					}
					if status, isOK := params.Args["status"].(int); isOK == true {
						find.Status = status
					}

					var todos []models.Todo
					db.Where(find).Find(&todos)

					return todos, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type: models.TodoType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var todo models.Todo
					todo.Name = params.Args["name"].(string)
					todo.Status = 1
					db.Create(&todo)

					return todo, nil
				},
			},
			"updateTodo": &graphql.Field{
				Type: models.TodoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)

					var todo models.Todo
					if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
						return nil, err
					}

					if name, isOK := params.Args["name"].(string); isOK == true {
						todo.Name = name
					}

					if status, isOK := params.Args["status"].(int); isOK == true {
						todo.Status = status
					}

					db.Save(&todo)

					return todo, nil
				},
			},
			"deleteTodo": &graphql.Field{
				Type: models.TodoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)

					var todo models.Todo
					if err := db.Where("id = ?", id).First(&todo).Delete(&todo).Error; err != nil {
						return nil, err
					}

					return todo, nil
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	schema, _ := graphql.NewSchema(schemaConfig) // Query

	params := graphql.Params{Schema: schema, RequestString: rBody.Query}
	r := graphql.Do(params)

	c.JSON(http.StatusOK, r)
}
