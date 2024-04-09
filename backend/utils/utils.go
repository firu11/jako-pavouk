package utils

import (
	"backend/databaze"
	cryptoRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	godiacritics "gopkg.in/Regis24GmbH/go-diacritics.v2"
)

var RegexJmeno *regexp.Regexp
var MaxCisloZaJmeno int // 10_000

// validuje email
func ValidFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var validate = validator.New()

// přetvoří request body do požadovaného structu
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	return err
}

/*
Vraci id, error.
id = 0 znamena ze se bud neco pokazilo nebo je autentizace nepovinna
*/
func Autentizace(c *fiber.Ctx, povinna bool) (uint, error) {
	if len(c.Get("Authorization")) >= 10 { // treba deset proste at tam neco je
		var token string = c.Get("Authorization")[7:]
		spravnej, id, err := ValidovatToken(token)

		if spravnej && err == nil {
			return id, nil
		} else if !spravnej {
			if povinna {
				return 0, errors.New("spatny token")
			} else {
				return 0, nil
			}
		} else {
			return 0, err
		}
	} else {
		if povinna {
			return 0, errors.New("je potreba autentizace (JWT Token)")
		}
		return 0, nil
	}
}

// vrací průměr floatů z pole
func Prumer(arr []float32) float32 {
	var soucet float32 = 0
	for _, v := range arr {
		soucet += v
	}
	if len(arr) == 0 {
		return 0
	}
	return soucet / float32(len(arr))
}

// ošetřuje escape charaktery v url (%C5%A1 -> š)
func DecodeURL(s string) (string, error) {
	x, err := url.QueryUnescape(s)
	if err != nil {
		log.Print(err)
		return "", fiber.ErrBadRequest
	}
	return x, nil
}

// vrací 5ti místný string kód
func GenKod() string {
	var kod string
	for i := 0; i < 5; i++ {
		cislo, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(10))
		if err != nil {
			cislo = big.NewInt(int64(mathRand.Intn(10))) // kdyby se něco pokazilo?
		}
		kod += fmt.Sprintf("%v", cislo)
	}
	return kod
}

// porovná ověřovací kód
func CheckKod(kod1 string, kod2 string) bool {
	// Timing attack: nebudu porovnávat stringy ale inty
	kodInt, err := strconv.Atoi(kod1)
	kodInt2, err2 := strconv.Atoi(kod2)
	if err != nil || err2 != nil {
		return false
	}
	return kodInt == kodInt2
}

// pošle mi na telefon notigikaci, chci vědět když se někdo zaregistruje :)
func MobilNotifikace(s string) {
	http.Post("https://ntfy.sh/novy_uzivatel115115jakopavouk", "text/plain", strings.NewReader(s))
}

// počítá délku textu z pole ["slovo ", "slovo "]
func DelkaTextuArray(a []string) int {
	var x int
	for _, v := range a {
		x += len(v)
	}
	return x
}

// z googlu vrací email, jmeno, error
func GoogleTokenNaData(token string) (string, string, error) {
	res, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%v", token))
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	m := make(map[string]string)
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", "", err
	}

	if m["aud"] != os.Getenv("GOOGLE_CLIENT_ID") {
		return "", "", errors.New("fake token")
	}

	jmeno, err := volbaJmena(m["name"])
	if err != nil {
		return "", "", err
	}

	return m["email"], jmeno, err
}

// vybírá jméno pro uživatele který se zaregistroval přes google
//
// zkusí kombinace google jména a náhodného čísla, poté Pavouk a číslo
//
// číslo přidávám k jménu abych minimalizoval šanci, že takový uživatel již existuje a musím vytvářet nové jméno a znovu kontrolovat v db
func volbaJmena(celeJmeno string) (string, error) {
	celeJmeno = godiacritics.Normalize(celeJmeno)
	var jmeno []string = strings.Fields(celeJmeno) // rozdělim na jmeno a prijimeni

	for range 20 { // vic než 20x to zkoušet nebudu
		var cislo int = mathRand.Intn(MaxCisloZaJmeno-1) + 1

		var jmenoNaTest string
		if len(jmeno) >= 1 {
			jmenoNaTest = fmt.Sprintf("%s%d", jmeno[0], cislo)
			if RegexJmeno.MatchString(jmenoNaTest) {
				_, err := databaze.GetUzivByJmeno(jmenoNaTest)
				if err != nil {
					return jmenoNaTest, nil
				}
			}
		}
		if len(jmeno) == 2 {
			jmenoNaTest = fmt.Sprintf("%s%d", jmeno[1], cislo)
			if RegexJmeno.MatchString(jmenoNaTest) {
				_, err := databaze.GetUzivByJmeno(jmenoNaTest)
				if err != nil { // ještě neexistuje
					return jmenoNaTest, nil
				}
			}
		}
		jmenoNaTest = fmt.Sprintf("Pavouk%d", cislo)
		if RegexJmeno.MatchString(jmenoNaTest) {
			_, err := databaze.GetUzivByJmeno(jmenoNaTest)
			if err != nil { // ještě neexistuje
				return jmenoNaTest, nil
			}
		}
	}

	return "", errors.New("konec sveta nenašel jsem jméno")
}
