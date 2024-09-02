package main

import (
	"backend/databaze"
	"backend/utils"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

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
	bodyPridatPraci struct {
		Text    string `json:"text" validate:"required,min=1"`
		Cas     int    `json:"cas" validate:"required"`
		TridaID uint   `json:"trida_id" validate:"required"`
	}
	bodyGetText struct {
		Typ string `json:"typ" validate:"required"`
	}

	praceProStudenta struct {
		ID       uint      `json:"id"`
		TridaID  uint      `json:"-"`
		Text     string    `json:"-"`
		Cas      int       `json:"cas"`
		Datum    time.Time `json:"datum"`
		Cpm      float64   `json:"cpm"`
		Presnost float64   `json:"presnost"`
	}
)

// typy uživatelů
// 1 - basic
// 2 - učitel
func setupSkolniRouter(api *fiber.Router) {
	skolaApi := (*api).Group("/skola")

	skolaApi.Post("/create-trida", createTrida)
	skolaApi.Get("/tridy", tridy)
	skolaApi.Get("/trida", tridaStudent)
	skolaApi.Get("/trida/:id", trida)
	skolaApi.Get("/test-tridy/:kod", testTridy)
	skolaApi.Post("/zmena-tridy", zmenaTridy)

	skolaApi.Post("/pridat-praci", pridatPraci)
	skolaApi.Get("/get-praci/:id", getPraci)
	skolaApi.Post("/dokoncit-praci/:id", dokoncitPraci)

	skolaApi.Post("/text", getText)

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

	for i := 0; i < 5; i++ {
		err = databaze.CreateTrida(body.Jmeno, id, utils.GenTridaKod())
		if err == nil {
			return c.SendStatus(fiber.StatusOK)
		}
		// kod asi neni unikatni tak zkusim znovu
	}

	// několikrat se to pokazilo tak hodim error
	return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
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
	prace, err := databaze.GetVsechnyPrace(trida.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"trida": trida, "studenti": studenti, "prace": prace})
}

func tridaStudent(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	trida, err := databaze.GetTridaByStudentID(uziv.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	prace, err := databaze.GetVsechnyPrace(trida.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	cpmka, presnost, err := databaze.GetDokoncenePrace(uziv.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	var vysledek []praceProStudenta
	for _, p := range prace {
		cpm, ok := cpmka[p.ID]
		if !ok || cpm < 0 {
			cpm = -1
		}
		pres, ok := presnost[p.ID]
		if !ok || pres < 0 {
			pres = -1
		}
		vysledek = append(vysledek, praceProStudenta{ID: p.ID, TridaID: p.TridaID, Text: p.Text, Cas: p.Cas, Datum: p.Datum, Cpm: cpm, Presnost: pres})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"trida": trida, "prace": vysledek})
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

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	student, err := databaze.GetUzivByID(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	presnost, cpm, daystreak, chybyPismenka, err := databaze.GetUdaje(student.ID)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(student.ID)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	cpmkaVPracich, presnostVPracich, err := databaze.GetDokoncenePrace(student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"student": fiber.Map{
			"email":            student.Email,
			"jmeno":            student.SkolniJmeno,
			"daystreak":        daystreak,
			"uspesnost":        presnost,
			"rychlost":         utils.Prumer(cpm),
			"dokonceno":        dokonceno,
			"nejcastejsiChyby": chybyPismenka,
			"klavesnice":       student.Klavesnice,
		},
		"cpmVPracich":      cpmkaVPracich,
		"presnostVPracich": presnostVPracich,
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

func pridatPraci(c *fiber.Ctx) error {
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

	var body bodyPridatPraci
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}

	if err = databaze.PridatPraci(body.Text, body.Cas, body.TridaID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	return c.SendStatus(fiber.StatusOK)
}

func getText(c *fiber.Ctx) error {
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

	var body bodyGetText
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Spatny body"))
	}

	var text strings.Builder

	if body.Typ == "Věty z pohádek" {
		vety, err := databaze.GetVsechnyVety(int(pocetZnaku / 85)) // cca 85 znaku na vetu
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}

		for _, v := range vety {
			text.WriteString(v)
			text.WriteString(" ")
		}
	} else {
		txt, err := databaze.GetRandomProcvic(body.Typ)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Takový typ nemáme"))
		}
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}

		text.WriteString(txt)
	}

	var vyslednyText string = text.String()
	if string(vyslednyText[len(vyslednyText)-1]) == " " {
		vyslednyText = vyslednyText[:len(vyslednyText)-1]
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": vyslednyText})
}

func getPraci(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	praceID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	text, cas, err := databaze.GetPrace(uint(praceID), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	var vyslednyText []string

	slova := strings.Split(text, " ")
	if slova[len(slova)-1] == "" {
		slova = slova[:len(slova)-1]
	}
	for _, v := range slova {
		vyslednyText = append(vyslednyText, v+" ")
	}
	utils.SmazatMezeruNaKonci(vyslednyText)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": vyslednyText, "cas": cas, "klavesnice": uziv.Klavesnice})
}

func dokoncitPraci(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	var body = bodyDokoncit{}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	praceID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	log.Println(praceID)

	err = databaze.DokoncitPraci(uint(praceID), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	return c.SendStatus(fiber.StatusOK)
}
