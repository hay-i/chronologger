package models

import (
	"net/url"
	"time"

	"github.com/hay-i/gooal/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Template struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"created_at"`
	Username    string             `bson:"username"`
	Questions   []Question         `bson:"questions,omitempty"`
}

func (t Template) FromForm(formValues url.Values) Template {
	t.Title = formValues.Get("title")
	t.Description = formValues.Get("description")
	t.Username = formValues.Get("username")
	t.CreatedAt = time.Now()

	formValues.Del("title")
	formValues.Del("description")
	formValues.Del("username")

	t.Questions = QuestionsFromForm(formValues)

	return t
}

func (t Template) Save(database *mongo.Database) {
	db.Save(database, "templates", t)
}

func GetTemplate(database *mongo.Database, id string) Template {
	return db.Get(database, "templates", id).(Template)
}
