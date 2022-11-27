package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/abassGarane/todos/domain"
	"github.com/abassGarane/todos/ports/api/rest"
	"github.com/abassGarane/todos/ports/repository/mongodb"
	"github.com/abassGarane/todos/ports/repository/mysql"
	"github.com/abassGarane/todos/ports/repository/postgres"
	"github.com/abassGarane/todos/ports/repository/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func httpPort() string {
	if os.Getenv("PORT") == "" {
		return ":8000"
	}
	return fmt.Sprintf(":%s", os.Getenv("PORT"))
}

func chooseRepository(ctx context.Context) domain.TodoRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisURL, ctx)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mysql":
		mysqlURL := os.Getenv("MYSQL_URL")
		repo, err := mysql.NewMysqlRepository(mysqlURL, ctx)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "postgres":
		postgresURL := os.Getenv("POSTGRES_URL")
		repo, err := postgres.NewPostgresRepository(postgresURL, ctx)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		timeOut, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongodb.NewMongoRepository(mongoURL, mongoDB, ctx, timeOut)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	default:
		log.Fatal("No database choosen")
		return nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repo := chooseRepository(ctx)
	service := domain.NewTodoService(repo)
	handler := rest.NewTodoHandler(service)
	app := fiber.New()
	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())

	//Todo handlers
	app.Get("/:id", handler.Get)
	app.Delete("/:id", handler.Delete)
	app.Patch("/:id", handler.Update)
	app.Post("/", handler.Post)

	//Server controls
	errs := make(chan error, 2)
	go func() {
		fmt.Printf("Running server on http://localhost%s", httpPort())
		errs <- app.Listen(httpPort())
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated server :: %s\n", <-errs)
}
