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
	"regexp"
	"strconv"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// validuje email
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
		return -1
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
	http.Post(os.Getenv("MOBIL_NOTIFIKACE_URL"), "text/plain", strings.NewReader(s))
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

func SmazatMezeruNaKonci(text []string) {
	posledniSlovo := text[len(text)-1]
	if string(posledniSlovo[len(posledniSlovo)-1]) == " " {
		text[len(text)-1] = text[len(text)-1][:len(posledniSlovo)-1]
	}
}

var uvozovka string = `"`
var pomlcka string = `-`
var tabulka = map[rune]string{
	'„': uvozovka,
	'”': uvozovka,
	'“': uvozovka,
	'‟': uvozovka,
	'❞': uvozovka,
	'❝': uvozovka,
	'🙷': uvozovka,
	'🙶': uvozovka,
	'＂': uvozovka,
	'‚': uvozovka,
	'’': uvozovka,
	'‘': uvozovka,
	'‛': uvozovka,
	'❛': uvozovka,
	'❜': uvozovka,
	'⹂': uvozovka,
	'…': "...",
	'‑': pomlcka,
	'‐': pomlcka,
	'־': pomlcka,
	'‒': pomlcka,
	'–': pomlcka,
	'—': pomlcka,
	'―': pomlcka,
	'⸺': pomlcka,
	'⸻': pomlcka,
	'﹘': pomlcka,
	'﹣': pomlcka,
	'－': pomlcka,
	'֊': pomlcka,
	'᠆': pomlcka,
	'ä': "a", // nemecke
	'ö': "o",
	'ü': "u",
	'ć': "c", // polske
	'ą': "a",
	'ę': "e",
	'ł': "l",
	'ń': "n",
	'ó': "o",
	'ź': "z",
	'ż': "z",
	'ĺ': "l", //slovenske
	'ľ': "l",
	'ô': "o",
	'ŕ': "r",
	'è': "e", // ani nevim
	'ã': "a",
	'ë': "e",
	'æ': "ae",
}
var tabulka2 map[string]string = map[string]string{
	"…": "...",
}
var mezery = regexp.MustCompile("( {2,})")

func UpravaTextu(txt string) string {
	var v strings.Builder
	for _, ch := range txt {
		if tabulka[ch] != "" {
			v.WriteString(tabulka[ch])
		} else {
			v.WriteRune(ch)
		}
	}

	for o, n := range tabulka2 {
		txt = strings.ReplaceAll(txt, o, n)
	}

	vysledek := v.String()
	vysledek = mezery.ReplaceAllString(vysledek, " ") // odstrani vice mezer

	return strings.TrimSpace(vysledek)
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
