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
		var isLoggedIn bool
		cookie, err := c.Cookie("token")
		if err != nil {
			isLoggedIn = false
		} else {
			_, err = parseToken(cookie.Value)

			if err != nil {
				isLoggedIn = false
			} else {
				isLoggedIn = true
			}
		}

		component := components.Home(isLoggedIn)

		return render(c, component)
	}
}

func Templates(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		templates := db.GetDefaultTemplates(requestContext, database)
		component := components.Templates(templates)

		return render(c, component)
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

		return render(c, component)
	}
}

func Modal(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		id := c.Param("id")
		template := db.GetTemplate(requestContext, database, id)
		component := components.Modal(template)

		return render(c, component)
	}
}

func Start(database *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestContext := c.Request().Context()
		id := c.Param("id")
		template := db.GetTemplate(requestContext, database, id)
		component := components.Start(template)

		return render(c, component)
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

		return redirect(c, "/templates/"+templateId+"?success=true", http.StatusFound)
	}
}
