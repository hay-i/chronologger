package controllers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hay-i/chronologger/components"
	"github.com/hay-i/chronologger/db"
	"github.com/hay-i/chronologger/models"
	"github.com/labstack/echo/v4"
)

func Home(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		component := components.Home()

		return component.Render(requestContext, c.Response().Writer)
	}
}

func Templates(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		templates := db.GetDefaultTemplates(requestContext, database)
		component := components.Templates(templates)

		return component.Render(requestContext, c.Response().Writer)
	}
}

func DismissModal() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func Template(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		id := c.Param("id")
		template := db.GetTemplate(requestContext, database, id)
		answers := db.GetAnswers(requestContext, database, id)
		success := c.QueryParam("success")
		component := components.Template(template, answers, success == "true")

		return component.Render(requestContext, c.Response().Writer)
	}
}

func Modal(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		id := c.Param("id")
		template := db.GetTemplate(requestContext, database, id)
		component := components.Modal(template)

		return component.Render(requestContext, c.Response().Writer)
	}
}

func Start(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		id := c.Param("id")
		template := db.GetTemplate(requestContext, database, id)
		component := components.Start(template)

		return component.Render(requestContext, c.Response().Writer)
	}
}

func Response(database *mongo.Database, client *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		answersCollection := database.Collection("answers")
		templateId := c.Param("id")
		req := c.Request()
		requestContext := req.Context()

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		session, err := client.StartSession()
		if err != nil {
			panic(err)
		}
		defer session.EndSession(requestContext)
		err = session.StartTransaction()
		if err != nil {
			panic(err)
		}
		defer session.AbortTransaction(requestContext)
		defer session.CommitTransaction(requestContext)

		templateObjectId, err := primitive.ObjectIDFromHex(templateId)
		if err != nil {
			panic(err)
		}

		for questionId, value := range req.Form {
			questionObjectId, err := primitive.ObjectIDFromHex(questionId)
			if err != nil {
				panic(err)
			}
			_, insertErr := answersCollection.InsertOne(requestContext, models.Answer{
				TemplateID: templateObjectId,
				QuestionID: questionObjectId,
				Answer:     value[0],
			})

			if insertErr != nil {
				panic(insertErr)
			}
		}

		// Ideally I wanted to send in http.StatusCreated, but it seems that the redirect doesn't work with that status code
		// See: https://github.com/labstack/echo/issues/229#issuecomment-1518502318
		return c.Redirect(http.StatusFound, "/templates/"+templateId+"?success=true")
	}
}
