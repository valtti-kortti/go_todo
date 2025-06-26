package api

import (
	"go_todo/internal/api/middleware"
	"go_todo/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	// Настройка CORS
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PATCH, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	// Группа маршрутов
	apiGroup := app.Group("/v1", middleware.Authorization(token))
	apiGroup.Post("/tasks", r.Service.CreateTask)
	apiGroup.Get("/tasks", r.Service.GetTasks)
	apiGroup.Get("/tasks/:id", r.Service.GetTaskById)
	apiGroup.Delete("/tasks/:id", r.Service.DeleteTask)
	apiGroup.Patch("/tasks/:id", r.Service.UpdateTask)

	return app
}
