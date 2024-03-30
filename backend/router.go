package main

import (
	"backend/databaze"
	"backend/utils"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"zgo.at/isbot"
)

type (
	bodyDokoncitCvic struct {
		CPM        float32 `json:"cpm" validate:"required"`
		Preklepy   int     `json:"preklepy" validate:"min=0"` //sus reqired nebere nulu takze min=0 asi ok
		Cas        float32 `json:"cas" validate:"required"`
		DelkaTextu int     `json:"delkaTextu" validate:"required"`
	}

	bodyPoslatEmail struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
	}

	bodyRegistrace struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=12"`
		Heslo string `json:"heslo" validate:"required,min=5,max=128"`
	}

	bodyPrihlaseni struct {
		EmailNeboJmeno string `json:"email" validate:"required"`
		Heslo          string `json:"heslo" validate:"required,min=5,max=25"`
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
		Heslo string `json:"heslo" validate:"required,min=5,max=128"`
	}

	bodyGoogle struct {
		AccessToken string `json:"access_token"`
	}

	bodyTestPsani struct {
		Typ   string `json:"typ" validate:"required"`
		Delka int    `json:"delka" validate:"min=1,max=200"`
	}
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/lekce", getVsechnyLekce)
	api.Get("/lekce/:pismena", getCviceniVLekci)
	api.Get("/cvic/:pismena/:cislo", getCviceni)
	api.Post("/dokonceno/:pismena/:cislo", dokoncitCvic)
	api.Get("/procvic", getVsechnyProcvic)
	api.Get("/procvic/:cisloProcvic/:cislo", getProcvic)
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
	api.Get("/testovaci-get-request", test)
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
	/* databaze.PushSlovnik() */
	/* databaze.PushPohadky() */
	return c.JSON("Vypadni pavouku")
}

func testPsani(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}

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
		text, err = databaze.GetVsechnySlova(body.Delka)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}

		for i := 0; i < len(text)-1; i++ {
			text[i] += " "
		}

		for i := 0; i < len(text); i++ {
			r := []rune(text[i])
			if i%5 == 0 { //kazdy paty velkym
				text[i] = fmt.Sprintf("%c%s", unicode.ToUpper(r[0]), string(r[1:]))
			} else {
				text[i] = fmt.Sprintf("%c%s", r[0], string(r[1:]))
			}
		}

	case "vety":
		vety, err := databaze.GetVsechnyVety(body.Delka)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
		}
		for i := 0; i < len(vety); i++ {
			slova := strings.Split(vety[i], " ")
			for _, v := range slova {
				text = append(text, v+" ")
			}
		}
		text[len(text)-1] = text[len(text)-1][:len(text[len(text)-1])-1] // smazat mezeru na konci

	case "nacas":

	default:
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Spatny typ testu psani"))
	}

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "klavesnice": "qwertz"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "klavesnice": u.Klavesnice})
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

func getCviceni(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
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

	var text []string

	switch vsechnyCviceni[cislo-1].Typ {
	case "nova":
		if pismena == "zbylá diakritika" {
			pismena = "óďťň"
		}
		var pismenaRuny []rune = []rune(pismena)

		var slovo strings.Builder
		for i := 0; i < pocetSlov; i++ {
			for j := 0; j < pocetPismenVeSlovu; j++ {
				r := rand.Intn(len(pismenaRuny)) // utf-8 jsou sus
				slovo.WriteRune(pismenaRuny[r])
			}
			slovo.WriteRune(' ')
			text = append(text, slovo.String())
			slovo.Reset()
		}
	case "naucena":
		var naucenaPismena string
		if pismena == "velká písmena (shift)" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťóFJGHDKSLAŮTZRUEIWOQPÚVBCNXMYĚŠČŘŽÝÁÍÉŇĎŤÓ"
		} else if pismena == "zbylá diakritika" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťó"
		} else if pismena == "čísla" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťó1234567890"
		} else {
			naucenaPismena, err = databaze.GetNaucenaPismena(id, pismena)
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
			}
		}

		var slovo strings.Builder
		for i := 0; i < pocetSlov*2/3; i++ { // kratší ať to není taková bolest
			for j := 0; j < pocetPismenVeSlovu; j++ {
				r := rand.Intn(utf8.RuneCountInString(naucenaPismena)) // utf-8 jsou sus
				slovo.WriteRune([]rune(naucenaPismena)[r])
			}
			slovo.WriteRune(' ')
			text = append(text, slovo.String())
			slovo.Reset()
		}
	case "slova":
		var slova []string
		slova, err = databaze.GetSlovaProLekci(id, pismena, pocetSlov)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}

		var pocetSlovKMani int = len(slova)
		var i int = 0
		for {
			if utils.DelkaTextuArray(text) >= delkaTextu-3 { // priblizne idk
				break
			}
			text = append(text, slova[i]+" ")

			i++
			if i >= pocetSlovKMani {
				i = 0
			}
		}

		if pismena == "velká písmena (shift)" { // dam kazdy prvni pismeno velkym
			for i := 0; i < len(text); i++ {
				r := []rune(text[i])
				text[i] = fmt.Sprintf("%c%s", unicode.ToUpper(r[0]), string(r[1:]))
			}
		} else if pismena == "čísla" {
			for i := 1; i < len(text); i += 2 {
				r := rand.Intn(8999) + 1000
				text[i] = strconv.Itoa(r) + " "
			}
		}
	case "programator":
		var slova []string
		slova, err = databaze.GetProgramatorSlova()
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
		}
		var pocetSlovKMani int = len(slova)

		if pismena == "závorky" {
			var zavorky []string = []string{"[\u005D", "()", "{}", "<>"}
			var zavorkyLen = len(zavorky)
			var i int = 0
			var zi int = 0
			rand.Shuffle(zavorkyLen, func(i, j int) { zavorky[i], zavorky[j] = zavorky[j], zavorky[i] })
			for {
				if utils.DelkaTextuArray(text) >= delkaTextu-3 { // priblizne idk
					break
				}
				text = append(text, fmt.Sprintf("%s%v%s ", string([]rune(zavorky[zi])[0]), slova[i], string([]rune(zavorky[zi])[1])))

				i++
				zi++
				if i >= pocetSlovKMani {
					i = 0
				}
				if zi >= zavorkyLen {
					zi = 0
				}
			}
		} else if pismena == "operátory" {
			var oper []string = []string{"=", "==", "!=", "<=", ">=", "<", ">", "+", "-", "*", "/", "%", "+=", "-=", "*=", "/="}
			var operLen = len(oper)
			var i int = 0
			var zi int = 0
			rand.Shuffle(operLen, func(i, j int) { oper[i], oper[j] = oper[j], oper[i] })
			text = append(text, slova[pocetSlovKMani-1]+" ")
			for {
				if utils.DelkaTextuArray(text) >= delkaTextu-3 { // priblizne idk
					break
				}
				text = append(text, fmt.Sprintf("%s %v ", string([]rune(oper[zi])), slova[i]))

				i++
				zi++
				if i >= pocetSlovKMani {
					i = 0
				}
				if zi >= operLen {
					zi = 0
				}
			}
		}
	default:
		log.Print("Cviceni ma divnej typ")
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}

	var posledni bool = int(cislo-1) == len(vsechnyCviceni)-1

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

	err = databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.CPM, body.Preklepy, body.Cas, body.DelkaTextu)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	return c.SendStatus(fiber.StatusOK)
}

func getVsechnyProcvic(c *fiber.Ctx) error {
	texty, err := databaze.GetTexty()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	return c.JSON(fiber.Map{"texty": texty})
}

func getProcvic(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	var klavesnice string
	if id == 0 { // neni prihlaseny
		klavesnice = "qwertz"
	} else {
		u, err := databaze.GetUzivByID(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
		klavesnice = u.Klavesnice
	}
	cisloProcvic, err := strconv.Atoi(c.Params("cisloProcvic")) // str -> int
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}
	cislo := c.Params("cislo") // str -> int
	nazev, text, err := databaze.GetProcvicovani(cisloProcvic, cislo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"text": text, "jmeno": nazev, "klavesnice": klavesnice})
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
		go databaze.SmazatPoLimitu()
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Cas pro overeni vyprsel. Zkus to prosim znovu 1"))
	}

	// Timing attack: nebudu porovnávat stringy ale inty
	kodInt, err := strconv.Atoi(cekajiciUziv.Kod)
	kodInt2, err2 := strconv.Atoi(body.Kod)
	if err != nil || err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(chyba("Divny kod"))
	}

	if time.Now().Unix() <= cekajiciUziv.Cas && kodInt != kodInt2 { // vsechno dobry ale spatnej kod
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
	go utils.MobilNotifikace(cekajiciUziv.Jmeno + " - " + body.Email)
	go databaze.OdebratOvereni(cekajiciUziv.Email)
	go databaze.SmazatPoLimitu()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

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
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Invalidni email"))
	}
	if _, err := databaze.GetUzivByEmail(body.Email); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s timto emailem jiz existuje"))
	}
	if _, err := databaze.GetUzivByJmeno(body.Jmeno); err == nil { // uz existuje
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Uzivatel s timto jmenem jiz existuje"))
	}
	if !regexJmeno.MatchString(body.Jmeno) {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("Invalidni jmeno"))
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
			return c.Status(fiber.StatusUnauthorized).JSON(chyba("Účet je registrován přes google"))
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

func google(c *fiber.Ctx) error {
	var body bodyGoogle
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
	}

	email, jmeno, err := utils.GoogleTokenNaData(body.AccessToken)
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
		go utils.MobilNotifikace(jmeno + " - " + email + " (google)")
	} else {
		token, err = utils.GenerovatToken(email, uziv.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(chyba("Token se pokazil"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

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

	if time.Now().Unix() <= uziv.Cas && uziv.Kod != body.Kod { // vsechno dobry ale spatnej kod
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
	preklepy, cpm, daystreak, cas, delkaVsechTextu, err := databaze.GetUdaje(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(chyba(""))
	}
	uspesnost := float32((delkaVsechTextu - preklepy))
	if uspesnost < 0 {
		uspesnost = 0 // kvuli adamovi kterej big troulin a měl -10%
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"email":           uziv.Email,
		"jmeno":           uziv.Jmeno,
		"daystreak":       daystreak,
		"uspesnost":       uspesnost / float32(delkaVsechTextu) * 100,
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
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"jePotrebaVymenit": true})
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

func navsteva(c *fiber.Ctx) error {
	var httpRequest http.Request
	err := fasthttpadaptor.ConvertRequest(c.Context(), &httpRequest, false)
	if !isbot.Is(isbot.Bot(&httpRequest)) && err == nil {
		databaze.NovaNavsteva()
	}
	return c.SendStatus(fiber.StatusOK)
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
		if !regexJmeno.MatchString(body.Hodnota) {
			return c.Status(fiber.StatusBadRequest).JSON(chyba("Jmeno obsahuje nepovolene znaky nebo ma spatnou delku"))
		}
		err = databaze.PrejmenovatUziv(id, body.Hodnota)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(chyba(err.Error()))
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(chyba("prázdný request"))
	}
	return c.SendStatus(fiber.StatusOK)
}
