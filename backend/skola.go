package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type (
	bodyCreateTrida struct {
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
	bodyZmenaTridy struct {
		TridaID uint   `json:"trida_id" validate:"required"`
		Zmena   string `json:"zmena"`
		Hodnota string `json:"hodnota"`
	}
	bodyZmenaJmenaStudenta struct {
		ID    uint   `json:"id" validate:"required"`
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
	bodyZapis struct {
		Kod   string `json:"kod" validate:"required"`
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
)

// typy uživatelů
// 1 - basic
// 2 - učitel
func setupSkolniRouter(api *fiber.Router) {
	skolaApi := (*api).Group("/skola")

	skolaApi.Post("/create-trida", createTrida)
	skolaApi.Get("/tridy", tridy)
	skolaApi.Get("/trida/:id", trida)
	skolaApi.Get("/test-tridy/:kod", testTridy)
	skolaApi.Post("/zmena-tridy", zmenaTridy)

	skolaApi.Get("/student/:id", student)
	skolaApi.Post("/student", studentPrejmenovat)
	skolaApi.Post("/zapis", zapis)
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

	var body bodyCreateTrida
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("body musí mít jméno třídy"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	err = databaze.CreateTrida(body.Jmeno, id, utils.GenTridaKod())
	if err != nil { // kod asi neni unikatni
		err = databaze.CreateTrida(body.Jmeno, id, utils.GenTridaKod())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func tridy(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tridy muze videt pouze ucitel"))
	}

	tridy, err := databaze.GetTridy(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	output := make(map[string][]databaze.TridaInfo)
	for _, trida := range tridy {
		cislo := strings.Split(strings.Split(trida.Jmeno, " ")[0], ".")[0]
		if output[cislo] == nil {
			output[cislo] = []databaze.TridaInfo{trida}
		} else {
			output[cislo] = append(output[cislo], trida)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"tridy": output})
}

func trida(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tridy muze videt pouze ucitel"))
	}

	tridaID, err := strconv.ParseInt(c.Params("id"), 10, 8)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	trida, err := databaze.GetTrida(uint(tridaID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	if trida.UcitelID != uziv.ID {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("To neni tvoje trida"))
	}
	studenti, err := databaze.GetStudentyZeTridy(trida.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"trida": trida, "studenti": studenti})
}

func testTridy(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	_, err = databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	trida, err := databaze.GetTridaByKod(strings.ToUpper(c.Params("kod")))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Takova trida neexistuje"))
	}
	if trida.Zamknuta {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Trida je zamcena"))
	}

	return c.SendStatus(fiber.StatusOK)
}

func zmenaTridy(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tridy muze upravovat pouze ucitel"))
	}

	var body bodyZmenaTridy
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}

	switch body.Zmena {
	case "zamek":
		err := databaze.ZamknoutTridu(body.TridaID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}
	case "jmeno":
		if body.Hodnota != "" {
			err := databaze.PrejmenovatTridu(body.TridaID, body.Hodnota)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Potrebuju hodnotu"))
		}
	case "smazat":
		err := databaze.SmazatTridu(body.TridaID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func student(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tohle muze pouze ucitel"))
	}

	studentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	student, err := databaze.GetUzivByID(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	presnost, cpm, daystreak, _, chybyPismenka, err := databaze.GetUdaje(uint(studentID))
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"email":            student.Email,
		"jmeno":            student.SkolniJmeno,
		"daystreak":        daystreak,
		"uspesnost":        presnost,
		"medianRychlosti":  utils.Median(cpm),
		"dokonceno":        dokonceno,
		"nejcastejsiChyby": chybyPismenka,
		"klavesnice":       student.Klavesnice,
	})
}

func studentPrejmenovat(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if uziv.Role != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Tohle muze pouze ucitel"))
	}

	var body bodyZmenaJmenaStudenta
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}

	err = databaze.PrejmenovatStudenta(body.ID, body.Jmeno)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.SendStatus(fiber.StatusOK)
}

func zapis(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	var body bodyZapis
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}

	err = databaze.ZapsatStudenta(body.Kod, id, body.Jmeno)
	if err != nil {
		log.Println(err)
		if err.Error() == "uz je ve tride" {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Uz jsi ve tride"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.SendStatus(fiber.StatusOK)
}
