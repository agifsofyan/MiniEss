package handlers

import (
	"mini-ess/internal/schemas"
	"mini-ess/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CheckHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to API Service",
	})
}

func CreateEmployee(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req schemas.CheckReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		userID := c.Locals("user_id").(int64)

		photoURL := ""

		res, err := db.Exec(`INSERT INTO attendance_records
        (user_id,type,occurred_at_utc,latitude,longitude,location_accuracy,photo_url,device_id,source,status)
        VALUES (?, 'checkin', ?, ?, ?, ?, ?, ?, 'mobile', 'pending')`,
			userID, utils.CreateTimeIdn(time.Now()), req.Latitude, req.Longitude, req.Accuracy, photoURL, req.DeviceID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		id, _ := res.LastInsertId()

		return c.Status(201).JSON(fiber.Map{"id": id, "status": "pending"})
	}
}
