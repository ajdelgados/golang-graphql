package models

import (
	"github.com/graphql-go/graphql"
)

type Todo struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

var TodoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
