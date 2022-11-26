package domain

type TodoService interface {
	Add(todo *Todo) error
	Update(todo *Todo, id string) (*Todo, error)
	Delete(id string) error
	Find(id string) (*Todo, error)
}
