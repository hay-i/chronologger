package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionAnswer struct {
	QuestionID primitive.ObjectID `bson:"question_id"`
	Answer     string             `bson:"answer"`
}

type Answer struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	TemplateID      primitive.ObjectID `bson:"template_id"`
	Username        string             `bson:"username"`
	CreatedAt       time.Time          `bson:"created_at"`
	QuestionAnswers []QuestionAnswer   `bson:"questions,omitempty"`
}

func (a Answer) FromForm(templateID primitive.ObjectID, username string, questionViews []QuestionView) Answer {
	questionAnswers := make([]QuestionAnswer, len(questionViews))

	for i, q := range questionViews {
		questionAnswers[i] = QuestionAnswer{
			QuestionID: q.ID,
			Answer:     q.Value,
		}
	}

	a.TemplateID = templateID
	a.Username = username
	a.CreatedAt = time.Now()
	a.QuestionAnswers = questionAnswers

	return a
}
