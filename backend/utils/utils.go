package utils

import (
	"backend/databaze"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	godiacritics "gopkg.in/Regis24GmbH/go-diacritics.v2"
)

var RegexJmeno *regexp.Regexp

func ValidFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var validate = validator.New()
var t transform.Transformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

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

func DecodeURL(s string) (string, error) {
	x, err := url.QueryUnescape(s)
	if err != nil {
		log.Print(err)
		return "", fiber.ErrBadRequest
	}
	return x, nil
}

func GenKod() string {
	var kod string
	for i := 0; i < 5; i++ {
		kod += fmt.Sprintf("%v", rand.Intn(10))
	}
	return kod
}

func MobilNotifikace(s string) {
	http.Post("https://ntfy.sh/novy_uzivatel115115jakopavouk", "text/plain", strings.NewReader(s))
}

func DelkaTextuArray(a []string) int {
	var x int
	for _, v := range a {
		x += len(v)
	}
	return x
}

/* vrací email, jmeno, error */
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

func volbaJmena(celeJmeno string) (string, error) {
	celeJmeno = godiacritics.Normalize(celeJmeno)
	var jmeno []string = strings.Fields(celeJmeno) // rozdělim na jmeno a prijimeni

	vsechnyJmena, err := databaze.GetVsechnyJmenaUziv()
	if err != nil {
		return "", err
	}

	var mapa = make(map[string]bool) // použijeme hash mapu žejoooo O(1)
	for _, v := range vsechnyJmena {
		mapa[v] = true
	}

	if testJmena(jmeno[0]+" "+jmeno[1], mapa) { // zkusim cely jmeno
		return jmeno[0] + " " + jmeno[1], nil
	}
	if testJmena(jmeno[0]+jmeno[1], mapa) { // zkusim cely jmeno bez mezery
		return jmeno[0] + jmeno[1], nil
	}
	if testJmena(jmeno[0], mapa) { // zkusim jmeno
		return jmeno[0], nil
	}
	if testJmena(jmeno[1], mapa) { // zkusim prijimeni
		return jmeno[1], nil
	}

	for i := 1; i < 999_999; i++ { // zkusim PavoukXXXXXX
		var j string = fmt.Sprintf("Pavouk%v", i)
		if !mapa[j] {
			return j, nil
		}
	}
	return "", errors.New("konec sveta")
}

func testJmena(s string, mapa map[string]bool) bool {
	return RegexJmeno.MatchString(s) && !mapa[s]
}

func BezDiakritiky(text []string, mezera bool) []string {
	for i := 0; i < len(text); i++ {
		slovo, _, err := transform.String(t, text[i])
		if err != nil {
			return text
		}
		if mezera {
			text[i] = slovo + " "
		} else {
			text[i] = slovo
		}
	}
	text[len(text)-1] = text[len(text)-1][:len(text[len(text)-1])-1] // smazat mezeru na konci
	return text
}
