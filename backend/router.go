package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
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
	/* databaze.PushSlovnik() */
	return c.JSON("sussy")
}

func getVsechnyLekce(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}

	lekce, err := databaze.GetLekce()
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	if id != 0 {
		dokoncene, err := databaze.GetDokonceneLekce(id)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": dokoncene})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": []int{}})
}

func getCviceniVLekci(c *fiber.Ctx) error {
	uzivID, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return HTTPerr
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(pismena)
	if err != nil {
		log.Print("Takova lekce neexistuje")
		return fiber.ErrBadRequest
	}
	id, _ := databaze.GetLekceIDbyPismena(c.Params("pismena"))
	doko, err := databaze.GetDokonceneCvicVLekci(uzivID, id)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cviceni": cvic, "dokoncene": doko})
}

func getCviceni(c *fiber.Ctx) error {
	_, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(pismena)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	cislo, err := strconv.Atoi(c.Params("cislo")) // str -> int
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	if cislo > len(vsechnyCviceni) {
		return c.Status(http.StatusBadRequest).JSON("Cviceni neexistuje")
	}

	var text []string

	switch vsechnyCviceni[cislo-1].Typ {
	case "nova":
		for i := 0; i < pocetSlov; i++ {
			var slovo string = ""
			for j := 0; j < pocetPismenVeSlovu; j++ {
				r := rand.Intn(utf8.RuneCountInString(pismena)) // utf-8 jsou sus
				slovo += string([]rune(pismena)[r])
			}
			slovo += " "
			text = append(text, slovo)
		}
	case "naucena":
		naucenaPismena, err := databaze.GetNaucenaPismena(pismena)
		if err != nil {
			return err
		}

		for i := 0; i < pocetSlov; i++ {
			var slovo string = ""
			for j := 0; j < pocetPismenVeSlovu; j++ {
				r := rand.Intn(utf8.RuneCountInString(naucenaPismena)) // utf-8 jsou sus
				slovo += string([]rune(naucenaPismena)[r])
			}
			slovo += " "
			text = append(text, slovo)
		}
	case "slova":
		slova, err := databaze.GetSlovaProLekci(pismena)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		for i := 0; i < pocetSlov; i++ {
			text = append(text, slova[rand.Intn(len(slova))]+" ")
		}
	default:
		log.Print("Cviceni ma divnej typ")
		return fiber.ErrInternalServerError
	}

	text[len(text)-1] = text[len(text)-1][:len(text[len(text)-1])-1] // smazat mezeru na konci
	return c.Status(http.StatusOK).JSON(fiber.Map{"text": text})
}

func dokoncitCvic(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	var body struct {
		CPM      float32 `json:"cpm" validate:"required"`
		Preklepy int     `json:"preklepy" validate:"min=0"` //sus reqired nebere nulu takze min=0 asi ok
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	cislo, err := strconv.ParseUint(c.Params("cislo"), 10, 32)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}

	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(c.Params("pismena"))
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	if int(cislo-1) >= len(vsechnyCviceni) { // error index out of range nebude
		log.Print("Takovy cviceni neni")
		return fiber.ErrBadRequest
	}
	if err := databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.CPM, body.Preklepy); err != nil {
		err = databaze.OdebratDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		err = databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.CPM, body.Preklepy)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
	}

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
		log.Print(err)
		return fiber.ErrInternalServerError
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
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		id, err := databaze.CreateUziv(body.Email, hesloHASH, body.Jmeno)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		token, err := utils.GenerovatToken(body.Email, id)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		} else {
			return c.Status(http.StatusOK).JSON(fiber.Map{"token": token, "expiryTime": tokenTimeDuration / 1_000_000})
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
		log.Print(err)
		return fiber.ErrInternalServerError
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
			return c.Status(http.StatusOK).JSON(fiber.Map{"token": token, "expiryTime": tokenTimeDuration / 1_000_000})
		}
	}
}

func prehled(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	preklepy, cpm, err := databaze.GetPreklepyACPM(id)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"email":           uziv.Email,
		"jmeno":           uziv.Jmeno,
		"uspesnost":       (float32(delkaTextu) - utils.Prumer(preklepy)) / float32(delkaTextu) * 100,
		"prumerRychlosti": utils.Prumer(cpm),
		"dokonceno":       dokonceno,
	})
}
