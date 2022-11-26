package rest

import (
	"fmt"

	"github.com/abassGarane/todos/domain"
	"github.com/abassGarane/todos/ports/serializers/json"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type TodoHandler interface {
	Get(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type todoHandler struct {
	service domain.TodoService
}

// // Configures the response header and writes into the response
//
//	func setUpResponse(c *fiber.Ctx, contentType string, statusCode int, body []byte) error {
//		c.Set("Content-Type", contentType)
//		c.Status(statusCode)
//		_, err := c.Write(body)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func NewTodoHandler(todoService domain.TodoService) TodoHandler {
	return &todoHandler{
		service: todoService,
	}
}

func (h todoHandler) Post(c *fiber.Ctx) error {
	// contentType := c.GetReqHeaders()["Content-Type"]
	body := c.Body()
	if len(body) == 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Body Must be non-null",
		})
	}
	// Temporarily using json only serializer
	t := &json.Todo{}
	todo, err := t.Decode(body)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}
	err = h.service.Add(todo)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(todo)
}
func (h todoHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	todo, err := h.service.Find(id)
	if err != nil {
		if errors.Cause(err) == domain.ErrorTodoNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("Todo with id %s Not found", id),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

func (h todoHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	body := c.Body()
	td, err := h.service.Find(id)
	if err != nil {
		if errors.Cause(err) == domain.ErrorTodoNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("Todo with id %s Not found", id),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}
	if len(body) == 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Body Must be non-null",
		})
	}
	// Temporarily using json only serializer
	t := &json.Todo{}
	todo, err := t.Decode(body)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}

	td.Content = todo.Content
	td.Status = todo.Status
	if err = h.service.Add(td); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"todo":    td,
		"message": "Successfuly updated",
	})
}

func (h *todoHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		if errors.Cause(err) == domain.ErrorTodoNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("Todo with id %s Not found", id),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
