package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	ParentID  *primitive.ObjectID `bson:"parent_id,omitempty"`
	CreatedAt primitive.DateTime  `bson:"created_at"`
}
