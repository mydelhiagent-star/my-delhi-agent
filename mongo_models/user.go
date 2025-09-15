package mongo_models

import(
	 "go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
	Area     string `json:"area" bson:"area"`
	Role     string `json:"role" bson:"role"`
}
