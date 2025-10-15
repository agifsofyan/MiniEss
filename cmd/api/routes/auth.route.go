package routes

import (
	"mini-ess/cmd/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func LoginRoute(route fiber.Router, db *sqlx.DB) {
	route.Post("/login", handlers.LoginHandler(db))
}
