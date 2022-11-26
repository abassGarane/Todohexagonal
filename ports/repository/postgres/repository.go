package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/abassGarane/todos/domain"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	ctx context.Context
	db  *sql.DB
}

func NewPostgresRepository(dsn string, ctx context.Context) (domain.TodoRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic("Unable to start mysql database", err)
		return nil, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "repository.postgres.NewRepository")
	}
	return &postgresRepository{
		db:  db,
		ctx: ctx,
	}, nil
}

func (m *postgresRepository) Add(todo *domain.Todo) error {
	stmt := `insert into todos(id, content, status, created_at)values ($1, $1, $3)`
	_, err := m.db.ExecContext(m.ctx, stmt, todo.ID, todo.Content, todo.Status, todo.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Add")
	}
	return nil

}
func (m *postgresRepository) Update(todo *domain.Todo, id string) (*domain.Todo, error) {
	stmt := `UPDATE todos set content=$1, status=$2 WHERE todo.id=$3`
	_, err := m.db.ExecContext(m.ctx, stmt, todo.Content, todo.Status, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return todo, nil
}
func (m *postgresRepository) Delete(id string) error {
	stmt := `DELETE FROM todos WHERE id=$1`
	_, err := m.db.ExecContext(m.ctx, stmt, id)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Delete")
	}
	return nil
}
func (m *postgresRepository) Find(id string) (*domain.Todo, error) {
	stmt := `SELECT * FROM todos WHERE id=$1`
	todo := &domain.Todo{}
	if err := m.db.QueryRow(stmt, id).Scan(&todo); err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return todo, nil
}
