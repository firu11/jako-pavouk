package handlers

import (
	"backend/databaze"
	"backend/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/labstack/echo/v4"
	"zgo.at/isbot"
)

func testPsani(c echo.Context) error {
	id := c.Get("uzivID").(uint)

	var body = bodyTestPsani{}
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	var text []string
	switch body.Typ {
	case "slova":
		var err error
		text, err = databaze.GetVsechnySlova(int(PocetZnaku / 7.5)) // cca 8.5 znaku na slovo
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}

		for i := 0; i < len(text); i++ {
			text[i] += " "

			r := []rune(text[i])
			if i%5 == 0 { //kazdy paty velkym
				text[i] = fmt.Sprintf("%c%s", unicode.ToUpper(r[0]), string(r[1:]))
			}
		}

	case "vety":
		vety, err := databaze.GetVsechnyVety(int(PocetZnaku / 85)) // cca 85 znaku na vetu
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}
		for i := range vety {
			slova := strings.Split(vety[i], " ")
			for _, v := range slova {
				text = append(text, v+" ")
			}
		}

	default:
		return c.JSON(http.StatusBadRequest, chyba("Spatny typ testu psani"))
	}

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{"text": text})
	}
	return c.JSON(http.StatusOK, echo.Map{"text": text, "klavesnice": u.Klavesnice})
}

func getVsechnyLekce(c echo.Context) error {
	id := c.Get("uzivID").(uint)

	lekce, err := databaze.GetLekce(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	if id != 0 {
		dokoncene, err := databaze.GetDokonceneLekce(id)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}

		lekcePismena, cisloCvic, err := databaze.GetDalsiCviceni(id)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Println(err)
			}
			return c.JSON(http.StatusOK, echo.Map{"lekce": lekce, "dokoncene": dokoncene, "dalsi_cviceni": ""})
		}
		return c.JSON(http.StatusOK, echo.Map{"lekce": lekce, "dokoncene": dokoncene, "dalsi_cviceni": fmt.Sprintf("/%s/%d", lekcePismena, cisloCvic)})
	}
	return c.JSON(http.StatusOK, echo.Map{"lekce": lekce, "dokoncene": []int{}, "dalsi_cviceni": ""})
}

func getCviceniVLekci(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	pismena, HTTPerr := utils.DecodeURL(c.Param("pismena"))
	if HTTPerr != nil {
		return c.JSON(http.StatusBadRequest, chyba(HTTPerr.Error()))
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba("Takova lekce neexistuje"))
	}
	doko, err := databaze.GetDokonceneCvicVLekci(id, 0, pismena)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	return c.JSON(http.StatusOK, echo.Map{"cviceni": cvic, "dokoncene": doko})
}

func getCviceni(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	pismena, HTTPerr := utils.DecodeURL(c.Param("pismena"))
	if HTTPerr != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	cislo, err := strconv.Atoi(c.Param("cislo")) // str -> int
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if cislo > len(vsechnyCviceni) {
		return c.JSON(http.StatusBadRequest, chyba("Cviceni neexistuje"))
	}

	text, _ := utils.GenerovatTextCviceni(pismena, vsechnyCviceni[cislo-1].Typ, id, 0, PocetZnaku)

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	return c.JSON(http.StatusOK, echo.Map{"text": text, "klavesnice": u.Klavesnice, "typ": vsechnyCviceni[cislo-1].Typ, "posledni": int(cislo-1) == len(vsechnyCviceni)-1})
}

func dokoncitCvic(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	var body = bodyDokoncit{}

	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	cislo, err := strconv.ParseUint(c.Param("cislo"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	pismena, HTTPerr := utils.DecodeURL(c.Param("pismena"))
	if HTTPerr != nil {
		return HTTPerr
	}
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(id, pismena)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if int(cislo-1) >= len(vsechnyCviceni) { // error index out of range nebude
		log.Print("Takovy cviceni neni")
		return c.NoContent(http.StatusBadRequest)
	}

	err = databaze.PridatDokonceneCvic(uint(vsechnyCviceni[cislo-1].ID), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	return c.NoContent(http.StatusOK)
}

func dokoncitProcvic(c echo.Context) error {
	id := c.Get("uzivID").(uint)

	var body = bodyDokoncit{}

	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	cislo, err := strconv.ParseUint(c.Param("cislo"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	vsechnyProcvic, err := databaze.GetTexty()
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if cislo > uint64(len(vsechnyProcvic)) { // error index out of range nebude
		log.Print("Takovy procvicovani neni")
		return c.NoContent(http.StatusBadRequest)
	}
	err = databaze.PridatDokonceneProcvic(uint(cislo), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	return c.NoContent(http.StatusOK)
}

func getVsechnyProcvic(c echo.Context) error {
	id := c.Get("uzivID").(uint)

	texty, err := databaze.GetTexty()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	if id == 0 {
		for i := range texty {
			texty[i].CPM = -1
		}
		return c.JSON(http.StatusOK, echo.Map{"texty": texty, "testPsaniCPM": -1})
	}

	rychlosti, err := databaze.GetRychlostiProcvic(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	for i := range texty {
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

	return c.JSON(http.StatusOK, echo.Map{"texty": texty, "testPsaniCPM": testPsaniCPM})
}

func getProcvic(c echo.Context) error {
	id := c.Get("uzivID").(uint)

	typ, err := strconv.Atoi(c.Param("cisloProcvic")) // str -> int
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	cislo, err := strconv.Atoi(c.Param("neCislo")) // str -> int
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	nazev, podnazev, text, cislo, err := databaze.GetProcvicovani(typ, cislo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	u, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{"text": text, "typ": nazev, "jmeno": podnazev, "cislo": cislo})
	}

	return c.JSON(http.StatusOK, echo.Map{"text": text, "typ": nazev, "jmeno": podnazev, "cislo": cislo, "klavesnice": u.Klavesnice})
}

func overitEmail(c echo.Context) error {
	var body bodyPoslatEmail = bodyPoslatEmail{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.JSON(http.StatusBadRequest, chyba("Invalidni email"))
	}

	cekajiciUziv, err := databaze.GetNeoverenyUziv(body.Email)
	if err != nil {
		go databaze.SmazatPoLimitu()
		return c.JSON(http.StatusBadRequest, chyba("Cas pro overeni vyprsel. Zkus to prosim znovu 1"))
	}

	if time.Now().Unix() <= cekajiciUziv.Cas && !utils.CheckKod(cekajiciUziv.Kod, body.Kod) { // vsechno dobry ale spatnej kod
		databaze.DalSpatnyKod(body.Email)
		return c.JSON(http.StatusBadRequest, chyba("Spatny kod"))
	} else if time.Now().Unix() > cekajiciUziv.Cas { // vyprselo
		return c.JSON(http.StatusBadRequest, chyba("Cas pro overeni vyprsel. Zkus to prosim znovu 2"))
	}

	uzivID, err := databaze.CreateUziv(cekajiciUziv.Email, cekajiciUziv.Heslo, cekajiciUziv.Jmeno)
	if err != nil {
		log.Println(err, uzivID)
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	token, err := utils.GenerovatToken(body.Email, uzivID)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	go databaze.OdebratOvereni(cekajiciUziv.Email)
	go databaze.SmazatPoLimitu()
	return c.JSON(http.StatusOK, echo.Map{"token": token, "jmeno": cekajiciUziv.Jmeno, "email": cekajiciUziv.Email})
}

func registrace(c echo.Context) error {
	var body bodyRegistrace

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	body.Email = strings.ToLower(body.Email)
	if err := utils.ValidaceEmailu(body.Email); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(fmt.Sprintf("Invalidni email: %s", err)))
	}
	if _, err := databaze.GetUzivByEmail(body.Email); err == nil { // uz existuje
		return c.JSON(http.StatusBadRequest, chyba("Uzivatel s timto emailem jiz existuje"))
	}

	if !RegexJmeno.MatchString(body.Jmeno) {
		return c.JSON(http.StatusBadRequest, chyba("Invalidni jmeno"))
	}
	if _, err := databaze.GetUzivByJmeno(body.Jmeno); err == nil { // uz existuje
		return c.JSON(http.StatusBadRequest, chyba("Uzivatel s timto jmenem jiz existuje"))
	}

	hesloHASH, err := utils.HashPassword(body.Heslo)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	var randomKod string = utils.GenKod()

	if err := databaze.CreateNeoverenyUziv(body.Email, hesloHASH, body.Jmeno, randomKod, time.Now().Add(10*time.Minute).Unix()); err != nil {
		go databaze.SmazatPoLimitu()
		return c.JSON(http.StatusBadRequest, chyba("Uzivatel s timto jmenem docasne existuje"))
	}
	if err := utils.PoslatOverovaciEmail(body.Email, randomKod); err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	go databaze.SmazatPoLimitu()
	return c.NoContent(http.StatusOK)
}

func prihlaseni(c echo.Context) error {
	var body bodyPrihlaseni

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	var uziv databaze.Uzivatel

	if !utils.ValidFormat(body.EmailNeboJmeno) { //predpokladam ze je to jmeno kdyz se to nepodobá emailu
		uziv, err = databaze.GetUzivByJmeno(body.EmailNeboJmeno)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chyba("Email/jmeno je spatne (Jmeno)"))
		}
	} else {
		uziv, err = databaze.GetUzivByEmail(body.EmailNeboJmeno)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chyba("Email/jmeno je spatne (Email)"))
		}
	}

	if err := utils.CheckPassword(body.Heslo, uziv.Heslo); err != nil {
		if err.Error() == "ucet je pres google" {
			return c.JSON(http.StatusUnauthorized, chyba("Ucet je registrovan pres google"))
		}
		return c.JSON(http.StatusUnauthorized, chyba("Heslo je spatne"))
	} else {
		token, err := utils.GenerovatToken(uziv.Email, uziv.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba("Token se pokazil"))
		} else {
			trida, _ := databaze.GetTridaByUziv(uziv.ID)
			return c.JSON(http.StatusOK, echo.Map{"token": token, "jmeno": uziv.Jmeno, "email": uziv.Email, "role": utils.GetRole(uziv.UcitelVeSkoleID, trida.ID)})
		}
	}
}

func google(c echo.Context) error {
	var body bodyGoogle
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	email, jmeno, err := utils.GoogleTokenNaData(body.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	var token string
	var novy bool = false
	uziv, err := databaze.GetUzivByEmail(email)
	if err != nil { // neexistuje
		novy = true
		id, err := databaze.CreateUziv(email, "google", jmeno)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chyba(err.Error()))
		}
		token, err = utils.GenerovatToken(email, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba("Token se pokazil"))
		}

		uziv, err = databaze.GetUzivByEmail(email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
	} else {
		token, err = utils.GenerovatToken(email, uziv.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba("Token se pokazil"))
		}
	}

	trida, _ := databaze.GetTridaByUziv(uziv.ID)
	return c.JSON(http.StatusOK, echo.Map{"token": token, "novy": novy, "jmeno": uziv.Jmeno, "email": uziv.Email, "role": utils.GetRole(uziv.UcitelVeSkoleID, trida.ID)})
}

func zmenaHesla(c echo.Context) error {
	var body bodyZmenaHesla
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	_, err = databaze.GetUzivByEmail(body.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba("Email je spatne"))
	}

	kod := utils.GenKod()
	if err := utils.PoslatOverovaciEmail(body.Email, kod); err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	if err = databaze.CreateZapomenuteHeslo(body.Email, kod, time.Now().Add(10*time.Minute).Unix()); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	go databaze.SmazatPoLimitu()

	return c.NoContent(http.StatusOK)
}

func overitZmenuHesla(c echo.Context) error {
	var body bodyOvereniZmenaHesla
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	err := utils.ValidateStruct(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	uziv, err := databaze.GetZmenuHesla(body.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	hesloHASH, err := utils.HashPassword(body.Heslo)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	if time.Now().Unix() <= uziv.Cas && !utils.CheckKod(uziv.Kod, body.Kod) { // vsechno dobry ale spatnej kod
		return c.JSON(http.StatusBadRequest, chyba("Spatny kod"))
	} else if time.Now().Unix() > uziv.Cas { // vyprselo
		return c.JSON(http.StatusBadRequest, chyba("Cas pro overeni vyprsel. Zkuste to prosim znovu"))
	}

	err = databaze.ZmenitHeslo(uziv.Email, hesloHASH)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	go databaze.OdebratZmenuHesla(uziv.Email)
	go databaze.SmazatPoLimitu()

	return c.NoContent(http.StatusOK)
}

func nastaveni(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	trida, _ := databaze.GetTridaByUziv(uziv.ID)
	return c.JSON(http.StatusOK, echo.Map{
		"id":         uziv.ID,
		"email":      uziv.Email,
		"jmeno":      uziv.Jmeno,
		"role":       utils.GetRole(uziv.UcitelVeSkoleID, trida.ID),
		"klavesnice": uziv.Klavesnice,
	})
}

func statistiky(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	daystreak, err := databaze.GetDaystreak(id)
	if err != nil {
		log.Print(err, 1)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	presnost, rychlost, chybyPismenka, cas, napsanychPismen, err := databaze.GetUdaje(id)
	if err != nil {
		log.Print(err, 2)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(id)
	if err != nil {
		log.Print(err, 3)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	rychlosti, presnosti, err := databaze.GetUdajeProGraf(id)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"daystreak":        daystreak,
		"postupVKurzu":     dokonceno,
		"uspesnost":        presnost,
		"rychlost":         rychlost,
		"cas":              cas,
		"nejcastejsiChyby": chybyPismenka,
		"napsanychPismen":  napsanychPismen,
		"rychlosti":        rychlosti,
		"presnosti":        presnosti,
	})
}

func testVyprseniTokenu(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if len(token) < 10 { // treba deset proste at tam neco je
		return c.JSON(http.StatusUnauthorized, chyba(""))
	}
	jePotrebaVymenit, err := utils.ValidovatExpTokenu(token[7:])
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{"jePotrebaVymenit": true})
	}
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil && !jePotrebaVymenit {
		jePotrebaVymenit = true
	}
	trida, _ := databaze.GetTridaByUziv(uziv.ID)
	return c.JSON(http.StatusOK, echo.Map{"jmeno": uziv.Jmeno, "email": uziv.Email, "jePotrebaVymenit": jePotrebaVymenit, "role": utils.GetRole(uziv.UcitelVeSkoleID, trida.ID)})
}

func navsteva(c echo.Context) error {
	if !isbot.Is(isbot.Bot(c.Request())) {
		databaze.NovaNavsteva()
	}
	return c.NoContent(http.StatusOK)
}

func upravaUctu(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	var body = bodyUprava{}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	if body.Zmena == "smazat" {
		if err := databaze.SmazatUzivatele(id); err != nil {
			return err
		}
	} else if body.Zmena == "klavesnice" {
		databaze.ZmenitKlavesnici(id, strings.ToLower(body.Hodnota))
	} else if body.Zmena == "jmeno" {
		if !RegexJmeno.MatchString(body.Hodnota) {
			return c.JSON(http.StatusBadRequest, chyba("Jmeno obsahuje nepovolene znaky nebo ma spatnou delku"))
		}
		err := databaze.PrejmenovatUziv(id, body.Hodnota)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chyba(err.Error()))
		}
	} else {
		return c.JSON(http.StatusBadRequest, chyba("prázdný request"))
	}
	return c.NoContent(http.StatusOK)
}
