package routes

import (
	"mini-ess/cmd/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CheckInRoute(route fiber.Router, db *sqlx.DB) {
	route.Post("/attendance/checkin", handlers.CheckInHandler(db))
	route.Post("/attendance/checkout", handlers.CheckOutHandler(db))
}
