package main

import (
	"fmt"
	"mini-ess/cmd/api/middlewares"
	"mini-ess/cmd/api/routes"
	"mini-ess/configs"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
)

func main() {
	cmsPort := configs.GetEnv("API_PORT")

	app := fiber.New(fiber.Config{
		AppName:       configs.GetEnv("APP_NAME"),
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ReadTimeout:   15 * time.Second,
		WriteTimeout:  7 * time.Second,
		IdleTimeout:   30 * time.Second,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	db, err := configs.Connection()
	if err != nil {
		log.Errorf("DB Connection failed: %s", err.Error())
	}
	defer db.Close()

	// Routes
	openRoute(app, db)
	app.Use(middlewares.AuthMiddleware(db))
	restrictRoute(app, db)

	app.Use(middlewares.NotFoundMiddleware)

	addr := strings.ReplaceAll(fmt.Sprintf(":%s", cmsPort), " ", "")
	log.Fatal(app.Listen(addr))
}

func openRoute(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api").Group("/v1")
	routes.LoginRoute(api, db)
	routes.CheckHealth(api)
	routes.CreateEmployee(api, db)
}

func restrictRoute(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api").Group("/v1")
	routes.CheckInRoute(api, db)
}
