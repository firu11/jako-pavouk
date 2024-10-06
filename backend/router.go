package main

import (
	"backend/databaze"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"zgo.at/isbot"
)

// struct body requestu
type (
	bodyDokoncit struct {
		Preklepy         int            `json:"neopravenePreklepy" validate:"min=0"` //sus reqired nebere nulu takze min=0 asi ok
		Cas              float32        `json:"cas" validate:"required"`
		DelkaTextu       int            `json:"delkaTextu" validate:"required"`
		NejcastejsiChyby map[string]int `json:"nejcastejsiChyby" validate:"required"`
	}

	bodyPoslatEmail struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
	}

	bodyRegistrace struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=16"`
		Heslo string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyPrihlaseni struct {
		EmailNeboJmeno string `json:"email" validate:"required"`
		Heslo          string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyUprava struct {
		Zmena   string `json:"zmena"`
		Hodnota string `json:"hodnota"`
	}

	bodyZmenaHesla struct {
		Email string `json:"email" validate:"required,email"`
	}

	bodyOvereniZmenaHesla struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
		Heslo string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyGoogle struct {
		AccessToken string `json:"access_token"`
	}

	bodyTestPsani struct {
		Typ string `json:"typ" validate:"required"`
	}
)

// vytvoří skupinu /api a v ní všechny endpointy
func SetupRouter(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/lekce", getVsechnyLekce)
	api.Get("/lekce/:pismena", getCviceniVLekci)
	api.Get("/cvic/:pismena/:cislo", getCviceni)
	api.Post("/dokonceno/:pismena/:cislo", dokoncitCvic)
	api.Post("/dokonceno-procvic/:cislo", dokoncitProcvic)
	api.Get("/procvic", getVsechnyProcvic)
	api.Get("/procvic/:cisloProcvic/:neCislo", getProcvic)
	api.Post("/test-psani", testPsani)

	api.Post("/overit-email", overitEmail)
	api.Post("/registrace", registrace)
	api.Post("/prihlaseni", prihlaseni)
	api.Post("/zmena-hesla", zmenaHesla)
	api.Post("/overeni-zmeny-hesla", overitZmenuHesla)
	api.Post("/google", google)

	api.Get("/ja", prehled)
	api.Post("/ucet-zmena", upravaUctu)

	api.Get("/token-expirace", testVyprseniTokenu)
	api.Post("/navsteva", navsteva)

	setupSkolniRouter(&api)
}

// standardní chybový výstup
func chyba(msg string) fiber.Map {
	if msg == "" {
		msg = "Neco se pokazilo"
	}
	return fiber.Map{"error": msg}
}

// vygeneruje text pro test psaní
//
// potřebuje délku textu a jeho typ: slova / věty
func testPsani(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)

	var body = bodyTestPsani{}
	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	var text []string
	switch body.Typ {
	case "slova":
		var err error
		text, err = databaze.GetVsechnySlova(int(pocetZnaku / 7.5)) // cca 8.5 znaku na slovo
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}

		for i := 0; i < len(text); i++ {
			text[i] += " "

			r := []rune(text[i])
			if i%5 == 0 { //kazdy paty velkym
				text[i] = fmt.Sprintf("%c%s", unicode.ToUpper(r[0]), string(r[1:]))
			}
		}

	case "vety":
		vety, err := databaze.GetVsechnyVety(int(pocetZnaku / 85)) // cca 85 znaku na vetu
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}
		for i := 0; i < len(vety); i++ {
			slova := strings.Split(vety[i], " ")
			for _, v := range slova {
				text = append(text, v+" ")
			}
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny typ testu psani"))
	}

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "klavesnice": u.Klavesnice})
}

// vrací seznam všech lekcí v závislosti na klávesnici
func getVsechnyLekce(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)

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

		lekcePismena, cisloCvic, err := databaze.GetDalsiCviceni(id)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": dokoncene, "dalsi_cviceni": fmt.Sprintf("/%s/%d", lekcePismena, cisloCvic)})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"lekce": lekce, "dokoncene": []int{}, "dalsi_cviceni": ""})
}

// vrací všechny cvičení v lekci podle písmen lekce z parametru url
func getCviceniVLekci(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if HTTPerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(HTTPerr.Error()))
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Takova lekce neexistuje"))
	}
	doko, err := databaze.GetDokonceneCvicVLekci(id, 0, pismena)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cviceni": cvic, "dokoncene": doko})
}

// generuje texty pro cvičení
//
// text vrací v závislosti na jeho typu: nové písmena, naučená písmena, slova, + nějaké speciální (programator...)
func getCviceni(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)
	if id == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if HTTPerr != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	cislo, err := strconv.Atoi(c.Params("cislo")) // str -> int
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if cislo > len(vsechnyCviceni) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cviceni neexistuje"))
	}

	text, err := generovatTextCviceni(pismena, vsechnyCviceni[cislo-1].Typ, id, 0, pocetZnaku)

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "klavesnice": u.Klavesnice, "typ": vsechnyCviceni[cislo-1].Typ, "posledni": int(cislo-1) == len(vsechnyCviceni)-1})
}

// přidá do databáze záznam o tom jak uživatel cvičení napsal
//
// potřebuje token uživatele, rychlost, preklepy, cas, delku textu
func dokoncitCvic(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)
	if id == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
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
	cislo, err := strconv.ParseUint(c.Params("cislo"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	pismena, HTTPerr := utils.DecodeURL(c.Params("pismena"))
	if HTTPerr != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if int(cislo-1) >= len(vsechnyCviceni) { // error index out of range nebude
		log.Print("Takovy cviceni neni")
		return fiber.ErrBadRequest
	}

	err = databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	return c.SendStatus(fiber.StatusOK)
}

// přidá do databáze záznam o tom jak uživatel cvičení napsal
//
// potřebuje token uživatele, rychlost, preklepy, cas, delku textu
func dokoncitProcvic(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)

	var body = bodyDokoncit{}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	cislo, err := strconv.ParseUint(c.Params("cislo"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	vsechnyProcvic, err := databaze.GetTexty()
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	if cislo > uint64(len(vsechnyProcvic)) { // error index out of range nebude
		log.Print("Takovy procvicovani neni")
		return fiber.ErrBadRequest
	}
	err = databaze.PridatDokonceneProcvic(uint(cislo), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	return c.SendStatus(fiber.StatusOK)
}

// vrátí seznam textů k procvičování
func getVsechnyProcvic(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)

	texty, err := databaze.GetTexty()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	if id == 0 {
		for i := 0; i < len(texty); i++ {
			texty[i].CPM = -1
		}
		return c.JSON(fiber.Map{"texty": texty, "testPsaniCPM": -1})
	}

	rychlosti, err := databaze.GetRychlostiProcvic(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	for i := 0; i < len(texty); i++ {
		r := rychlosti[int(texty[i].ID)-1]
		if r == 0 {
			texty[i].CPM = -1
		} else {
			texty[i].CPM = r
		}
	}

	var testPsaniCPM float32
	if rychlosti[-1] == 0 {
		testPsaniCPM = -1
	} else {
		testPsaniCPM = rychlosti[-1]
	}

	return c.JSON(fiber.Map{"texty": texty, "testPsaniCPM": testPsaniCPM})
}

// vrací text k odpovídajícímu procvičování
func getProcvic(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)

	typ, err := strconv.Atoi(c.Params("cisloProcvic")) // str -> int
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	cislo, err := strconv.Atoi(c.Params("neCislo")) // str -> int
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	nazev, podnazev, text, cislo, err := databaze.GetProcvicovani(typ, cislo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "typ": nazev, "jmeno": podnazev, "cislo": cislo})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "typ": nazev, "jmeno": podnazev, "cislo": cislo, "klavesnice": u.Klavesnice})
}

// porovná kód který byl uživateli zaslán na email s tím který mu přišel
//
// kontroluje také zda nevypršel čas ba ověření a maže asynchroně ty, kterým čas vypršel
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
		go databaze.SmazatPoLimitu()
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkus to prosim znovu 1"))
	}

	if time.Now().Unix() <= cekajiciUziv.Cas && !utils.CheckKod(cekajiciUziv.Kod, body.Kod) { // vsechno dobry ale spatnej kod
		databaze.DalSpatnyKod(body.Email)
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny kod"))
	} else if time.Now().Unix() > cekajiciUziv.Cas { // vyprselo
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkus to prosim znovu 2"))
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
	go databaze.OdebratOvereni(cekajiciUziv.Email)
	go databaze.SmazatPoLimitu()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

// # Registrace
//  1. přidá neověřeného uživatele do db
//  2. vygeneruje a odešle kód na email
//  3. smaže neověřené uživ. po limitu
func registrace(c *fiber.Ctx) error {
	// overeni spravnych dat co prijdou
	var body bodyRegistrace

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	body.Email = strings.ToLower(body.Email)
	if err := utils.ValidaceEmailu(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(fmt.Sprintf("Invalidni email: %s", err)))
	}
	if _, err := databaze.GetUzivByEmail(body.Email); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s timto emailem jiz existuje"))
	}

	if !regexJmeno.MatchString(body.Jmeno) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Invalidni jmeno"))
	}
	if _, err := databaze.GetUzivByJmeno(body.Jmeno); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s timto jmenem jiz existuje"))
	}

	hesloHASH, err := utils.HashPassword(body.Heslo)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	var randomKod string = utils.GenKod()

	if err := databaze.CreateNeoverenyUziv(body.Email, hesloHASH, body.Jmeno, randomKod, time.Now().Add(10*time.Minute).Unix()); err != nil {
		go databaze.SmazatPoLimitu()
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s timto jmenem docasne existuje"))
	}
	if err := utils.PoslatOverovaciEmail(body.Email, randomKod); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	go databaze.SmazatPoLimitu()
	return c.SendStatus(fiber.StatusOK)
}

// kontroluje hashe hesel z db a z frontendu. v případě že se shodují, vygeneruje a vrátí token
//
// také ošetřuje účty které jsou registrované přes google -> nemám jejich heslo
func prihlaseni(c *fiber.Ctx) error {
	var body bodyPrihlaseni

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
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
		if err.Error() == "ucet je pres google" {
			return c.Status(fiber.StatusUnauthorized).JSON(chyba("Ucet je registrovan pres google"))
		}
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

// # Google přihlášení / registrace
func google(c *fiber.Ctx) error {
	var body bodyGoogle
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	email, jmeno, err := databaze.GoogleTokenNaData(body.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	var token string
	uziv, err := databaze.GetUzivByEmail(email)
	if err != nil { // neexistuje
		id, err := databaze.CreateUziv(email, "google", jmeno)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
		token, err = utils.GenerovatToken(email, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba("Token se pokazil"))
		}
	} else {
		token, err = utils.GenerovatToken(email, uziv.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba("Token se pokazil"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

// posílá ověřovací email s kódem pro ověření hesla
func zmenaHesla(c *fiber.Ctx) error {
	var body bodyZmenaHesla
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	_, err = databaze.GetUzivByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Email je spatne"))
	}

	kod := utils.GenKod()
	if err := utils.PoslatOverovaciEmail(body.Email, kod); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	if err = databaze.CreateZapomenuteHeslo(body.Email, kod, time.Now().Add(10*time.Minute).Unix()); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	go databaze.SmazatPoLimitu()

	return c.SendStatus(fiber.StatusOK)
}

// ověřuje kód který byl zaslán na email s kódem z frontendu
func overitZmenuHesla(c *fiber.Ctx) error {
	var body bodyOvereniZmenaHesla
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	uziv, err := databaze.GetZmenuHesla(body.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	hesloHASH, err := utils.HashPassword(body.Heslo)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	if time.Now().Unix() <= uziv.Cas && !utils.CheckKod(uziv.Kod, body.Kod) { // vsechno dobry ale spatnej kod
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny kod"))
	} else if time.Now().Unix() > uziv.Cas { // vyprselo
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkuste to prosim znovu"))
	}

	err = databaze.ZmenitHeslo(uziv.Email, hesloHASH)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	go databaze.OdebratZmenuHesla(uziv.Email)
	go databaze.SmazatPoLimitu()

	return c.SendStatus(fiber.StatusOK)
}

// vrací statistiky o uživateli podle id z tokenu
func prehled(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)
	if id == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	presnost, cpm, daystreak, chybyPismenka, err := databaze.GetUdaje(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	trida, _ := databaze.GetTridaByUziv(uziv.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"email":            uziv.Email,
		"jmeno":            uziv.Jmeno,
		"daystreak":        daystreak,
		"uspesnost":        presnost,
		"rychlost":         utils.Prumer(cpm),
		"dokonceno":        dokonceno,
		"nejcastejsiChyby": chybyPismenka,
		"klavesnice":       uziv.Klavesnice,
		"role":             utils.GetRole(uziv.Role, trida.ID),
	})
}

// endpoint který vrací zda je potřeba token co nejdříve vyměnit
func testVyprseniTokenu(c *fiber.Ctx) error {
	if len(c.Get("Authorization")) < 10 { // treba deset proste at tam neco je
		return c.Status(fiber.StatusUnauthorized).JSON(chyba(""))
	}
	jePotrebaVymenit, err := utils.ValidovatExpTokenu(c.Get("Authorization")[7:])
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"jePotrebaVymenit": true})
	}
	id := c.Locals("uzivID").(uint)
	if id == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil && !jePotrebaVymenit {
		jePotrebaVymenit = true
	}
	trida, _ := databaze.GetTridaByUziv(uziv.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"jePotrebaVymenit": jePotrebaVymenit, "role": utils.GetRole(uziv.Role, trida.ID)})
}

// monitoruje přibližně návštěvníky které na stránku chodí
//
// snaží se filtrovat requesty od botů a crawlerů
func navsteva(c *fiber.Ctx) error {
	var httpRequest http.Request
	err := fasthttpadaptor.ConvertRequest(c.Context(), &httpRequest, false)
	if !isbot.Is(isbot.Bot(&httpRequest)) && err == nil {
		databaze.NovaNavsteva()
	}
	return c.SendStatus(fiber.StatusOK)
}

// mění buď jméno uživatele nebo klávesnici
func upravaUctu(c *fiber.Ctx) error {
	id := c.Locals("uzivID").(uint)
	if id == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
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
		databaze.ZmenitKlavesnici(id, strings.ToLower(body.Hodnota))
	} else if body.Zmena == "jmeno" {
		if !regexJmeno.MatchString(body.Hodnota) {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Jmeno obsahuje nepovolene znaky nebo ma spatnou delku"))
		}
		err := databaze.PrejmenovatUziv(id, body.Hodnota)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("prázdný request"))
	}
	return c.SendStatus(fiber.StatusOK)
}
