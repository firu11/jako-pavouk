package databaze

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rickb777/date"
)

type Lekce struct {
	ID      uint   `json:"id"`
	Pismena string `json:"pismena"`
	// Skupina uint        nepouzivame ale je tam
	// Klavesnice string
}

type Cviceni struct {
	ID  uint   `json:"id"`
	Typ string `json:"typ"`
}

type Uzivatel struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Jmeno      string `json:"jmeno"`
	Heslo      string `json:"heslo"`
	Klavesnice string `json:"klavesnice"`
}

type NeoUziv struct {
	Email string `json:"email"`
	Jmeno string `json:"jmeno"`
	Heslo string `json:"heslo"`
	Kod   string `json:"kod"`
	Cas   int64  `json:"cas"`
}

type ZmenaHeslaUziv struct {
	Email string `json:"email"`
	Kod   string `json:"kod"`
	Cas   int64  `json:"cas"`
}

type Slovnik struct {
	ID    uint   `json:"id"`
	Slovo string `json:"slovo"`
}

type Dokoncene struct {
	ID         uint      `json:"id"`
	UzivID     uint      `json:"uziv_id"`
	CviceniID  uint      `json:"cviceni_id"`
	CPM        float32   `json:"cpm"`
	Preklepy   uint      `json:"preklepy"`
	Cas        int       `json:"cas"`
	DelkaTextu int       `json:"delka_textu"`
	Datum      date.Date `json:"datum"`
}

func GetLekce(uzivID uint) ([][]Lekce, error) {
	var lekce [][]Lekce = [][]Lekce{}

	rows, err := DB.Queryx(`SELECT id, pismena, skupina FROM lekce WHERE klavesnice = 'oboje' OR klavesnice = COALESCE((SELECT klavesnice FROM uzivatel WHERE id = $1), 'qwertz') ORDER BY id ASC;`, uzivID)
	if err != nil {
		return lekce, err
	}

	defer rows.Close()

	var skupina []Lekce
	var cisloSkupiny uint = 1
	var jednaLekce Lekce

	for rows.Next() {
		jednaLekce = Lekce{}
		var skup uint

		err := rows.Scan(&jednaLekce.ID, &jednaLekce.Pismena, &skup)
		if err != nil {
			return lekce, err
		}
		if cisloSkupiny == skup {
			skupina = append(skupina, jednaLekce)
		} else if cisloSkupiny < skup {
			lekce = append(lekce, skupina)
			skupina = []Lekce{jednaLekce}
			cisloSkupiny += 1
		} else {
			lekce[skup] = append(lekce[skup], jednaLekce)
		}
	}
	lekce = append(lekce, skupina)
	return lekce, nil
}

func GetDokonceneLekce(uzivID uint) ([]int32, error) {
	var vysledek []int32 = []int32{}
	// zjistim kolik ma kazda lekce cviceni
	rows, err := DB.Queryx(`WITH vsechny_cviceni AS (SELECT lekce_id, c.id as cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = (SELECT klavesnice FROM uzivatel WHERE id = $1) OR l.klavesnice = 'oboje'), moje_dokonceny AS (SELECT 1 as dokonceno, d.cviceni_id FROM dokoncene d WHERE d.uziv_id = $1) SELECT lekce_id FROM vsechny_cviceni vc LEFT JOIN moje_dokonceny d ON vc.cviceni_id = d.cviceni_id GROUP BY lekce_id HAVING (COUNT(*)) = (COUNT (*) FILTER (WHERE d.dokonceno IS NOT NULL));`, uzivID)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return vysledek, err
		}
		vysledek = append(vysledek, int32(id))
	}
	return vysledek, nil
}

func GetTexty() ([]string, error) {
	var texty []string

	rows, err := DB.Queryx(`SELECT jmeno FROM texty ORDER BY id;`)
	if err != nil {
		return texty, err
	}

	for rows.Next() {
		var t string
		if err = rows.Scan(&t); err != nil {
			return texty, err
		}
		texty = append(texty, t)
	}
	return texty, nil
}

func GetProcvicovani(id int) (string, []string, error) {
	var text string
	var nazev string
	var randomText string = strconv.Itoa(rand.Intn(5-1) + 1)

	err := DB.QueryRowx(`SELECT jmeno, text`+randomText+` FROM texty WHERE id = $1;`, id).Scan(&nazev, &text)
	if err != nil {
		return "", []string{}, err
	}

	var textArr []string = strings.Fields(text)
	for i := 0; i < len(textArr)-1; i++ {
		textArr[i] += " "
	}

	return nazev, textArr, nil
}

type Cvic struct {
	Id       int
	Cpm      float32
	Presnost float32
}

func GetDokonceneCvicVLekci(uzivID uint, lekceID uint, pismena string) ([]Cvic, error) {
	var cviceniIDs []Cvic = []Cvic{}
	var rows *sqlx.Rows
	var err error

	if pismena != "" {
		lekceID, err = GetLekceIDbyPismena(pismena)
		if err != nil {
			return cviceniIDs, err
		}
	}
	rows, err = DB.Queryx(`SELECT d.cviceni_id, d.cpm, d.delka_textu, d.preklepy FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1 AND uziv_id = $2;`, lekceID, uzivID)
	if err != nil {
		return cviceniIDs, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, delkaTextu int
		var preklepy, cpm float32
		if err = rows.Scan(&id, &cpm, &delkaTextu, &preklepy); err != nil {
			return cviceniIDs, err
		}
		cviceniIDs = append(cviceniIDs, Cvic{id, cpm, (float32(delkaTextu) - preklepy) / float32(delkaTextu) * 100})
	}

	return cviceniIDs, nil
}

func GetLekceIDbyPismena(pismena string) (uint, error) {
	var id uint
	err := DB.QueryRowx(`SELECT id FROM lekce WHERE pismena = $1;`, pismena).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCviceniVLekciByID(lekceID uint) ([]Cviceni, error) {
	var cviceni []Cviceni

	rows, err := DB.Queryx(`SELECT id, typ FROM cviceni WHERE lekce_id = $1;`, lekceID)
	if err != nil {
		return cviceni, err
	}
	defer rows.Close()

	for rows.Next() {
		var jednoCviceni Cviceni
		err := rows.StructScan(&jednoCviceni)
		if err != nil {
			return cviceni, err
		}

		cviceni = append(cviceni, jednoCviceni)
	}

	return cviceni, nil
}

func GetCviceniVLekciByPismena(uzivID uint, pismena string) ([]Cviceni, error) {
	var cviceni []Cviceni

	rows, err := DB.Queryx(`SELECT id, typ FROM cviceni WHERE lekce_id = (SELECT id FROM lekce where pismena = $1 LIMIT 1);`, pismena)
	if err != nil {
		return cviceni, err
	}

	defer rows.Close()

	for rows.Next() {
		var jednoCviceni Cviceni
		err := rows.StructScan(&jednoCviceni)
		if err != nil {
			return cviceni, err
		}

		cviceni = append(cviceni, jednoCviceni)
	}
	if len(cviceni) == 0 {
		return cviceni, errors.New("nejsou zadny takovy cviceni")
	}
	return cviceni, nil
}

func GetUzivByID(uzivID uint) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRowx(`SELECT * FROM uzivatel WHERE id = $1;`, uzivID).StructScan(&uziv)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRowx(`SELECT * FROM uzivatel WHERE email = $1;`, email).StructScan(&uziv)
	return uziv, err
}

func GetUzivByJmeno(jmeno string) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRowx(`SELECT * FROM uzivatel WHERE jmeno = $1;`, jmeno).StructScan(&uziv)
	return uziv, err
}

func SmazatUzivatele(id uint) error {
	_, err := DB.Exec(`DELETE FROM uzivatel WHERE id = $1;`, id)
	return err
}

func ZmenitKlavesnici(id uint, novaKlavesnice string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET klavesnice = $1 WHERE id = $2;`, novaKlavesnice, id)
	return err // buď nil nebo error
}

func PrejmenovatUziv(id uint, noveJmeno string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET jmeno = $1 WHERE id = $2;`, noveJmeno, id)
	return err // buď nil nebo error
}

/* preklepy, cpm, daystreak, cas, delka */
func GetUdaje(uzivID uint) (int, []float32, int, float32, int, error) {
	var preklepy int
	var delkaVsechTextu int = 0
	var cpm []float32
	var daystreak int = 0
	var celkovyCas float32 = 0

	var poslednich int = 10
	rows, err := DB.Queryx(`SELECT preklepy, cpm, delka_textu FROM dokoncene WHERE uziv_id = $1 ORDER BY den DESC LIMIT $2;`, uzivID, poslednich)
	if err != nil {
		return preklepy, cpm, daystreak, celkovyCas, delkaVsechTextu, err
	}
	defer rows.Close()

	for rows.Next() {
		var preklep int
		var cpmko float32
		var delka int
		err := rows.Scan(&preklep, &cpmko, &delka)
		if err != nil {
			return preklepy, cpm, daystreak, celkovyCas, delkaVsechTextu, err
		}

		preklepy += preklep
		cpm = append(cpm, cpmko)
		delkaVsechTextu += delka
	}

	rows, err = DB.Queryx(`SELECT den, cas FROM dokoncene WHERE uziv_id = $1 ORDER BY den DESC;`, uzivID)
	if err != nil {
		return preklepy, cpm, daystreak, celkovyCas, delkaVsechTextu, err
	}
	defer rows.Close()

	var dny []date.Date
	for rows.Next() {
		var c float32
		var d date.Date
		if err := rows.Scan(&d, &c); err != nil {
			return preklepy, cpm, daystreak, celkovyCas, delkaVsechTextu, err
		}
		celkovyCas += c
		dny = append(dny, d)
	}

	var hledanyDen date.Date = date.Today().Add(-1)
	for _, d := range dny {
		if hledanyDen.Equal(d) {
			hledanyDen = hledanyDen.Add(-1)
			daystreak++
		} else if d.Equal(date.Today()) {
			daystreak = 1
		} else if hledanyDen.Sub(d) != -1 { // pokud dalsi je uz víc davno nez vcera
			break
		}
	}
	if delkaVsechTextu == 0 {
		delkaVsechTextu = 1
	}
	return preklepy, cpm, daystreak, celkovyCas, delkaVsechTextu, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var x float32
	err := DB.QueryRowx(`WITH vsechny_cviceni AS (SELECT lekce_id, c.id as cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = (SELECT klavesnice FROM uzivatel WHERE id = $1) OR l.klavesnice = 'oboje'), moje_dokonceny AS (SELECT 1 as dokonceno, d.cviceni_id FROM dokoncene d JOIN vsechny_cviceni vc ON d.cviceni_id = vc.cviceni_id WHERE d.uziv_id = $1) SELECT (SELECT COUNT(*)::float FROM moje_dokonceny) / (SELECT COUNT(*) FROM vsechny_cviceni) as x;`, uzivID).Scan(&x)
	if err != nil {
		return 0, err
	}

	return x * 100, nil
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	var uzivID uint
	err := DB.QueryRowx(`INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3) RETURNING id;`, email, jmeno, hesloHash).Scan(&uzivID)
	if err != nil {
		return 0, err
	}
	return uzivID, nil
}

func PridatDokonceneCvic(cvicID, uzivID uint, cpm float32, preklepy int, cas float32, delkaTextu int) error {
	cpm = float32(math.Round(float64(cpm)*100) / 100) // zaokrouhlit na 2 desetiny cisla
	cas = float32(math.Round(float64(cas)*100) / 100)
	_, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy, cas, delka_textu) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT unikatni DO UPDATE SET cpm = EXCLUDED.cpm, preklepy = EXCLUDED.preklepy, cas = EXCLUDED.cas, delka_textu = EXCLUDED.delka_textu, den = CURRENT_TIMESTAMP;`, uzivID, cvicID, cpm, preklepy, cas, delkaTextu)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}

func GetSlovaProLekci(uzivID uint, pismena string, pocet int) ([]string, error) {
	var vysledek []string
	var rows *sqlx.Rows

	if pismena == "Velká písmena (Shift)" {
		var err error
		rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwertz_id <= (SELECT id from lekce WHERE pismena = $1) ORDER BY RANDOM() LIMIT $2;`, pismena, pocet)
		if err != nil {
			return vysledek, err
		}
	} else {
		var k string
		err := DB.QueryRowx(`SELECT klavesnice FROM uzivatel WHERE id = $1;`, uzivID).Scan(&k)
		if err != nil {
			return vysledek, err
		}

		if k == "qwertz" {
			rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwertz_id = (SELECT id from lekce WHERE pismena = $1) ORDER BY RANDOM() LIMIT $2;`, pismena, pocet)
		} else {
			rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwerty_id = (SELECT id from lekce WHERE pismena = $1) ORDER BY RANDOM() LIMIT $2;`, pismena, pocet)
		}
		if err != nil {
			return vysledek, err
		}
	}

	defer rows.Close()

	var slovo string
	for rows.Next() {
		slovo = ""
		err := rows.Scan(&slovo)
		if err != nil {
			return vysledek, err
		}
		vysledek = append(vysledek, slovo)
	}
	return vysledek, nil
}

func GetNaucenaPismena(uzivID uint, pismena string) (string, error) {
	var vysledek strings.Builder
	rows, err := DB.Queryx(`SELECT pismena FROM lekce WHERE id <= (SELECT id from lekce WHERE pismena = $1) AND (klavesnice = COALESCE((SELECT klavesnice FROM uzivatel WHERE id = $2), 'qwertz') OR klavesnice = 'oboje');`, pismena, uzivID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var pismenaJedny string
	for rows.Next() {
		pismenaJedny = ""
		err := rows.Scan(&pismenaJedny)
		if err != nil {
			return "", err
		}

		vysledek.WriteString(pismenaJedny)
	}

	return vysledek.String(), nil
}

func CreateNeoverenyUziv(email, hesloHASH, jmeno, kod string, cas int64) error {
	_, err := DB.Exec(`INSERT INTO overeni (email, jmeno, heslo, kod, cas) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO UPDATE SET jmeno = EXCLUDED.jmeno, heslo = EXCLUDED.heslo, kod = EXCLUDED.kod, cas = EXCLUDED.cas;`, email, jmeno, hesloHASH, kod, cas)
	return err
}

func GetNeoverenyUziv(email string) (NeoUziv, error) {
	var uziv NeoUziv
	err := DB.QueryRowx(`SELECT * FROM overeni WHERE email = $1;`, email).StructScan(&uziv)
	return uziv, err
}

func OdebratOvereni(email string) error {
	_, err := DB.Exec(`DELETE FROM overeni WHERE email = $1`, email)
	return err
}

func SmazatPoLimitu() error {
	var now int64 = time.Now().Unix()

	_, err := DB.Exec(`DELETE FROM overeni WHERE cas < $1;`, now)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`DELETE FROM zmena_hesla WHERE cas < $1;`, now)
	return err
}

func CreateZapomenuteHeslo(email, kod string, cas int64) error {
	_, err := DB.Exec(`INSERT INTO zmena_hesla (email, kod, cas) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET kod = EXCLUDED.kod, cas = EXCLUDED.cas;`, email, kod, cas)
	return err
}

func OdebratZmenuHesla(email string) error {
	_, err := DB.Exec(`DELETE FROM zmena_hesla WHERE email = $1`, email)
	return err
}

func GetZmenuHesla(email string) (ZmenaHeslaUziv, error) {
	var uziv ZmenaHeslaUziv
	err := DB.QueryRowx(`SELECT * FROM zmena_hesla WHERE email = $1;`, email).StructScan(&uziv)
	return uziv, err
}

func ZmenitHeslo(email, hesloHASH string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET heslo = $1 WHERE email = $2`, hesloHASH, email)
	return err
}
