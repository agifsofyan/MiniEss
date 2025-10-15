package middlewares

import (
	"fmt"
	"mini-ess/configs"
	"mini-ess/internal/schemas"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

func AuthMiddleware(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			auth = c.Cookies("access_token")
			if auth == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
			}
		} else {
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
			}
			auth = parts[1]
		}

		token, err := jwt.ParseWithClaims(auth, &configs.TokenClaim{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return configs.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(*configs.TokenClaim)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
		}

		uid := int64(claims.Sub)
		var u schemas.User
		if err := db.Get(&u, `
			SELECT u.id, u.employee_id, u.email, u.name, u.role, 
				CASE
					WHEN EXISTS (SELECT 1 FROM attendance a WHERE a.user_id  = u.id) THEN 1
					ELSE 0
				END AS is_check_in
			FROM users u 
			WHERE u.id = ?
		`, uid); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}

		c.Locals("user", &u)
		c.Locals("user_id", u.ID)
		c.Locals("is_check_in", u.IsCheckIn)

		return c.Next()
	}
}

func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("Sorry, can't find that!")
}
