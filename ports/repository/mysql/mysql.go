package mysql

import (
	"context"
	"database/sql"
	"log"

	"github.com/abassGarane/todos/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type mysqlRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewMysqlRepository(dsn string, ctx context.Context) (domain.TodoRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("Unable to start mysql database", err)
		return nil, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "repository.mysql.NewRepository")
	}
	return &mysqlRepository{
		db:  db,
		ctx: ctx,
	}, nil
}

func (m *mysqlRepository) Add(todo *domain.Todo) error {
	stmt := `insert into todos(id, content, status, created_at)values (?, ?, ?)`
	_, err := m.db.ExecContext(m.ctx, stmt, todo.ID, todo.Content, todo.Status, todo.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Add")
	}
	return nil

}
func (m *mysqlRepository) Update(todo *domain.Todo, id string) (*domain.Todo, error) {
	stmt := `UPDATE todos set content=?, status=? WHERE todo.id=?`
	_, err := m.db.ExecContext(m.ctx, stmt, todo.Content, todo.Status, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return todo, nil
}
func (m *mysqlRepository) Delete(id string) error {
	stmt := `DELETE FROM todos WHERE id=?`
	_, err := m.db.ExecContext(m.ctx, stmt, id)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Delete")
	}
	return nil
}
func (m *mysqlRepository) Find(id string) (*domain.Todo, error) {
	stmt := `SELECT * FROM todos WHERE id=?`
	todo := &domain.Todo{}
	if err := m.db.QueryRow(stmt, id).Scan(&todo); err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return todo, nil
}
