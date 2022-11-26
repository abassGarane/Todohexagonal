package domain

type Todo struct {
	ID        string `json:"id" bson:"id" msgpack:"id"`
	Content   string `json:"content" bson:"content" msgpack:"content" validate:"empty=false"`
	Status    string `json:"status" bson:"status" msgpack:"status"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
