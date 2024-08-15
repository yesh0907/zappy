package routes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"zappy.sh/database"
	"zappy.sh/models"
	"zappy.sh/repositories"
)

type Handler struct {
	AliasRepository   *repositories.AliasRepository
	RequestRepository *repositories.RequestRepository
}

func NewHandler() *Handler {
	return &Handler{
		AliasRepository:   repositories.NewAliasRepository(database.DBConn),
		RequestRepository: repositories.NewRequestRepository(database.DBConn),
	}
}

func (h *Handler) GetAlias(c *fiber.Ctx) error {
	name := strings.ToLower(c.Params("alias"))
	userId := strings.ToLower(c.Query("id"))
	source := strings.ToLower(c.Query("source"))

	req := new(models.Request)

	req.AliasName = name
	req.IP = fmt.Sprintf("ip: %s, x-forwarded-by: %s, fly-client-ip: %s", c.IP(), c.Get("X-Forwarded-For"), c.Get("Fly-Client-IP"))
	req.UserAgent = c.Get("User-Agent")

	if userId != "" {
		req.UserId = userId
	}

	if source != "" {
		req.Referer = source
	}

	// Set headers to not cache response
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	// Set headers to expire immediately
	c.Set("Expires", "0")

	alias, err := h.AliasRepository.GetAlias(name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   "alias not found",
				"data":    nil,
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
		}
	}

	// Log request
	if err := h.RequestRepository.CreateRequest(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"data":    nil,
		})
	}

	// Redirect to the alias url
	return c.Redirect(alias.Url, 301)
}

func (h *Handler) CreateAlias(c *fiber.Ctx) error {
	alias := new(models.Alias)

	if err := c.BodyParser(alias); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"created": false,
			"error":   err.Error(),
		})
	}

	if err := h.AliasRepository.CreateAlias(alias); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(200).JSON(fiber.Map{
				"created": false,
				"error":   nil,
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"created": false,
				"error":   err.Error(),
			})
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"created": true,
		"error":   nil,
	})
}

func (h *Handler) AllRequests(c *fiber.Ctx) error {

	aliasName := strings.ToLower(c.Params("alias"))

	if aliasName == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "alias name is required",
			"data":    nil,
		})
	}

	requests, err := h.RequestRepository.GetAllRequests(aliasName)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   "alias not found",
				"data":    nil,
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"error":   nil,
		"data":    requests,
	})
}
