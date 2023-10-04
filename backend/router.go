package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"math/rand"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type (
	bodyDokoncitCvic struct {
		CPM      float32 `json:"cpm" validate:"required"`
		Preklepy int     `json:"preklepy" validate:"min=0"` //sus reqired nebere nulu takze min=0 asi ok
		Cas      float32 `json:"cas" validate:"required"`
	}

	bodyPoslatEmail struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
	}

	bodyRegistrace struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=25"`
		Heslo string `json:"heslo" validate:"required,min=5,max=25"`
	}

	bodyPrihlaseni struct {
		EmailNeboJmeno string `json:"email" validate:"required"`
		Heslo          string `json:"heslo" validate:"required,min=5,max=25"`
	}

	bodyUprava struct {
		Zmena   string `json:"zmena"`
		Hodnota string `json:"hodnota"`
	}
)

func SetupRouter(app *fiber.App) {
	app.Get("/api/lekce", getVsechnyLekce)
	app.Get("/api/lekce/:pismena", getCviceniVLekci)
	app.Get("/api/cvic/:pismena/:cislo", getCviceni)
	app.Post("/api/dokonceno/:pismena/:cislo", dokoncitCvic)
	app.Post("/api/overit-email", overitEmail)
	app.Post("/api/registrace", registrace)
	app.Post("/api/prihlaseni", prihlaseni)
	app.Get("/api/ja", prehled)
	app.Post("/api/ucet-zmena", upravaUctu)
	app.Get("/api/token-expirace", testVyprseniTokenu)
	app.Get("/api/test", test)
}

func chyba(msg string) fiber.Map {
	if msg == "" {
		msg = "Neco se pokazilo"
	}
	return fiber.Map{"error": msg}
}

func test(c *fiber.Ctx) error {
	/* log.Println(utils.UzivCekajiciNaOvereni)
	log.Println(utils.ValidFormat("firu"))
	utils.MobilNotifikace("Jmeno - email@ema.il") */
	/* databaze.PushCviceni() */
	return c.JSON("Vypadni")
}

func getVsechnyLekce(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}

	lekce, err := databaze.GetLekce(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	if id != 0 {
		dokoncene, err := databaze.GetDokonceneLekce(id)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": dokoncene})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": []int{}})
}

func getCviceniVLekci(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(HTTPerr.Error()))
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Takova lekce neexistuje"))
	}
	doko, err := databaze.GetDokonceneCvicVLekci(id, 0, pismena)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cviceni": cvic, "dokoncene": doko})
}

func getCviceni(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
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
		return c.Status(fiber.StatusBadRequest).JSON("Cviceni neexistuje")
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
		naucenaPismena, err := databaze.GetNaucenaPismena(id, pismena)
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
		slova, err := databaze.GetSlovaProLekci(id, pismena)
		if err != nil {
			log.Print(err)
			return fiber.ErrInternalServerError
		}
		if len(slova) < pocetSlov {
			log.Println("neni dost slov")
		}
		for i := 0; i < pocetSlov; i++ {
			text = append(text, slova[rand.Intn(len(slova))]+" ")
		}
	default:
		log.Print("Cviceni ma divnej typ")
		return fiber.ErrInternalServerError
	}

	var posledni bool = false
	if int(cislo-1) == len(vsechnyCviceni)-1 {
		posledni = true
	}

	text[len(text)-1] = text[len(text)-1][:len(text[len(text)-1])-1] // smazat mezeru na konci

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "klavesnice": u.Klavesnice, "posledni": posledni})
}

func dokoncitCvic(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	var body = bodyDokoncitCvic{}

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
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		log.Print(err)
		return fiber.ErrInternalServerError
	}
	if int(cislo-1) >= len(vsechnyCviceni) { // error index out of range nebude
		log.Print("Takovy cviceni neni")
		return fiber.ErrBadRequest
	}

	err = databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.CPM, body.Preklepy, body.Cas)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	return c.SendStatus(fiber.StatusOK)
}

func overitEmail(c *fiber.Ctx) error {
	var body bodyPoslatEmail = bodyPoslatEmail{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Invalidni email"))
	}

	cekajiciUziv, err := databaze.GetNeoverenyUziv(body.Email)
	if err != nil {
		go databaze.SmazatNeoverenyPoLimitu()
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkuste to prosim znovu"))
	}

	if time.Now().Unix() <= cekajiciUziv.Cas && cekajiciUziv.Kod != body.Kod { //vsechno dobry ale spatnej kod
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny kod"))
	} else if time.Now().Unix() > cekajiciUziv.Cas { //vyprselo
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkuste to prosim znovu"))
	}

	uzivID, err := databaze.CreateUziv(cekajiciUziv.Email, cekajiciUziv.Heslo, cekajiciUziv.Jmeno)
	if err != nil {
		log.Println(err, uzivID)
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	token, err := utils.GenerovatToken(body.Email, uzivID)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	go utils.MobilNotifikace(cekajiciUziv.Jmeno + " - " + body.Email)
	go databaze.OdebratOvereni(cekajiciUziv.Email)
	go databaze.SmazatNeoverenyPoLimitu()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func registrace(c *fiber.Ctx) error {
	// overeni spravnych dat co prijdou
	var body bodyRegistrace = bodyRegistrace{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Invalidni email"))
	}
	hesloHASH, err := utils.HashPassword(body.Heslo)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	if _, err = databaze.GetUzivByEmail(body.Email); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s tímto emailem jiz existuje"))
	}

	if _, err = databaze.GetUzivByJmeno(body.Jmeno); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s tímto jmenem jiz existuje"))
	}

	var randomKod string = utils.GenKod()

	if err := databaze.CreateNeoverenyUziv(body.Email, hesloHASH, body.Jmeno, randomKod, time.Now().Add(10*time.Minute).Unix()); err != nil {
		go databaze.SmazatNeoverenyPoLimitu()
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s tímto jmenem docasne existuje"))
	}
	if err := utils.PoslatOverovaciEmail(body.Email, randomKod); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	go databaze.SmazatNeoverenyPoLimitu()
	return c.SendStatus(fiber.StatusOK)
}

func prihlaseni(c *fiber.Ctx) error {
	// validace body dat
	var body bodyPrihlaseni = bodyPrihlaseni{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	var uziv databaze.Uzivatel

	// validace emailu
	if !utils.ValidFormat(body.EmailNeboJmeno) { //predpokladam ze je to jmeno kdyz se to nepodobá emailu
		uziv, err = databaze.GetUzivByJmeno(body.EmailNeboJmeno)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Email/jmeno je spatne (Jmeno)"))
		}
	} else {
		uziv, err = databaze.GetUzivByEmail(body.EmailNeboJmeno)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Email/jmeno je spatne (Email)"))
		}
	}

	if err := utils.CheckPassword(body.Heslo, uziv.Heslo); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba("Heslo je spatne"))
	} else {
		token, err := utils.GenerovatToken(uziv.Email, uziv.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba("Token se pokazil"))
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
		}
	}
}

func prehled(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	preklepy, cpm, daystreak, cas, err := databaze.GetUdaje(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	uspesnost := float32(delkaTextu) - utils.Prumer(preklepy)
	if uspesnost < 0 {
		uspesnost = 0 // kvuli adamovi kterej big troulin a měl -10%
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"email":           uziv.Email,
		"jmeno":           uziv.Jmeno,
		"daystreak":       daystreak,
		"uspesnost":       uspesnost / float32(delkaTextu) * 100,
		"prumerRychlosti": utils.Prumer(cpm),
		"celkovyCas":      cas,
		"dokonceno":       dokonceno,
		"klavesnice":      uziv.Klavesnice,
	})
}

func testVyprseniTokenu(c *fiber.Ctx) error {
	if len(c.Get("Authorization")) >= 10 { // treba deset proste at tam neco je
		jePotrebaVymenit, err := utils.ValidovatExpTokenu(c.Get("Authorization")[7:])
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}
		id, err := utils.Autentizace(c, true)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}
		_, err = databaze.GetUzivByID(id)
		if err != nil && !jePotrebaVymenit {
			jePotrebaVymenit = true
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"jePotrebaVymenit": jePotrebaVymenit})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(""))
	}
}

func upravaUctu(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	var body = bodyUprava{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	if body.Zmena == "smazat" {
		if err := databaze.SmazatUzivatele(id); err != nil {
			return err
		}
	} else if body.Zmena == "klavesnice" {
		databaze.ZmenitKlavesnici(id, body.Hodnota)
	} else if body.Zmena == "jmeno" {
		err := databaze.PrejmenovatUziv(id, body.Hodnota)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
	}
	return nil
}
