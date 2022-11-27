package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/abassGarane/todos/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	timeout  time.Duration
	database string
	ctx      context.Context
}

func newClient(mongoUrl string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Print("Unable to disconnect client", err)
		}
	}()
	return client, nil
}

func NewMongoRepository(mongoUrl, mongoDB string, ctx context.Context, mongoTimeout int) (domain.TodoRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
		ctx:      ctx,
	}
	client, err := newClient(mongoUrl, ctx)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client
	return repo, nil
}

func (m *mongoRepository) Add(todo *domain.Todo) error {
	col := m.client.Database(m.database).Collection("todos")
	bs := bson.M{
		"id":         todo.ID,
		"created_at": todo.CreatedAt,
		"status":     todo.Status,
		"content":    todo.Content,
	}
	_, err := col.InsertOne(m.ctx, bs)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Add")
	}
	return nil
}

func (m *mongoRepository) Update(todo *domain.Todo, id string) (*domain.Todo, error) {
	col := m.client.Database(m.database).Collection("todos")
	filter := bson.M{
		"id": id,
	}
	_, err := col.UpdateOne(m.ctx, filter, &todo)
	if err == mongo.ErrNoDocuments {
		return nil, errors.Wrap(err, "repository.Todo.Update")
	}
	return todo, nil
}

func (m *mongoRepository) Delete(id string) error {
	col := m.client.Database(m.database).Collection("todos")
	filter := bson.M{
		"id": id,
	}
	_, err := col.DeleteOne(m.ctx, filter)
	if err != nil {
		return errors.Wrap(err, "repository.Todo.Delete")
	}
	return nil
}

func (m *mongoRepository) Find(id string) (*domain.Todo, error) {
	todo := &domain.Todo{}
	filter := bson.M{
		"id": id,
	}
	col := m.client.Database(m.database).Collection("todos")
	if err := col.FindOne(m.ctx, filter).Decode(todo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrorTodoNotFound, "repository.Todo.Find")
		}
		return nil, errors.Wrap(err, "repository.Todo.Find")
	}
	return todo, nil
}
