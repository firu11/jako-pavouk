package main

import (
	"backend/databaze"
	"backend/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

type (
	bodyCreateTrida struct {
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
)

// typy uživatelů
// 1 - basic
// 2 - učitel
func SetupSkolniRouter(api *fiber.Router) {
	skolaApi := (*api).Group("/skola")

	skolaApi.Get("/create-trida", createTrida)
	skolaApi.Get("/tridy", tridy)
}

func createTrida(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tridu muze vytvaret pouze ucitel"))
	}

	var body = bodyCreateTrida{}
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("body musí mít jméno třídy"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	databaze.CreateTrida(body.Jmeno, id, utils.GenTridaKod())

	return c.SendStatus(fiber.StatusOK)
}

func tridy(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString(utils.GenTridaKod())
}
