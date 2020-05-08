# Creando un servidor GraphQL con Golang y PostgreSQL

Crearemos un servidor web con el lenguaje de consulta GraphQL, Golang y en base de datos PostgreSQL. Las librerías o frameworks que estaremos usando son Gin, GORM y graphql-go. Gin es el framework web, GORM es el ORM para la base de datos PostgreSQL y graphql-go es para la resolución de las consultas.

## Creando el proyecto TODO

Lo primero es crear un directorio llamado todo, nos colocamos dentro del directorio creado e iniciamos la gestión de dependencias.

```
mkdir todo
cd todo
go mod init todo
```

Luego obtenemos las dependencias necesarias para el proyecto, en nuestro caso son Gin, GORM y graphql-go.

```
go get github.com/gin-gonic/gin github.com/graphql-go/graphql github.com/jinzhu/gorm
```

Con la línea anterior se descargaran la dependencias, además se incluirán en nuestro archivo go.mod.

## Iniciamos el servidor GraphQL

Creamos el archivo main.go e iniciamos el servicio web.

```
package main
  
import (
        "todo/models"
        "todo/controllers"

        "github.com/gin-gonic/gin"
)

func main() {
        r := gin.Default()

        db := models.SetupModels()
        defer db.Close()

        r.Use(func(c *gin.Context) {
                c.Set("db", db)
                c.Next()
        })

        r.GET("/graphql", controllers.Graphql)

        r.Run()
}
```

En el código anterior, iniciamos el servidor con las configuraciones por defecto de Gin, levantamos la conexión a la base de datos con el método SetupModels que crearemos en la próxima sección, indicamos que al finalizar la función, se ejecute la función para cerrar la conexión a la base de datos.

También creamos un middleware donde establecemos la asociación entre la conexión anteriormente creada con una variable en el servidor web, nos servirá a la hora de hacer el o los controladores. En la línea 20, creamos la ruta o endpoint para la resolución de las peticiones en nuestro controlador.

## Generando el modelo TODO y conexión a la base de datos

Dentro del proyecto crearemos un directorio llamado models y dentro un archivo llamada todo.go, donde especificaremos la estructura del modelo y el type para GraphQL.

```
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
```

En models también creamos la conexión hacia la base de datos en el archivo setup.go, con los parametros necesarios para ingresar al PostgreSQL y una migración de la estructura.

```
package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SetupModels() *gorm.DB {
	db, err := gorm.Open("postgres", "user= dbname= sslmode=disable")

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	return db
}
```

En el directorio models dejamos las definiciones con las estructuras necesarias.

## Controlador para resolver la petición

En está sección desarrollaremos el procesamiento del query o mutation de la petición. Procedemos a hacer un directorio controllers, allí creamos un archivo llamado graphql.go donde tendremos una pequeña estructura para obtener el query de la petición y la función para procesar la petición del cliente.

```
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo/models"

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
```

En la función Graphql lo primero es tomar la conexión hacia la base de datos y procesar body de la petición, definimos los query con la variable rootQuery y los mutation con rootMutation. Configuramos el schema y en la línea 184 hacemos la ejecución de los parametros de la petición y el schema.

La librería graphql-go es bastante versátil y vemos validaciones como NewNonNull para hacer obligatorio algún parametro en la definición del query o mutation.

## Conclusión

En este ejemplo vemos lo rápido y fácil generar un API de tipo GraphQL con Golang, usando unos cuantos módulos de tantos que tiene Golang para hacer este tipo de proyectos.