package routes

import (
	"mini-ess/cmd/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CheckHealth(route fiber.Router) {
	route.Get("/", handlers.CheckHealth)
}

func CreateEmployee(route fiber.Router, db *sqlx.DB) {
	route.Post("/employee", handlers.CreateEmployee(db))
}
