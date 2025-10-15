package handlers

import (
	"mini-ess/internal/schemas"
	"mini-ess/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CheckInHandler(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*schemas.User)

		var req struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			DeviceID  string  `json:"device_id"`
			Note      string  `json:"note"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
		}

		if utils.IsWeekendWIB(time.Now()) {
			return c.Status(400).JSON(fiber.Map{"error": "Check in is only available on weekdays (Monday - Friday) Asia/Jakarta timezone"})
		}

		endedCheckInTime := utils.IsAfterOrEqualTimeInWIB(time.Now(), 9, 1)
		if endedCheckInTime {
			return c.Status(400).JSON(fiber.Map{"error": "Check in time has ended"})
		}

		if hasCheckIn(db, user.ID) {
			return c.Status(400).JSON(fiber.Map{"error": "You have checked in today"})
		}

		res, err := db.Exec(`
        INSERT INTO attendance (user_id, latitude, longitude, status, note)
        VALUES (?, ?, ?, 'completed', ?)`,
			user.ID, req.Latitude, req.Longitude, req.Note,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to insert attendance"})
		}

		attID, _ := res.LastInsertId()

		return c.Status(201).JSON(fiber.Map{
			"message":       "Check-in success",
			"attendance_id": attID,
			"status":        "completed",
		})
	}
}

func CheckOutHandler(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*schemas.User)

		var req struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			DeviceID  string  `json:"device_id"`
			Note      string  `json:"note"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
		}

		if !user.IsCheckIn {
			return c.Status(400).JSON(fiber.Map{"error": "You haven't checked in today"})
		}

		isAvailableCheckOutTime := utils.IsAfterOrEqualTimeInWIB(time.Now(), 17, 0)
		if !isAvailableCheckOutTime {
			return c.Status(400).JSON(fiber.Map{"error": "You can only check out at 17:00 WIB (Asia/Jakarta)"})
		}

		if hasCheckOut(db, user.ID) {
			return c.Status(400).JSON(fiber.Map{"error": "You have checkout in today"})
		}

		res, err := db.Exec(`
        UPDATE attendance SET check_out_at = ?
        WHERE user_id = ? AND DATE(check_in_at) = CURDATE()`,
			utils.CreateTimeIdn(time.Now()), user.ID,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to insert attendance"})
		}

		attID, _ := res.LastInsertId()

		return c.Status(201).JSON(fiber.Map{
			"message":       "Check-out success",
			"attendance_id": attID,
			"status":        "completed",
		})
	}
}

func hasCheck(db *sqlx.DB, userId int64, kind string) bool {
	type CheckData struct {
		Total int64 `db:"total"`
	}

	var checkData CheckData

	var checkOutValue string
	if kind == "in" {
		checkOutValue = "IS NULL"
	} else {
		checkOutValue = "IS NOT NULL"
	}

	err := db.Get(&checkData, "SELECT COUNT(*) AS total FROM attendance WHERE user_id = ? AND DATE(check_in_at) = CURDATE() AND check_out_at "+checkOutValue, userId)
	if err != nil {
		return false
	}

	return checkData.Total > 0
}

func hasCheckIn(db *sqlx.DB, userId int64) bool {
	return hasCheck(db, userId, "in")
}

func hasCheckOut(db *sqlx.DB, userId int64) bool {
	return hasCheck(db, userId, "out")
}
