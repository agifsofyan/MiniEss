package handlers

import (
	"database/sql"
	"time"

	"mini-ess/configs"
	"mini-ess/internal/schemas"
	"mini-ess/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

func LoginHandler(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		var user schemas.User

		err := db.Get(&user, "SELECT id, employee_id, email, name, password_hash, timezone, role FROM users WHERE email = ?", req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(401).JSON(fiber.Map{"error": "invalid"})
			}
			return err
		}

		isMatch := utils.Compare(user.Password, req.Password)
		if !isMatch {
			return c.Status(400).JSON(fiber.Map{"error": "password not match"})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  user.ID,
			"role": user.Role,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
		s, _ := token.SignedString(configs.JwtSecret)
		return c.JSON(fiber.Map{"access_token": s})
	}
}
