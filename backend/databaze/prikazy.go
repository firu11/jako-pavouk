// V tomto souboru jsou příkazy do databáze které buď SELECTujou data z databáze, nebo data v databázi upravují.
// Největší dokumentace je jméno funkcí samotných.
package databaze

import (
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rickb777/date"
	godiacritics "gopkg.in/Regis24GmbH/go-diacritics.v2"
)

var RegexJmeno *regexp.Regexp
var MaxCisloZaJmeno int // 10_000
var poslednich int = 15

type (
	Lekce struct {
		ID      uint   `json:"id" db:"id"`
		Pismena string `json:"pismena" db:"pismena"`
		// Skupina uint        nepouzivame ale je tam
		// Klavesnice string
	}

	Cviceni struct {
		ID  uint   `json:"id" db:"id"`
		Typ string `json:"typ" db:"typ"`
	}

	Uzivatel struct {
		ID          uint      `json:"id" db:"id"`
		Email       string    `json:"email" db:"email"`
		Jmeno       string    `json:"jmeno" db:"jmeno"`
		Heslo       string    `json:"heslo" db:"heslo"`
		Klavesnice  string    `json:"klavesnice" db:"klavesnice"`
		Datum       date.Date `json:"datum" db:"datum"`
		Role        int       `json:"role" db:"role"`
		SkolniJmeno string    `json:"skolni_jmeno" db:"skolni_jmeno"`
	}

	NeoUziv struct {
		Email  string `json:"email" db:"email"`
		Jmeno  string `json:"jmeno" db:"jmeno"`
		Heslo  string `json:"heslo" db:"heslo"`
		Kod    string `json:"kod" db:"kod"`
		Cas    int64  `json:"cas" db:"cas"`
		Pokusy int    `json:"pokusy" db:"pokusy"`
	}

	ZmenaHeslaUziv struct {
		Email string `json:"email" db:"email"`
		Kod   string `json:"kod" db:"kod"`
		Cas   int64  `json:"cas" db:"cas"`
	}

	Slovnik struct {
		ID    uint   `json:"id" db:"id"`
		Slovo string `json:"slovo" db:"slovo"`
	}

	Dokoncene struct {
		ID            uint           `json:"id" db:"id"`
		UzivID        uint           `json:"uziv_id" db:"uziv_id"`
		CviceniID     uint           `json:"cviceni_id" db:"cviceni_id"`
		Neopravene    uint           `json:"neopravene" db:"neopravene"`
		Cas           float32        `json:"cas" db:"cas"`
		DelkaTextu    int            `json:"delka_textu" db:"delka_textu"`
		Datum         date.Date      `json:"datum" db:"datum"`
		ChybyPismenka map[string]int `json:"chyby_pismenka" db:"chyby_pismenka"`
	}

	Trida struct {
		ID       uint   `json:"id" db:"id"`
		Jmeno    string `json:"jmeno" db:"jmeno"`
		UcitelID uint   `json:"ucitel_id" db:"ucitel_id"`
		Kod      string `json:"kod" db:"kod"`
		Zamknuta bool   `json:"zamknuta" db:"zamknuta"`
		Smazana  bool   `json:"smazana" db:"smazana"`
	}
)

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
				_, err := GetUzivByJmeno(jmenoNaTest)
				if err != nil {
					return jmenoNaTest, nil
				}
			}
		}
		if len(jmeno) == 2 {
			jmenoNaTest = fmt.Sprintf("%s%d", jmeno[1], cislo)
			if RegexJmeno.MatchString(jmenoNaTest) {
				_, err := GetUzivByJmeno(jmenoNaTest)
				if err != nil { // ještě neexistuje
					return jmenoNaTest, nil
				}
			}
		}
		jmenoNaTest = fmt.Sprintf("Pavouk%d", cislo)
		if RegexJmeno.MatchString(jmenoNaTest) {
			_, err := GetUzivByJmeno(jmenoNaTest)
			if err != nil { // ještě neexistuje
				return jmenoNaTest, nil
			}
		}
	}

	return "", errors.New("konec sveta nenašel jsem jméno")
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

func GetProcvicovani(id int, cislo string) (string, []string, error) {
	var text string
	var nazev string

	err := DB.QueryRowx(`SELECT jmeno, text`+cislo+` FROM texty WHERE id = $1 ORDER by id;`, id).Scan(&nazev, &text)
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
	Id  int     `json:"id"`
	Cpm float32 `json:"cpm"`
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
	rows, err = DB.Queryx(`SELECT cviceni_id, MAX(((d.delka_textu - 10 * d.neopravene) / d.cas) * 60) AS cpm FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1 AND uziv_id = $2 GROUP BY d.cviceni_id;`, lekceID, uzivID)
	if err != nil {
		return cviceniIDs, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var cpm float32
		if err = rows.Scan(&id, &cpm); err != nil {
			return cviceniIDs, err
		}
		cviceniIDs = append(cviceniIDs, Cvic{id, cpm})
	}

	return cviceniIDs, nil
}

func GetDokonceneProcvic(uzivID uint) (map[int]float32, error) {
	var rychlosti map[int]float32 = make(map[int]float32)
	var rows *sqlx.Rows
	var err error

	rows, err = DB.Queryx(`SELECT procvic_id, MAX(((delka_textu - 10 * neopravene) / cas) * 60) AS cpm FROM dokoncene_procvic WHERE uziv_id = $1 GROUP BY procvic_id ORDER BY procvic_id DESC;`, uzivID)
	if err != nil {
		return rychlosti, err
	}

	defer rows.Close()

	for rows.Next() {
		var id sql.NullInt16
		var cpm float32
		if err = rows.Scan(&id, &cpm); err != nil {
			return rychlosti, err
		}
		rychlosti[int(id.Int16)-1] = cpm
	}

	return rychlosti, nil
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

func GetVsechnyJmenaUziv() ([]string, error) {
	var uzivatele []string
	rows, err := DB.Queryx(`SELECT jmeno FROM uzivatel;`)
	if err != nil {
		return uzivatele, err
	}
	defer rows.Close()

	for rows.Next() {
		var uziv Uzivatel
		err := rows.StructScan(&uziv)
		if err != nil {
			return uzivatele, err
		}

		uzivatele = append(uzivatele, uziv.Jmeno)
	}

	return uzivatele, nil
}

func SmazatUzivatele(id uint) error {
	_, err := DB.Exec(`DELETE FROM uzivatel WHERE id = $1 AND role != 2;`, id)
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

/*                          presnost,  cpm, daystreak, chybyPismenka */
func GetUdaje(uzivID uint) (float32, []float64, int, map[string]int, error) {
	var presnost float32
	var delkaVsechTextu int = 0
	var cpm []float64
	var daystreak int = 0
	var chybyPismenka map[string]int = make(map[string]int)

	rows, err := DB.Queryx(`SELECT neopravene, delka_textu, cas, chyby_pismenka, datum FROM dokoncene WHERE uziv_id = $1 UNION SELECT neopravene, delka_textu, cas, chyby_pismenka, datum FROM dokoncene_procvic WHERE uziv_id = $1 ORDER BY datum DESC LIMIT $2;`, uzivID, poslednich)
	if err != nil {
		return presnost, cpm, daystreak, chybyPismenka, err
	}
	defer rows.Close()

	for rows.Next() {
		var neopravene, delka int
		var cas float64
		var chybyPismenkaRowByte []byte
		var datumNezajima date.Date
		err := rows.Scan(&neopravene, &delka, &cas, &chybyPismenkaRowByte, &datumNezajima)
		if err != nil {
			return presnost, cpm, daystreak, chybyPismenka, err
		}

		var chybyPismenkaRow map[string]int
		err = json.Unmarshal(chybyPismenkaRowByte, &chybyPismenkaRow)
		if err == nil {
			for key, value := range chybyPismenkaRow {
				chybyPismenka[key] += value //když to ještě neexistuje, default value je 0
			}
		}
		cpm = append(cpm, utils.CPM(delka, cas, neopravene))
		delkaVsechTextu += delka
	}

	if delkaVsechTextu != 0 {
		var soucetChyb int = 0
		for _, hodnota := range chybyPismenka {
			soucetChyb += hodnota
		}
		presnost = float32(delkaVsechTextu-soucetChyb) / float32(delkaVsechTextu) * 100
		if presnost < 0 {
			presnost = 0 // kvuli adamovi kterej big troulin a měl -10%
		}
	}

	// daystreak
	rows, err = DB.Queryx(`SELECT datum, cas FROM dokoncene WHERE uziv_id = $1 UNION SELECT datum, cas FROM dokoncene_procvic WHERE uziv_id = $1 ORDER BY datum DESC;`, uzivID)
	if err != nil {
		return presnost, cpm, daystreak, chybyPismenka, err
	}
	defer rows.Close()

	var posledniHledanyDen date.Date = date.Today().AddDate(0, 0, -1) // vcera
	for rows.Next() {
		var c float32
		var d date.Date
		if err := rows.Scan(&d, &c); err != nil {
			return presnost, cpm, daystreak, chybyPismenka, err
		}

		if d == date.Today() && daystreak == 0 {
			daystreak++
		} else if posledniHledanyDen == d {
			daystreak++
			posledniHledanyDen = posledniHledanyDen.AddDate(0, 0, -1) // dalsi den
		} else {
			break
		}
	}
	if delkaVsechTextu == 0 {
		delkaVsechTextu = 1
	}
	return presnost, cpm, daystreak, chybyPismenka, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var x float32
	err := DB.QueryRowx(`WITH vsechny_cviceni AS (SELECT lekce_id, c.id as cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = (SELECT klavesnice FROM uzivatel WHERE id = $1) OR l.klavesnice = 'oboje'), moje_dokonceny AS (SELECT DISTINCT d.cviceni_id FROM dokoncene d JOIN vsechny_cviceni vc ON d.cviceni_id = vc.cviceni_id WHERE d.uziv_id = $1) SELECT (SELECT COUNT(*)::float FROM moje_dokonceny) / (SELECT COUNT(*) FROM vsechny_cviceni) as x;`, uzivID).Scan(&x)
	if err != nil {
		return 0, err
	}

	return x * 100, nil
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	email = strings.ToLower(email)

	var uzivID uint
	err := DB.QueryRowx(`INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3) RETURNING id;`, email, jmeno, hesloHash).Scan(&uzivID)
	if err != nil {
		return 0, err
	}
	return uzivID, nil
}

func PridatDokonceneCvic(cvicID, uzivID uint, neopravene int, cas float32, delkaTextu int, chybyPismenka map[string]int) error {
	chybyPismenkaJSON, err := json.Marshal(chybyPismenka)
	if err != nil {
		return errors.New("konverze mapy chyb na json se nepovedla")
	}
	_, err = DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, neopravene, cas, delka_textu, chyby_pismenka) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT duplicitni DO NOTHING`, uzivID, cvicID, neopravene, cas, delkaTextu, chybyPismenkaJSON)
	return err
}

func PridatDokonceneProcvic(procvicID, uzivID uint, neopravene int, cas float32, delkaTextu int, chybyPismenka map[string]int) error {
	chybyPismenkaJSON, err := json.Marshal(chybyPismenka)

	// pokud je procvic 0 neboli je to test psaní, vložim NULL
	var procvicCislo = sql.NullString{String: fmt.Sprintf("%d", procvicID), Valid: true}
	if procvicID == 0 {
		procvicCislo = sql.NullString{}
	}
	var id = sql.NullInt32{Int32: int32(uzivID), Valid: true}
	if uzivID == 0 {
		id = sql.NullInt32{}
	}
	if err != nil {
		return errors.New("konverze mapy chyb na json se nepovedla")
	}
	_, err = DB.Exec(`INSERT INTO dokoncene_procvic (uziv_id, procvic_id, neopravene, cas, delka_textu, chyby_pismenka) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT duplicitni2 DO NOTHING`, id, procvicCislo, neopravene, cas, delkaTextu, chybyPismenkaJSON)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}

func GetVsechnySlova(pocet int) ([]string, error) {
	var vysledek []string
	var err error

	rows, err := DB.Queryx(`WITH prvni AS (SELECT slovo FROM slovnik WHERE nahodnost < $1 ORDER BY nahodnost DESC LIMIT $4), druhy AS (SELECT slovo FROM slovnik WHERE nahodnost < $2 ORDER BY nahodnost DESC LIMIT $4), treti AS (SELECT slovo FROM slovnik WHERE nahodnost < $3 ORDER BY nahodnost DESC LIMIT $4) SELECT slovo FROM prvni UNION ALL SELECT slovo FROM druhy UNION ALL SELECT slovo FROM treti;`, mathRand.Float64()+0.0035, mathRand.Float64()+0.0035, mathRand.Float64()+0.0035, (pocet+3)/3) // +0.0035 proto aby tam byl dostatek slov když random hodí 0.0 (když jo, máme 100 slov)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	var slovo string
	for rows.Next() {
		if len(vysledek) == pocet {
			break
		}
		slovo = ""
		err := rows.Scan(&slovo)
		if err != nil {
			return vysledek, err
		}
		vysledek = append(vysledek, slovo)
	}
	return vysledek, nil
}

func GetVsechnyVety(pocet int) ([]string, error) {
	var vysledek []string
	var err error

	rows, err := DB.Queryx(`SELECT veta FROM vety ORDER BY RANDOM() LIMIT $1;`, pocet) // tady si asi random muzu dovolit
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	var veta string
	for rows.Next() {
		veta = ""
		err := rows.Scan(&veta)
		if err != nil {
			return vysledek, err
		}
		vysledek = append(vysledek, veta)
	}
	return vysledek, nil
}

func GetSlovaProLekci(uzivID uint, pismena string, pocet int) ([]string, error) {
	var vysledek []string
	var rows *sqlx.Rows

	if pismena == "velká písmena (shift)" || pismena == "čísla" || pismena == "interpunkce" {
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

func GetProgramatorSlova() ([]string, error) {
	var slova []string

	rows, err := DB.Queryx(`SELECT slovo FROM slovnik_programator ORDER BY RANDOM();`) // tady si taky muzu random dovolit (jen super malo rows)
	if err != nil {
		return slova, err
	}
	defer rows.Close()

	for rows.Next() {
		var slovo string
		err := rows.Scan(&slovo)
		if err != nil {
			return slova, err
		}

		slova = append(slova, slovo)
	}

	return slova, nil
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

func DalSpatnyKod(email string) {
	var pokusy int
	DB.QueryRowx(`UPDATE overeni SET pokusy = pokusy - 1 WHERE email = $1 RETURNING pokusy;`, email).Scan(&pokusy)
	if pokusy <= 0 {
		OdebratOvereni(email)
	}
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

func NovaNavsteva() error {
	_, err := DB.Exec(`INSERT INTO navstevnost (den, pocet) VALUES (CURRENT_DATE, 1) ON CONFLICT (den) DO UPDATE SET pocet = navstevnost.pocet + 1;`)
	return err
}

func CreateTrida(jmeno string, ucitelID uint, kod string) error {
	_, err := DB.Exec(`INSERT INTO trida (jmeno, ucitel_id, kod) VALUES ($1, $2, $3);`, jmeno, ucitelID, kod)
	return err
}

func GetTrida(id uint) (Trida, error) {
	var trida Trida
	err := DB.QueryRowx(`SELECT * FROM trida WHERE id = $1;`, id).StructScan(&trida)
	return trida, err
}

type TridaInfo struct {
	ID            string `json:"id" db:"id"`
	Jmeno         string `json:"jmeno" db:"jmeno"`
	Kod           string `json:"kod" db:"kod"`
	Zamknuta      bool   `json:"zamknuta" db:"zamknuta"`
	PocetStudentu int    `json:"pocet_studentu" db:"pocet_studentu"`
}

func GetTridy(ucitelID uint) ([]TridaInfo, error) {
	var tridy []TridaInfo = []TridaInfo{}

	rows, err := DB.Queryx(`SELECT id, jmeno, kod, zamknuta, (SELECT COUNT(*) FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id WHERE s.trida_id = trida.id) as pocet_studentu FROM trida WHERE ucitel_id = $1 AND smazana = FALSE;`, ucitelID)
	if err != nil {
		return tridy, err
	}

	defer rows.Close()

	for rows.Next() {
		var jendaTrida TridaInfo
		err := rows.StructScan(&jendaTrida)
		if err != nil {
			return tridy, err
		}

		tridy = append(tridy, jendaTrida)
	}
	return tridy, nil
}

type Student struct {
	ID    uint    `json:"id" db:"id"`
	Jmeno string  `json:"jmeno" db:"skolni_jmeno"`
	Email string  `json:"email" db:"email"`
	CPM   float64 `json:"cpm" db:"cpm"`
}

func GetStudentyZeTridy(tridaID uint) ([]Student, error) {
	var zaci []Student = []Student{}

	// super ultra šílený query by LLM
	// předtím jsem pro každého studenta posílal query samostatně (30 žáků | 900ms -> 80ms)
	rows, err := DB.Queryx(`WITH cpm_data AS (SELECT datum, (((delka_textu - 10 * neopravene) / cas) * 60) AS cpm, uziv_id FROM dokoncene UNION ALL SELECT datum, (((delka_textu - 10 * neopravene) / cas) * 60) AS cpm, uziv_id FROM dokoncene_procvic), cpm_filtered AS (SELECT uziv_id, datum, CASE WHEN cpm < 0 THEN 0 ELSE cpm END AS cpm, ROW_NUMBER() OVER (PARTITION BY uziv_id ORDER BY datum DESC) AS rn FROM cpm_data), latest_15_cpm AS (SELECT uziv_id, datum, cpm FROM cpm_filtered WHERE rn <= 15), cpm AS (SELECT uziv_id, AVG(cpm) AS cpm FROM latest_15_cpm GROUP BY uziv_id) SELECT u.id, u.skolni_jmeno, u.email, COALESCE(mc.cpm, 0) as cpm FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id INNER JOIN trida t ON t.id = s.trida_id FULL OUTER JOIN latest_15_cpm l15 ON l15.uziv_id = u.id FULL OUTER JOIN cpm mc ON mc.uziv_id = u.id WHERE s.trida_id = $1 AND t.smazana = FALSE GROUP BY u.id, u.skolni_jmeno, u.email, mc.cpm;`, tridaID)
	if err != nil {
		return zaci, err
	}
	defer rows.Close()

	for rows.Next() {
		var zak Student
		err := rows.StructScan(&zak)
		if err != nil {
			return zaci, err
		}

		zaci = append(zaci, zak)
	}
	return zaci, nil
}

func ZamknoutTridu(tridaID uint) error {
	_, err := DB.Exec(`UPDATE trida SET zamknuta = NOT zamknuta WHERE id = $1;`, tridaID)
	return err
}

func PrejmenovatTridu(tridaID uint, noveJmeno string) error {
	_, err := DB.Exec(`UPDATE trida SET jmeno = $1 WHERE id = $2;`, noveJmeno, tridaID)
	return err
}

func SmazatTridu(id uint) error {
	_, err := DB.Exec(`UPDATE trida SET smazana = TRUE WHERE id = $1;`, id)
	return err
}

func PrejmenovatStudenta(id uint, skolniJmeno string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET skolni_jmeno = $1 WHERE id = $2;`, skolniJmeno, id)
	return err
}

func GetTridaByKod(kod string) (Trida, error) {
	var trida Trida
	err := DB.QueryRowx(`SELECT * FROM trida WHERE kod = $1 AND NOT smazana;`, kod).StructScan(&trida)
	return trida, err
}

func GetTridaByUziv(id uint) (Trida, error) {
	var trida Trida
	err := DB.QueryRowx(`SELECT t.* FROM trida t INNER JOIN student_a_trida s ON t.id = s.trida_id WHERE student_id = $1 AND NOT smazana;`, id).StructScan(&trida)
	return trida, err
}

func ZapsatStudenta(kod string, studentID uint, jmeno string) error {
	kod = strings.ToUpper(kod)

	var id int
	var smazana bool
	err := DB.QueryRowx(`SELECT s.trida_id, t.smazana FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id INNER JOIN trida t ON t.id = s.trida_id WHERE u.id = $1 AND NOT t.smazana;`, studentID).Scan(&id, &smazana)
	if err == nil {
		return errors.New("uz je ve tride")
	}
	if err != sql.ErrNoRows {
		return err
	}
	if smazana {
		return errors.New("trida je smazana")
	}

	_, err = DB.Exec(`INSERT INTO student_a_trida (student_id, trida_id) VALUES ($1, (SELECT id FROM trida WHERE kod = $2 AND NOT zamknuta)) ON CONFLICT DO NOTHING;`, studentID, kod)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`UPDATE uzivatel SET skolni_jmeno = $1 WHERE id = $2;`, jmeno, studentID)
	return err
}
