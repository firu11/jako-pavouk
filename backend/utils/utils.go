package utils

import (
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

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
func Prumer(arr []float64) float64 {
	var soucet float64 = 0
	for _, v := range arr {
		soucet += v
	}
	if len(arr) == 0 {
		return 0
	}
	return soucet / float64(len(arr))
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

// vzorec pro výpočet ryhclosti (založen na pravidlech státní zkoušky)
func CPM(delkaTextu int, cas float64, preklepy int) float64 {
	var cpm float64 = (float64(delkaTextu-10*preklepy) / cas) * 60
	if cpm < 0 {
		return 0
	}
	return cpm
}

// vrací 6ti místný string kód k třídě
func GenTridaKod() string {
	var kod string
	var znaky string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < 6; i++ {
		cislo, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(znaky))))
		if err != nil {
			cislo = big.NewInt(int64(mathRand.Intn(len(znaky)))) // kdyby se něco pokazilo?
		}
		kod += string(znaky[cislo.Int64()])
	}
	return kod
}

func GetRole(uzivRole int, trida_id uint) string {
	if uzivRole == 2 {
		return "ucitel"
	}
	if uzivRole == 1 {
		if trida_id != 0 {
			return "student"
		} else {
			return "basic"
		}
	}
	return "basic"
}
