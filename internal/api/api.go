package api

import (
	"github.com/1F47E/go-shaihulud/internal/core"
	"github.com/go-playground/validator/v10"
	zlog "github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type API struct {
	app      *fiber.App
	core     *core.Core
	validate *validator.Validate
}

func NewApi(core *core.Core) *API {
	// add recovery middleware
	app := fiber.New(fiber.Config{})

	app.Use(recover.New())
	app.Use(logger.New())

	api := &API{
		app:      app,
		core:     core,
		validate: validator.New(),
	}

	// app.Static("/", "./webui")
	app.Static("/", "./webui/react/chat/build")
	// app.Static("/chat", "./webui/chat")
	app.Get("/ping", api.Ping)
	app.Post("/chat/create", api.ChatCreate)
	app.Post("/chat/join", api.ChatJoin)
	app.Post("/chat/send", api.ChatSend)
	return api
}

// start server
func (a *API) Start(endpoint string) {
	zlog.Info().Msg("Starting server on " + endpoint)
	if err := a.app.Listen(endpoint); err != nil {
		zlog.Fatal().Err(err).Msg("Error starting server")
	}
}

// shutdown server
func (a *API) Shutdown() error {
	return a.app.Shutdown()
}

// ===== Handlers

func (a *API) Ping(c *fiber.Ctx) error {
	// body := c.Body()
	zlog.Info().Msg("ping")
	return c.SendString("pong")
}

// start the server
func (a *API) ChatCreate(c *fiber.Ctx) error {
	auth, err := a.core.Client.RunServer("")
	if err != nil {
		zlog.Error().Err(err).Msg("Error starting chat server")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(auth)
}

// join the server
type ChatJoinRequest struct {
	Key      string `json:"key" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *API) ChatJoin(c *fiber.Ctx) error {
	var req ChatJoinRequest
	if err := c.BodyParser(&req); err != nil {
		zlog.Error().Err(err).Msg("Error parsing join request")
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing join request")
	}

	if err := a.validate.Struct(req); err != nil {
		zlog.Error().Err(err).Msg("Validation error")
		return c.Status(fiber.StatusBadRequest).SendString("Validation error")
	}

	err := a.core.Client.AuthVerify(req.Key, req.Password)
	if err != nil {
		zlog.Error().Err(err).Msg("Error verifying auth")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = a.core.Client.RunClient()
	if err != nil {
		zlog.Error().Err(err).Msg("Error starting chat client")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

// send message
type ChatSendRequest struct {
	Message string `json:"message" validate:"required"`
}

func (a *API) ChatSend(c *fiber.Ctx) error {
	var req ChatSendRequest
	if err := c.BodyParser(&req); err != nil {
		zlog.Error().Err(err).Msg("Error parsing send request")
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing send request")
	}

	if err := a.validate.Struct(req); err != nil {
		zlog.Error().Err(err).Msg("Validation error")
		return c.Status(fiber.StatusBadRequest).SendString("Validation error")
	}

	err := a.core.Client.Send(req.Message)
	if err != nil {
		zlog.Error().Err(err).Msg("Error sending message")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
