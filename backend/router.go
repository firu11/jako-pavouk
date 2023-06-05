package main

import (
	"backend/databaze"
	"backend/utils"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setupRouter(app *fiber.App) {
	app.Get("/lekce", getVsechnyLekce)
	app.Get("/lekce/:pismena", getCviceniVLekci)
	app.Get("/cvic/:pismena/:cislo", getCviceni)
	app.Post("/dokonceno/:pismena/:cislo", dokoncitCvic)
	app.Post("/registrace", registrace)
	app.Post("/prihlaseni", prihlaseni)
	app.Get("/ja", prehled)
	app.Get("/test", test)
}

func test(c *fiber.Ctx) error {
	x, err := databaze.GetDokonceneCvicVLekci(2)
	if err != nil {
		return err
	}
	return c.JSON(x)
}

func getVsechnyLekce(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	lekce, err := databaze.GetLekce()
	if err != nil {
		return err
	}
	if id != 0 {
		dokoncene, err := databaze.GetDokonceneLekce(id)
		if err != nil {
			return err
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": dokoncene})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": []int{}})
}

func getCviceniVLekci(c *fiber.Ctx) error {
	_, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(c.Params("pismena"))
	if err != nil {
		return err
	}
	id, _ := databaze.GetLekceIDbyPismena(c.Params("pismena"))
	doko, err := databaze.GetDokonceneCvicVLekci(id)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"cviceni": cvic, "dokoncene": doko})
}

func getCviceni(c *fiber.Ctx) error {
	pismena := c.Params("pismena")
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(pismena)
	if err != nil {
		return err
	}
	cislo, err := strconv.Atoi(c.Params("cislo")) // str -> int
	if err != nil {
		return err
	}
	if cislo > len(vsechnyCviceni) {
		return c.Status(http.StatusBadRequest).JSON("Cviceni neexistuje")
	}
	if vsechnyCviceni[cislo-1].Typ == "nova" {
		var text []string = []string{}
		for i := 0; i < pocetSlov; i++ {
			var slovo string = ""
			for j := 0; j < pocetPismenVeSlovu; j++ {
				slovo += string(pismena[rand.Intn(len(pismena))])
			}
			slovo += " "
			text = append(text, slovo)
		}
		text[len(text)-1] = text[len(text)-1][:len(text[len(text)-1])-1] // smazat mezeru na konci
		return c.Status(http.StatusOK).JSON(fiber.Map{"text": text})
	} else {
		return c.Status(http.StatusOK).JSON(fiber.Map{"text": "Bruh GG"}) //TODO
	}
}

func dokoncitCvic(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	var body struct {
		CPM      float32 `json:"cpm" validate:"required"`
		Preklepy int     `json:"preklepy" validate:"required,number"` //sus nebere nulu
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	if err := utils.ValidateStruct(&body); err != nil {
		return err
	}
	cislo, err := strconv.ParseUint(c.Params("cislo"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Id cviceni je spatne")
	}

	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(c.Params("pismena"))
	if err != nil {
		return err
	}

	databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.CPM, body.Preklepy)

	return nil
}

func registrace(c *fiber.Ctx) error {
	// overeni spravnych dat co prijdou
	var body struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=25"`
		Heslo string `json:"heslo" validate:"required,min=8,max=25"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return err
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(http.StatusBadRequest).JSON("Invalidni email")
	}

	// overeni jestli nahodou uz neexistuje
	_, err = databaze.GetUzivByEmail(body.Email)
	if err != nil {
		hesloHASH, err := utils.HashPassword(body.Heslo)
		if err != nil {
			return err
		}
		id, err := databaze.CreateUziv(body.Email, hesloHASH, body.Jmeno)
		if err != nil {
			return err
		}
		token, err := utils.GenerovatToken(body.Email, id)
		if err != nil {
			return err
		} else {
			return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
		}
	} else {
		return c.Status(http.StatusBadRequest).JSON("Uzivatel jiz existuje")
	}
}

func prihlaseni(c *fiber.Ctx) error {
	// validace body dat
	var body struct {
		Email string `json:"email" validate:"required,email"`
		Heslo string `json:"heslo" validate:"required,min=8,max=25"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return err
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(http.StatusBadRequest).JSON("Invalidni email")
	}
	// pull z databaze
	uziv, err := databaze.GetUzivByEmail(body.Email)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Uzivatel neexistuje")
	}

	if err := utils.CheckPassword(body.Heslo, uziv.Heslo); err != nil {
		return c.Status(http.StatusUnauthorized).JSON("Spatne heslo")
	} else {
		token, err := utils.GenerovatToken(uziv.Email, uziv.ID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON("Token se pokazil")
		} else {
			return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
		}
	}
}

func prehled(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return err
	}
	/* prumerPreklepu := utils.Prumer(uziv.Preklepy)
	prumerRychlosti := utils.Prumer(uziv.Rychlosti) */
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"email":           uziv.Email,
		"jmeno":           uziv.Jmeno,
		"prumerPreklepu":  "prumerPreklepu",
		"prumerRychlosti": "prumerRychlosti",
	})
}
