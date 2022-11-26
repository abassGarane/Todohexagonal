package json

import (
	"encoding/json"

	"github.com/abassGarane/todos/domain"
	"github.com/pkg/errors"
)

type Todo struct {
}

func (t *Todo) Decode(input []byte) (*domain.Todo, error) {
	td := &domain.Todo{}
	if err := json.Unmarshal(input, td); err != nil {
		return nil, errors.Wrap(err, "serializers.Todo.Decode")
	}
	return td, nil
}

func (t *Todo) Encode(input *domain.Todo) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializers.Todo.Encode")
	}
	return raw, nil
}
