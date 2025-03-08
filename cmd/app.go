package main

import (
	"api/internal/config"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func testHandler(ctx context.Context, c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func main() {
	cfg, errCfg := config.LoadConfig()
	if errCfg != nil {
		log.Fatal().Msg(errCfg.Error())
		return
	}
	ctx := context.WithValue(context.Background(), "config", cfg)

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Get("/", testHandler(ctx, app))

	err := app.Listen(fmt.Sprintf("%s:%d", (*cfg).Host, (*cfg).Port))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("Listening host %s on port %d", (*cfg).Host, (*cfg).Port)
}
