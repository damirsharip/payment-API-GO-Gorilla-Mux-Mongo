package models

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)



type Payment struct {
	ID    		primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Name 		string            	`json:"name" bson:"name"`
	Price		int					`json:"price" bson:"price"`
	Date   		time.Time			`json:"date," bson:"date"`
	Typeof  	string				`json:"typeof" bson:"typeof"`
	Comment 	string				`json:"comment" bson:"comment"`
	Category 	string				`json:"category" bson:"category"`
}	