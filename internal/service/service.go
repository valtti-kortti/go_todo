package service

import (
	"encoding/json"
	"strconv"

	"go_todo/internal/api/dto"
	"go_todo/internal/repo"
	"go_todo/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
	GetTaskById(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
}

type service struct {
	log  *zap.SugaredLogger
	repo repo.Repository
}

func NewService(repo repo.Repository, log *zap.SugaredLogger) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}

	idTask := s.repo.CreateNewTask(&task)

	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"id_task": idTask},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTasks(ctx *fiber.Ctx) error {
	tasks := s.repo.GetAllTasks()
	response := dto.Response{
		Status: "success",
		Data:   tasks,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTaskById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	if id < 0 {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}
	task, ok := s.repo.GetTaskById(id)
	if !ok {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.NotFound(ctx, "task not found")
	}

	response := dto.Response{
		Status: "success",
		Data:   task,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	if id < 0 {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	ok := s.repo.DeleteTaskById(id)
	if !ok {
		return dto.NotFound(ctx, "task not found")
	}

	response := dto.Response{
		Status: "success",
		Data:   "task deleted",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	if id < 0 {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	var req TaskUpdate

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	newTask := s.repo.UpdateTaskById(id, task)
	if newTask == nil {
		return dto.NotFound(ctx, "task not found")
	}
	response := dto.Response{
		Status: "success",
		Data:   newTask,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
