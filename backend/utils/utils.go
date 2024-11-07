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
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-playground/validator/v10"
)

var verifier = emailverifier.NewVerifier()

func ValidaceEmailu(email string) error {
	ret, err := verifier.Verify(email)
	if err != nil {
		return err
	}
	if !ret.Syntax.Valid {
		return errors.New("email ma spatny syntax")
	}
	return nil
}

func ValidFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// přetvoří request body do požadovaného structu
var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	return err
}

// id = 0 znamena ze se neco pokazilo
func Autentizace(tokenHeader string) uint {
	if len(tokenHeader) < 7 {
		return 0
	}
	var token string = tokenHeader[7:]
	spravnej, id, err := ValidovatToken(token)

	if spravnej && err == nil {
		return id
	} else {
		return 0
	}
}

func Prumer(arr []float64) float64 {
	var soucet float64 = 0
	for _, v := range arr {
		soucet += v
	}
	if len(arr) == 0 {
		return -1
	}
	return soucet / float64(len(arr))
}

// ošetřuje escape charaktery v url (%C5%A1 -> š)
func DecodeURL(s string) (string, error) {
	x, err := url.QueryUnescape(s)
	if err != nil {
		log.Print(err)
		return "", err
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

// pošle mi na telefon notigikaci
func MobilNotifikace(s string) {
	http.Post(os.Getenv("MOBIL_NOTIFIKACE_URL"), "text/plain", strings.NewReader(s))
}

// počítá délku textu z pole ["slovo ", "slovo "]
func DelkaTextuArray(a []string) int {
	var x int
	for _, v := range a {
		x += utf8.RuneCountInString(v)
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

// vrací 4místný string kód k třídě
func GenTridaKod() string {
	var kod string
	var znaky string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 4; i++ {
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

func SmazatMezeruNaKonci(text []string) {
	posledniSlovo := text[len(text)-1]
	if string(posledniSlovo[len(posledniSlovo)-1]) == " " {
		text[len(text)-1] = text[len(text)-1][:len(posledniSlovo)-1]
	}
}

func SaveSkola(jmenoSkoly string, kontaktniEmail string, kontaktniTelefon string) error {
	f, err := os.OpenFile("skoly.log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(fmt.Sprintf("%s, %s, %s\n", jmenoSkoly, kontaktniEmail, kontaktniTelefon)))
	return err
}

var PismenaPobliz map[rune]rune = map[rune]rune{
	't': 'f',
	'r': 'f',
	'ř': 'f',
	'ž': 'f',
	'v': 'f',
	'b': 'f',
	'z': 'j',
	'u': 'j',
	'á': 'j',
	'ý': 'j',
	'n': 'j',
	'm': 'j',
	'č': 'd',
	'e': 'd',
	'c': 'd',
	'š': 's',
	'w': 's',
	'x': 's',
	'ě': 'a',
	'+': 'a',
	'q': 'a',
	'y': 'a',
	'í': 'k',
	'i': 'k',
	',': 'k',
	'é': 'l',
	'o': 'l',
	'.': 'l',
	'=': 'ů',
	'p': 'ů',
	'-': 'ů',
	'ú': 'ů',
	'(': 'ů',
	')': 'ů',
	'/': 'ů',
	'1': 'a',
	'2': 'a',
	'3': 's',
	'4': 'd',
	'5': 'f',
	'6': 'f',
	'7': 'j',
	'8': 'j',
	'9': 'k',
	'0': 'l',
	'ó': 'l',
	'ď': 'd',
	'ť': 'f',
	'ň': 'j',
}
