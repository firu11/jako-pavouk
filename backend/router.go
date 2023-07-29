package main

import (
	"backend/databaze"
	"backend/utils"
	"fmt"
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
	}

	bodyPoslatEmail struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
	}

	bodyRegistrace struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=25"`
		Heslo string `json:"heslo" validate:"required,min=8,max=25"`
	}

	bodyPrihlaseni struct {
		EmailNeboJmeno string `json:"email" validate:"required"`
		Heslo          string `json:"heslo" validate:"required,min=8,max=25"`
	}

	bodyUprava struct {
		Jmeno  string `json:"jmeno" validate:"min=3,max=12"`
		Smazat bool   `json:"smazat"`
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
	fmt.Println(utils.UzivCekajiciNaOvereni)
	fmt.Println(utils.ValidFormat("firu"))
	/* databaze.PushSlovnik() */
	return c.JSON("Vypadni")
}

func getVsechnyLekce(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}

	lekce, err := databaze.GetLekce()
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
	uzivID, err := utils.Autentizace(c, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(err.Error()))
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(HTTPerr.Error()))
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(pismena)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Takova lekce neexistuje"))
	}
	doko, err := databaze.GetDokonceneCvicVLekci(uzivID, 0, pismena)
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
		for i := 0; i < 14; i++ {
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "posledni": posledni})
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
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(pismena)
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

	// najdeme cekajici mail
	var deleteIndexy []int
	var vPoradku bool = false
	var cekajiciUzivatel utils.UzivCekajici

	for i, uziv := range utils.UzivCekajiciNaOvereni {
		if time.Now().Unix() > uziv.DobaTrvani { // uz je po dobe trvani
			deleteIndexy = append(deleteIndexy, i)
		}
		if uziv.Email == body.Email {
			if time.Now().Unix() <= uziv.DobaTrvani && uziv.Kod != body.Kod { //vsechno dobry ale spatnej kod
				utils.VycistitSeznam(deleteIndexy)
				return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny kod"))
			} else if time.Now().Unix() > uziv.DobaTrvani { //vyprselo
				utils.VycistitSeznam(deleteIndexy)
				return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkuste to prosim znovu"))
			} else { // dobry
				vPoradku = true
				cekajiciUzivatel = uziv
				utils.VycistitSeznam(deleteIndexy)
			}
		}
	}
	if vPoradku {
		id, err := databaze.CreateUziv(body.Email, cekajiciUzivatel.HesloHash, cekajiciUzivatel.Jmeno)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
		token, err := utils.GenerovatToken(body.Email, id)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Prvne zkuste endpoint pro registraci aby vám byl poslán email"))
	}
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
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Uzivatel s tímto emailem jiz existuje"))
	}

	if _, err = databaze.GetUzivByJmeno(body.Jmeno); err == nil { // uz existuje
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Uzivatel s tímto jmenem jiz existuje"))
	}

	var randomKod string = utils.GenKod()

	for i, v := range utils.UzivCekajiciNaOvereni {
		if v.Email == body.Email && v.Jmeno == body.Jmeno {
			utils.UzivCekajiciNaOvereni[i].Kod = randomKod
			utils.PoslatOverovaciEmail(body.Email, randomKod)
			return c.SendStatus(fiber.StatusOK)
		} else if body.Email == v.Email {
			utils.UzivCekajiciNaOvereni[i].Kod = randomKod
			utils.UzivCekajiciNaOvereni[i].Jmeno = body.Jmeno
			utils.PoslatOverovaciEmail(body.Email, randomKod)
			return c.SendStatus(fiber.StatusOK)
		} else if v.Jmeno == body.Jmeno {
			utils.UzivCekajiciNaOvereni[i].Kod = randomKod
			utils.UzivCekajiciNaOvereni[i].Email = body.Email
			utils.PoslatOverovaciEmail(body.Email, randomKod)
			return c.SendStatus(fiber.StatusOK)
		}
	} // pokud jsme nenašli shodu tak pridame novy
	utils.UzivCekajiciNaOvereni = append(utils.UzivCekajiciNaOvereni, utils.UzivCekajici{Email: body.Email, Jmeno: body.Jmeno, HesloHash: hesloHASH, Kod: randomKod, DobaTrvani: time.Now().Add(15 * time.Minute).Unix()})
	utils.PoslatOverovaciEmail(body.Email, randomKod)
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
	preklepy, cpm, err := databaze.GetPreklepyACPM(id)
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
		"daystreak":       uziv.DayStreak,
		"uspesnost":       uspesnost / float32(delkaTextu) * 100,
		"prumerRychlosti": utils.Prumer(cpm),
		"dokonceno":       dokonceno,
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
	if err := databaze.ZmenitUzivatele(body.Jmeno, "", body.Smazat, id); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(body)
}
