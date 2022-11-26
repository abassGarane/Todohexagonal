package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abassGarane/todos/domain"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisClient(redisURL string, ctx context.Context) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
func NewRedisRepository(redisURL string, ctx context.Context) (domain.TodoRepository, error) {
	repo := &redisRepository{}
	client, err := NewRedisClient(redisURL, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}
func (r *redisRepository) generateKey(id string) string {
	return fmt.Sprintf("todo:%s", id)
}
func (r *redisRepository) Find(id string) (*domain.Todo, error) {
	todo := &domain.Todo{}
	key := r.generateKey(id)
	data, err := r.client.HGetAll(r.client.Context(), key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(domain.ErrorTodoNotFound, "repository.Todo.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Find")
	}
	todo.Content = data["content"]
	todo.Status = data["status"]
	todo.ID = data["id"]
	todo.CreatedAt = createdAt
	return todo, nil
}

func (r *redisRepository) Update(todo *domain.Todo, id string) (*domain.Todo, error) {
	td, err := r.Find(id)
	if err != nil {
		if err == domain.ErrorTodoNotFound {
			return nil, errors.Wrap(domain.ErrorTodoNotFound, "repository.Todo.Update")
		}
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	td.Content = todo.Content
	td.Status = todo.Status
	key := r.generateKey(id)
	_, err = r.client.HSet(r.client.Context(), key, td).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return td, nil
}
func (r *redisRepository) Delete(id string) error {
	_, err := r.client.Del(r.client.Context(), r.generateKey(id)).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Delete")
	}
	return nil
}

func (r *redisRepository) Add(todo *domain.Todo) error {
	key := r.generateKey(todo.ID)
	data := map[string]interface{}{
		"content":    todo.Content,
		"id":         todo.ID,
		"status":     todo.Status,
		"created_at": todo.CreatedAt,
	}
	_, err := r.client.HMSet(r.client.Context(), key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Add")
	}
	return nil
}
