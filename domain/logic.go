package domain

import (
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrorTodoNotFound = errors.New("Todo Not Found")
	ErrorInvalidTodo  = errors.New("Invalid Todo")
)

type todoService struct {
	todoRepository TodoRepository
}

func NewTodoService(todoRepository TodoRepository) TodoService {
	return &todoService{
		todoRepository: todoRepository,
	}
}
func (t *todoService) Add(todo *Todo) error {
	if err := validate.Validate(todo); err != nil {
		return errors.Wrap(ErrorInvalidTodo, "service.Todo.Add")
	}
	todo.ID = shortid.MustGenerate()
	todo.CreatedAt = time.Now().UTC().Unix()
	return t.todoRepository.Add(todo)
}

func (t *todoService) Update(todo *Todo, id string) (*Todo, error) {
	return t.todoRepository.Update(todo, id)
}

func (t *todoService) Delete(id string) error {
	return t.todoRepository.Delete(id)
}

func (t *todoService) Find(id string) (*Todo, error) {
	return t.todoRepository.Find(id)
}
