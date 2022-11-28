package domain

import (
	"fmt"
)

type Todo struct {
	ID        string `json:"id" bson:"id" msgpack:"id"`
	Content   string `json:"content" bson:"content" msgpack:"content" validate:"empty=false"`
	Status    string `json:"status" bson:"status" msgpack:"status"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}

func (t Todo) String() string {
	return fmt.Sprintf("Todo{ID: %s, Content: %s, Status: %s, CreatedAt: %d}",
		t.ID, t.Content, t.Status, t.CreatedAt)
}
