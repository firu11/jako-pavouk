package main

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

	"github.com/labstack/echo/v4"
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
		ID      uint   `json:"id" validate:"required"`
		Jmeno   string `json:"jmeno" validate:"max=30"`
		TridaID int    `json:"trida_id"`
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
		TridaID uint    `json:"trida_id"`
		Typ     string  `json:"typ" validate:"required"`
		ZLekce  string  `json:"z_lekce"`
		Delka   float32 `json:"delka"`
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
func setupSkolniRouter(api *echo.Group) {
	skolaApi := api.Group("/skola")

	skolaApi.POST("/create-trida", createTrida)
	skolaApi.GET("/tridy", tridy)
	skolaApi.GET("/trida", tridaStudent)
	skolaApi.GET("/trida/:id", trida)
	skolaApi.GET("/zaci-stream/:id", zaciStream)
	skolaApi.GET("/test-tridy/:kod", testTridy)
	skolaApi.POST("/zmena-tridy", zmenaTridy)

	skolaApi.POST("/pridat-praci", pridatPraci)
	skolaApi.GET("/get-praci/:id", getPraci)
	skolaApi.POST("/dokoncit-praci/:id", dokoncitPraci)
	skolaApi.DELETE("/smazat-praci/:id", smazatPraci)

	skolaApi.POST("/text", getText)
	skolaApi.GET("/typy-cviceni", getTypyCviceni)

	skolaApi.GET("/student/:id", student)
	skolaApi.POST("/student", studentUprava)
	skolaApi.POST("/zapis", zapis)
}

var kanalyTrid = make(map[int]chan int)

func zaciStream(c echo.Context) error {
	tridaID, err := strconv.ParseInt(c.Param("id"), 10, 8)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	kanalyTrid[int(tridaID)] = make(chan int)

	fmt.Fprint(w, "Jedemee\n\n")
	w.Flush()

	for {
		select {
		case <-c.Request().Context().Done():
			delete(kanalyTrid, int(tridaID))
			return nil
		case <-kanalyTrid[int(tridaID)]:
			fmt.Fprint(w, "data: Nový žák\n\n")
			w.Flush()
		}
	}
}

func createTrida(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tridu muze vytvaret pouze ucitel"))
	}

	var body bodyCreateTrida
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("body musí mít jméno třídy"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	for i := 0; i < 5; i++ {
		err = databaze.CreateTrida(body.Jmeno, id, utils.GenTridaKod())
		if err == nil {
			return c.NoContent(http.StatusOK)
		}
		// kod asi neni unikatni tak zkusim znovu
	}

	// několikrat se to pokazilo tak hodim error
	return c.JSON(http.StatusBadRequest, chyba(err.Error()))
}

func tridy(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusUnauthorized, chyba("Tridy muze videt pouze ucitel"))
	}

	tridy, err := databaze.GetTridy(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
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

	return c.JSON(http.StatusOK, echo.Map{"tridy": output})
}

func trida(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tridy muze videt pouze ucitel"))
	}

	tridaID, err := strconv.ParseInt(c.Param("id"), 10, 8)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	ostatniTridy, err := databaze.GetTridy(uziv.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	trida, err := databaze.GetTrida(uint(tridaID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	if trida.UcitelID != uziv.ID {
		return c.JSON(http.StatusBadRequest, chyba("To neni tvoje trida"))
	}
	studenti, err := databaze.GetStudentyZeTridy(trida.ID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	prace, err := databaze.GetVsechnyPrace(trida.ID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	return c.JSON(http.StatusOK, echo.Map{"trida": trida, "studenti": studenti, "prace": prace, "ostatniTridy": ostatniTridy})
}

func tridaStudent(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	trida, err := databaze.GetTridaByStudentID(uziv.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	prace, err := databaze.GetVsechnyPrace(trida.ID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	cpmka, presnost, err := databaze.GetDokoncenePrace(uziv.ID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
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

	return c.JSON(http.StatusOK, echo.Map{"trida": trida, "prace": vysledek, "klavesnice": uziv.Klavesnice})
}

func testTridy(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	_, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	trida, err := databaze.GetTridaByKod(strings.ToUpper(c.Param("kod")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba("Takova trida neexistuje"))
	}
	if trida.Zamknuta {
		return c.JSON(http.StatusBadRequest, chyba("Trida je zamcena"))
	}

	return c.NoContent(http.StatusOK)
}

func zmenaTridy(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tridy muze upravovat pouze ucitel"))
	}

	var body bodyZmenaTridy
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}

	switch body.Zmena {
	case "zamek":
		err := databaze.ZamknoutTridu(body.TridaID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}
	case "jmeno":
		if body.Hodnota != "" {
			err := databaze.PrejmenovatTridu(body.TridaID, body.Hodnota)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
			}
		} else {
			return c.JSON(http.StatusBadRequest, chyba("Potrebuju hodnotu"))
		}
	case "smazat":
		err := databaze.SmazatTridu(body.TridaID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}
	case "klavesnice":
		err := databaze.ZmenitKlavesniciTridy(body.TridaID, body.Hodnota)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}
	}

	return c.NoContent(http.StatusOK)
}

func student(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tohle muze pouze ucitel"))
	}

	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	student, err := databaze.GetUzivByID(uint(studentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}
	presnost, cpm, daystreak, chybyPismenka, err := databaze.GetUdaje(student.ID)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	dokonceno, err := databaze.DokonceneProcento(student.ID)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	cpmkaVPracich, presnostVPracich, err := databaze.GetDokoncenePrace(student.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"student": echo.Map{
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

func studentUprava(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tohle muze pouze ucitel"))
	}

	var body bodyZmenaJmenaStudenta
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}

	if body.Jmeno != "" {
		err = databaze.PrejmenovatStudenta(body.ID, body.Jmeno)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
	}
	if body.TridaID == -1 {
		err = databaze.OdebratStudenta(body.ID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
	}
	if body.TridaID != 0 {
		err = databaze.PresunoutStudenta(body.ID, body.TridaID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
	}

	return c.NoContent(http.StatusOK)
}

func zapis(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	var body bodyZapis
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}

	tridaID, err := databaze.ZapsatStudenta(body.Kod, id, body.Jmeno)
	if err != nil {
		if err.Error() == "uz je ve tride" {
			return c.JSON(http.StatusBadRequest, chyba("Uz jsi ve tride"))
		}
		if err.Error() == "jako ucitel nemuzete byt ve tride" {
			return c.JSON(http.StatusBadRequest, chyba(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	ch, ok := kanalyTrid[tridaID]
	if ok {
		select {
		case ch <- 1:
		default:
		}
	}

	return c.NoContent(http.StatusOK)
}

func pridatPraci(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tohle muze pouze ucitel"))
	}

	var body bodyPridatPraci
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}

	if err = databaze.PridatPraci(body.Text, body.Cas, body.TridaID); err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	return c.NoContent(http.StatusOK)
}

func smazatPraci(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tohle muze pouze ucitel"))
	}

	praceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}
	err = databaze.SmazatPraci(uint(praceID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	return c.NoContent(http.StatusOK)
}

func getText(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	if uziv.Role != 2 {
		return c.JSON(http.StatusBadRequest, chyba("Tohle muze pouze ucitel"))
	}

	var body bodyGetText
	if err := c.Bind(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}
	if err := utils.ValidateStruct(&body); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba("Spatny body"))
	}

	var text string = ""
	if body.Typ == "1" {
		textPole, err := generovatTextCviceni(body.ZLekce, "nova", 0, body.TridaID, body.Delka+5)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
		text = strings.Join(textPole, "")
	} else if body.Typ == "2" {
		textPole, err := generovatTextCviceni(body.ZLekce, "naucena", 0, body.TridaID, body.Delka+5)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
		text = strings.Join(textPole, "")
	} else if body.Typ == "3" {
		textPole, err := generovatTextCviceni(body.ZLekce, "slova", 0, body.TridaID, body.Delka+5)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
		text = strings.Join(textPole, "")
	} else if body.Typ == "4" {
		textPole, err := generovatTextCviceni(body.ZLekce, "programator", 0, body.TridaID, body.Delka+5)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, chyba(""))
		}
		text = strings.Join(textPole, "")
	} else if body.Typ == "Věty z pohádek" {
		var textBuilder strings.Builder
		vety, err := databaze.GetVsechnyVety(int(pocetZnaku / 85)) // cca 85 znaku na vetu
		if err != nil {
			return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
		}

		for _, v := range vety {
			textBuilder.WriteString(v)
			textBuilder.WriteString(" ")
		}

		text = textBuilder.String()
	} else {
		var textBuilder strings.Builder
		txt, err := databaze.GetRandomProcvic(body.Typ)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, chyba("Takový typ nemáme"))
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, chyba(err.Error()))
		}

		textBuilder.WriteString(txt)
		text = textBuilder.String()
	}

	if len(text) != 0 && string(text[len(text)-1]) == " " {
		text = text[:len(text)-1]
	}

	return c.JSON(http.StatusOK, echo.Map{"text": text})
}

func getPraci(c echo.Context) error {
	id := c.Get("uzivID").(uint)
	if id == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	praceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, chyba(err.Error()))
	}

	text, cas, err := databaze.GetPrace(uint(praceID), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
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

	return c.JSON(http.StatusOK, echo.Map{"text": vyslednyText, "cas": cas, "klavesnice": uziv.Klavesnice})
}

func dokoncitPraci(c echo.Context) error {
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
	praceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	tridaID, err := databaze.DokoncitPraci(uint(praceID), id, body.Preklepy, body.Cas, body.DelkaTextu, body.NejcastejsiChyby)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(err.Error()))
	}

	ch, ok := kanalyTrid[tridaID]
	if ok {
		select {
		case ch <- 1:
		default:
		}
	}

	return c.NoContent(http.StatusOK)
}

func getTypyCviceni(c echo.Context) error {
	tridaID, err := strconv.Atoi(c.QueryParam("trida_id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}

	mapa, err := databaze.GetTypyCviceni(uint(tridaID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, chyba(""))
	}
	return c.JSON(http.StatusOK, mapa)
}
