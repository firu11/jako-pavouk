package databaze

import (
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

	"github.com/blockloop/scan/v2"
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
		ID              uint      `json:"id" db:"id"`
		Email           string    `json:"email" db:"email"`
		Jmeno           string    `json:"jmeno" db:"jmeno"`
		Heslo           string    `json:"heslo" db:"heslo"`
		Klavesnice      string    `json:"klavesnice" db:"klavesnice"`
		Datum           date.Date `json:"datum" db:"datum"`
		Role            int       `json:"role" db:"role"`
		UcitelVeSkoleID int       `json:"" db:"skola_id"`
		SkolniJmeno     string    `json:"skolni_jmeno" db:"skolni_jmeno"`
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
		ID         uint   `json:"id" db:"id"`
		Jmeno      string `json:"jmeno" db:"jmeno"`
		UcitelID   uint   `json:"-" db:"ucitel_id"`
		Kod        string `json:"kod" db:"kod"`
		Zamknuta   bool   `json:"zamknuta" db:"zamknuta"`
		Smazana    bool   `json:"smazana" db:"smazana"`
		Klavesnice string `json:"klavesnice" db:"klavesnice"`
	}

	Skola struct {
		ID          uint      `json:"id" db:"id"`
		Jmeno       string    `json:"jmeno" db:"jmeno"`
		DenZalozeni date.Date `json:"den_zalozeni" db:"den_zalozeni"`
		Aktivni     bool      `json:"aktivni" db:"aktivni"`
		// kontaktni_email
		// kontaktni_telefon
	}

	Ucitel struct {
		ID      uint `json:"id" db:"id"`
		SkolaID uint `json:"skola_id" db:"skola_id"`

		UzivID uint   `json:"uziv_id" db:"uziv_id"`
		Email  string `json:"email" db:"email"`

		PocetTrid int `json:"pocet_trid" db:"pocet_trid"`
		PocetZaku int `json:"pocet_zaku" db:"pocet_zaku"`
	}

	Prace struct {
		ID                uint      `json:"id" db:"id"`
		TridaID           uint      `json:"-" db:"trida_id"`
		Text              string    `json:"text" db:"text"`
		Cas               int       `json:"cas" db:"cas"`
		Datum             time.Time `json:"datum" db:"datum"`
		PrumerneCPM       float32   `json:"prumerne_cpm" db:"prumerne_cpm"`
		PrumernaPresnost  float32   `json:"prumerna_presnost" db:"prumerna_presnost"`
		StudentuDokoncilo int       `json:"studentu_dokoncilo" db:"studentu_dokoncilo"`
	}

	Procvic struct {
		ID        uint    `json:"id" db:"id"`
		Jmeno     string  `json:"jmeno" db:"jmeno"`
		Kategorie string  `json:"kategorie" db:"kategorie"`
		CPM       float32 `json:"cpm" db:"-"`
		Obtiznost int     `json:"obtiznost" db:"obtiznost"`
	}
)

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

	rows, err := DB.Query(`SELECT id, pismena, skupina FROM lekce WHERE klavesnice = 'oboje' OR klavesnice = COALESCE((SELECT klavesnice FROM uzivatel WHERE id = $1), 'qwertz') ORDER BY id ASC;`, uzivID)
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
	rows, err := DB.Query(`WITH vsechny_cviceni AS (SELECT lekce_id, c.id as cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = (SELECT klavesnice FROM uzivatel WHERE id = $1) OR l.klavesnice = 'oboje'), moje_dokonceny AS (SELECT 1 as dokonceno, d.cviceni_id FROM dokoncene d WHERE d.uziv_id = $1) SELECT lekce_id FROM vsechny_cviceni vc LEFT JOIN moje_dokonceny d ON vc.cviceni_id = d.cviceni_id GROUP BY lekce_id HAVING (COUNT(*)) = (COUNT (*) FILTER (WHERE d.dokonceno IS NOT NULL));`, uzivID)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	err = scan.Rows(&vysledek, rows)

	return vysledek, err
}

func GetDalsiCviceni(uzivID uint) (string, int32, error) {
	var cislo int32
	var lekcePismena string
	err := DB.QueryRow(`WITH klavesnice_user AS ( SELECT klavesnice FROM uzivatel WHERE id = $1 ), vsechny_cviceni AS ( SELECT lekce_id, c.id AS cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = ( SELECT klavesnice FROM klavesnice_user ) OR l.klavesnice = 'oboje' ), moje_dokonceny AS ( SELECT d.cviceni_id FROM dokoncene d WHERE d.uziv_id = $1 ), hledany_cvic_id AS ( SELECT vc.lekce_id, vc.cviceni_id FROM vsechny_cviceni vc LEFT JOIN moje_dokonceny d ON vc.cviceni_id = d.cviceni_id WHERE d.cviceni_id IS NULL ORDER BY vc.lekce_id ASC, vc.cviceni_id ASC LIMIT 1 ), cviceni_v_lekci AS ( SELECT c.id AS cviceni_id, l.id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE ( l.klavesnice = ( SELECT klavesnice FROM klavesnice_user ) OR l.klavesnice = 'oboje' ) AND l.id = ( SELECT lekce_id FROM hledany_cvic_id ) ), finalni AS ( SELECT l.pismena, ROW_NUMBER() OVER ( ORDER BY c.cviceni_id ASC ) AS cislo, c.cviceni_id FROM cviceni_v_lekci c JOIN lekce l ON l.id = c.id ) SELECT pismena, cislo FROM finalni WHERE cviceni_id = ( SELECT cviceni_id FROM hledany_cvic_id );`, uzivID).Scan(&lekcePismena, &cislo)
	if err != nil {
		return lekcePismena, cislo, err
	}

	return lekcePismena, cislo, nil
}

func GetTexty() ([]Procvic, error) {
	var texty []Procvic

	rows, err := DB.Query(`SELECT id, jmeno, kategorie, obtiznost FROM druhy_textu ORDER BY jmeno;`)
	if err != nil {
		return texty, err
	}

	err = scan.Rows(&texty, rows)

	return texty, err
}

func GetProcvicovani(id, cislo int) (string, string, []string, int, error) {
	var text, nazev, podnazev string
	var realCislo int

	r := DB.QueryRow(`WITH maximum AS ( SELECT MAX(cislo) as m FROM texty WHERE typ = $1) SELECT d.jmeno as nazev, t.jmeno as podnazev, t.txt, t.cislo FROM texty t JOIN druhy_textu d ON t.typ = d.id, maximum WHERE d.id = $1 AND t.cislo = ((($2 - 1) % maximum.m) + 1) LIMIT 1;`, id, cislo)
	err := r.Scan(&nazev, &podnazev, &text, &realCislo)
	if err != nil {
		return "", "", []string{}, 0, err
	}

	var textArr []string = strings.Fields(text)
	for i := 0; i < len(textArr)-1; i++ {
		textArr[i] += " "
	}

	return nazev, podnazev, textArr, realCislo, nil
}

func GetRandomProcvic(jmeno string) (string, error) {
	var text string

	err := DB.QueryRow(`SELECT t.txt FROM texty t INNER JOIN druhy_textu d ON t.typ = d.id WHERE d.jmeno = $1 ORDER BY RANDOM() LIMIT 1;`, jmeno).Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

type Cvic struct {
	Id  int     `json:"id"`
	Cpm float32 `json:"cpm"`
}

func GetDokonceneCvicVLekci(uzivID uint, lekceID uint, pismena string) ([]Cvic, error) {
	var cviceniIDs []Cvic = []Cvic{}
	var rows *sql.Rows
	var err error

	if pismena != "" {
		lekceID, err = GetLekceIDbyPismena(pismena)
		if err != nil {
			return cviceniIDs, err
		}
	}
	rows, err = DB.Query(`SELECT cviceni_id, MAX(((d.delka_textu - 10 * d.neopravene)::float / d.cas) * 60) AS cpm FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1 AND uziv_id = $2 GROUP BY d.cviceni_id;`, lekceID, uzivID)
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

func GetRychlostiProcvic(uzivID uint) (map[int]float32, error) {
	var rychlosti map[int]float32 = make(map[int]float32)
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`WITH sus AS ( SELECT * FROM ( SELECT ROW_NUMBER() OVER ( PARTITION BY typ_textu ORDER BY datum DESC ) AS r, d.* FROM dokoncene_procvic d WHERE uziv_id = $1 ) AS idk WHERE idk.r <= $2 ) SELECT typ_textu, GREATEST( ( (SUM(delka_textu) - 10 * SUM(neopravene)) / SUM(cas)::NUMERIC ) * 60, 0 ) AS cpm FROM sus GROUP BY typ_textu;`, uzivID, poslednich)
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
	err := DB.QueryRow(`SELECT id FROM lekce WHERE pismena = $1;`, pismena).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCviceniVLekciByID(lekceID uint) ([]Cviceni, error) {
	var cviceni []Cviceni

	rows, err := DB.Query(`SELECT id, typ FROM cviceni WHERE lekce_id = $1 ORDER BY id;`, lekceID)
	if err != nil {
		return cviceni, err
	}

	err = scan.Rows(&cviceni, rows)
	return cviceni, err
}

func GetCviceniVLekciByPismena(uzivID uint, pismena string) ([]Cviceni, error) {
	var cviceni []Cviceni

	rows, err := DB.Query(`SELECT id, typ FROM cviceni WHERE lekce_id = (SELECT id FROM lekce where pismena = $1 LIMIT 1) ORDER BY id;`, pismena)
	if err != nil {
		return cviceni, err
	}
	err = scan.Rows(&cviceni, rows)
	if err != nil {
		return cviceni, err
	}
	if len(cviceni) == 0 {
		return cviceni, errors.New("nejsou zadny takovy cviceni")
	}
	return cviceni, nil
}

func GetUzivByID(uzivID uint) (Uzivatel, error) {
	var uziv Uzivatel
	row, err := DB.Query(`SELECT uz.*, COALESCE(u.skola_id, 0) AS skola_id FROM uzivatel uz FULL OUTER JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.id = $1 AND NOT uz.smazany;`, uzivID)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv Uzivatel
	row, err := DB.Query(`SELECT uz.*, COALESCE(u.skola_id, 0) AS skola_id FROM uzivatel uz FULL OUTER JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.email = $1 AND NOT uz.smazany;`, email)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func GetUzivByJmeno(jmeno string) (Uzivatel, error) {
	var uziv Uzivatel
	row, err := DB.Query(`SELECT uz.*, COALESCE(u.skola_id, 0) AS skola_id FROM uzivatel uz FULL OUTER JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.jmeno = $1 AND NOT uz.smazany;`, jmeno)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func GetVsechnyJmenaUziv() ([]string, error) {
	var uzivatele []string
	rows, err := DB.Query(`SELECT jmeno FROM uzivatel WHERE NOT smazany;`)
	if err != nil {
		return uzivatele, err
	}
	err = scan.Rows(&uzivatele, rows)
	return uzivatele, err
}

func SmazatUzivatele(id uint) error {
	_, err := DB.Exec(`UPDATE uzivatel SET smazany = TRUE WHERE id = $1 AND role != 2;`, id)
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

func GetDaystreak(uzivID uint) (int, error) {
	var daystreak int = 0

	rows, err := DB.Query(`SELECT datum FROM dokoncene WHERE uziv_id = $1 UNION SELECT datum FROM dokoncene_procvic WHERE uziv_id = $1 ORDER BY datum DESC;`, uzivID)
	if err != nil {
		return daystreak, err
	}
	defer rows.Close()

	var posledni date.Date = date.Today()
	for rows.Next() {
		var d date.Date
		if err := rows.Scan(&d); err != nil {
			return daystreak, err
		}

		if d == posledni {
			if daystreak == 0 {
				daystreak++
			}
		} else if posledni.AddDate(0, 0, -1) == d {
			daystreak++
			posledni = d
		} else {
			break
		}
	}

	return daystreak, nil
}

/*                          presnost,  cpm,  chybyPismenka,  cas, napsanychPismen */
func GetUdaje(uzivID uint) (float32, float32, map[string]int, [3]int, [3]int, error) {
	var presnost float32 = -1
	var rychlost float32 = -1
	var chybyPismenkaJsonb []byte
	var chybyPismenka map[string]int

	var cas [3]int
	var napsanychPismen [3]int

	err := DB.QueryRow(`WITH default_hodnoty AS ( SELECT 0 AS neopravene, 0 AS delka_textu, 0 AS cas, CURRENT_DATE AS datum, 0 AS opravene, '{}'::jsonb AS chyby_pismenka ), vsechny_za_dva_tydny AS ( SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene, chyby_pismenka FROM dokoncene WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene, chyby_pismenka FROM dokoncene_procvic WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT * FROM default_hodnoty ), soucty_pismenek AS ( SELECT key AS pismeno, SUM(value::NUMERIC) AS soucet FROM vsechny_za_dva_tydny, jsonb_each_text(chyby_pismenka) GROUP BY key UNION ALL SELECT '' AS pismeno, 0 AS soucet WHERE NOT EXISTS ( SELECT 1 FROM vsechny_za_dva_tydny, jsonb_each_text(chyby_pismenka) ) ), vsechny AS ( SELECT delka_textu, cas, datum::date FROM dokoncene WHERE uziv_id = $1 UNION SELECT delka_textu, cas, datum::date FROM dokoncene_procvic WHERE uziv_id = $1 ), soucty_dnes AS ( SELECT SUM(cas) AS cas_dnes, SUM(delka_textu) AS napsanych_pismen_dnes FROM vsechny WHERE datum = CURRENT_DATE ), soucty_dva_tydny AS ( SELECT SUM(cas) AS cas_dva_tydny, SUM(delka_textu) AS napsanych_pismen_dva_tydny FROM vsechny WHERE datum > CURRENT_DATE - MAKE_INTERVAL(days => 14) ), soucty_celkem AS ( SELECT SUM(cas) AS cas_celkem, SUM(delka_textu) AS napsanych_pismen_celkem FROM vsechny ) SELECT GREATEST( ( (SUM(delka_textu) - 10 * SUM(neopravene)) / GREATEST(SUM(cas)::NUMERIC, 1) ) * 60, 0 ) AS rychlost, COALESCE( ( ( SUM(delka_textu) - SUM(neopravene) - COALESCE(SUM(opravene), 0) ) / GREATEST(SUM(delka_textu)::NUMERIC, 1) ) * 100, -1 ) AS presnost, COALESCE(jsonb_object_agg(pismeno, soucet), '{}'::jsonb) AS chyby_pismenka, COALESCE(max(cas_dnes), 0) AS cas_dnes, COALESCE(max(cas_dva_tydny), 0) AS cas_dva_tydny, COALESCE(max(cas_celkem), 0) AS cas_celkem, COALESCE(max(napsanych_pismen_dnes), 0) AS napsanych_pismen_dnes, COALESCE(max(napsanych_pismen_dva_tydny), 0) AS napsanych_pismen_dva_tydny, COALESCE(max(napsanych_pismen_celkem), 0) AS napsanych_pismen_celkem FROM soucty_pismenek, vsechny_za_dva_tydny, soucty_dnes, soucty_dva_tydny, soucty_celkem;`, uzivID).Scan(&rychlost, &presnost, &chybyPismenkaJsonb, &cas[0], &cas[1], &cas[2], &napsanychPismen[0], &napsanychPismen[1], &napsanychPismen[2])
	if err != nil {
		return presnost, rychlost, chybyPismenka, cas, napsanychPismen, err
	}
	err = json.Unmarshal(chybyPismenkaJsonb, &chybyPismenka)
	if err != nil {
		return presnost, rychlost, chybyPismenka, cas, napsanychPismen, err
	}

	return presnost, rychlost, chybyPismenka, cas, napsanychPismen, nil
}

func GetUdajeProGraf(uzivID uint) ([13]float32, [13]float32, error) {
	var rychlosti [13]float32
	var presnosti [13]float32

	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`WITH dny AS ( SELECT CURRENT_DATE - gs.n AS datum FROM generate_series(0, 12, 1) AS gs (n) ), vsechny_zaznamy AS ( SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene FROM dokoncene WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene FROM dokoncene_procvic WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) ), vypocteny_dny AS ( SELECT datum::date, GREATEST( ( ( SUM(delka_textu) - 10 * SUM(neopravene) ) / SUM(cas)::NUMERIC ) * 60, 0 ) AS rychlost, COALESCE( ( ( SUM(delka_textu) - SUM(neopravene) - COALESCE(SUM(opravene), 0) ) / SUM(delka_textu)::NUMERIC ) * 100, -1 ) AS presnost FROM vsechny_zaznamy GROUP BY datum::date ) SELECT dny.datum, COALESCE(vypocteny_dny.rychlost, -1) AS rychlost, COALESCE(vypocteny_dny.presnost, -1) AS presnost FROM dny LEFT JOIN vypocteny_dny ON dny.datum = vypocteny_dny.datum ORDER BY dny.datum;`, uzivID)
	if err != nil {
		return rychlosti, presnosti, err
	}
	defer rows.Close()

	for i := range 13 {
		rows.Next()

		var rychlost, presnost float32
		var datum date.Date
		err := rows.Scan(&datum, &rychlost, &presnost)
		if err != nil {
			return rychlosti, presnosti, err
		}

		rychlosti[i] = rychlost
		presnosti[i] = presnost
	}
	return rychlosti, presnosti, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var x float32
	err := DB.QueryRow(`WITH vsechny_cviceni AS (SELECT lekce_id, c.id as cviceni_id FROM cviceni c JOIN lekce l ON l.id = c.lekce_id WHERE l.klavesnice = (SELECT klavesnice FROM uzivatel WHERE id = $1) OR l.klavesnice = 'oboje'), moje_dokonceny AS (SELECT DISTINCT d.cviceni_id FROM dokoncene d JOIN vsechny_cviceni vc ON d.cviceni_id = vc.cviceni_id WHERE d.uziv_id = $1) SELECT (SELECT COUNT(*)::float FROM moje_dokonceny) / (SELECT COUNT(*) FROM vsechny_cviceni) as x;`, uzivID).Scan(&x)
	if err != nil {
		return 0, err
	}

	return x * 100, nil
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	email = strings.ToLower(email)

	var uzivID uint
	// kdyby náhodou uživatel už byl dříve zaregistrovaný, smažu všechen jeho progres pomocí WITH x2 a resetnu vsechny sloupce
	err := DB.QueryRow(`WITH id_uzivatele AS ( SELECT id FROM uzivatel WHERE email = $1 ), d AS ( UPDATE dokoncene SET uziv_id = NULL WHERE uziv_id = ( SELECT id FROM id_uzivatele ) ), dp AS ( UPDATE dokoncene_procvic SET uziv_id = NULL WHERE uziv_id = ( SELECT id FROM id_uzivatele ) ), ds AS ( UPDATE dokoncena_prace SET student_id = NULL WHERE student_id = ( SELECT id FROM id_uzivatele ) ) INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email, jmeno = EXCLUDED.jmeno, heslo = EXCLUDED.heslo, klavesnice = DEFAULT, datum = DEFAULT, skolni_jmeno = DEFAULT, smazany = DEFAULT RETURNING id;`, email, jmeno, hesloHash).Scan(&uzivID)
	if err != nil {
		return 0, err
	}
	return uzivID, nil
}

func PridatDokonceneCvic(cvicID, uzivID uint, neopravene int, cas int, delkaTextu int, chybyPismenka map[string]int) error {
	chybyPismenkaJSON, err := json.Marshal(chybyPismenka)
	if err != nil {
		return errors.New("konverze mapy chyb na json se nepovedla")
	}
	_, err = DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, neopravene, cas, delka_textu, chyby_pismenka) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT duplicitni DO NOTHING`, uzivID, cvicID, neopravene, cas, delkaTextu, chybyPismenkaJSON)
	return err
}

func PridatDokonceneProcvic(procvicID, uzivID uint, neopravene int, cas int, delkaTextu int, chybyPismenka map[string]int) error {
	chybyPismenkaJSON, err := json.Marshal(chybyPismenka)

	// pokud je procvic 0 neboli je to test psaní, vložim NULL
	var procvicCislo = sql.NullString{}
	if procvicID != 0 {
		procvicCislo = sql.NullString{String: fmt.Sprintf("%d", procvicID), Valid: true}
	}
	var id = sql.NullInt32{Int32: int32(uzivID), Valid: true}
	if uzivID == 0 {
		id = sql.NullInt32{}
	}
	if err != nil {
		return errors.New("konverze mapy chyb na json se nepovedla")
	}
	_, err = DB.Exec(`INSERT INTO dokoncene_procvic (uziv_id, typ_textu, neopravene, cas, delka_textu, chyby_pismenka) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT duplicitni2 DO NOTHING`, id, procvicCislo, neopravene, cas, delkaTextu, chybyPismenkaJSON)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}

func GetVsechnySlova(pocet int) ([]string, error) {
	var vysledek []string
	var err error

	rows, err := DB.Query(`WITH prvni AS (SELECT slovo FROM slovnik WHERE nahodnost < $1 ORDER BY nahodnost DESC LIMIT $4), druhy AS (SELECT slovo FROM slovnik WHERE nahodnost < $2 ORDER BY nahodnost DESC LIMIT $4), treti AS (SELECT slovo FROM slovnik WHERE nahodnost < $3 ORDER BY nahodnost DESC LIMIT $4) SELECT slovo FROM prvni UNION ALL SELECT slovo FROM druhy UNION ALL SELECT slovo FROM treti;`, mathRand.Float64()+0.0035, mathRand.Float64()+0.0035, mathRand.Float64()+0.0035, (pocet+3)/3) // +0.0035 proto aby tam byl dostatek slov když random hodí 0.0 (když jo, máme 100 slov)
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

	rows, err := DB.Query(`SELECT veta FROM vety ORDER BY RANDOM() LIMIT $1;`, pocet) // tady si asi random muzu dovolit
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

func GetSlovaProLekci(uzivID uint, pismena string, pocet int, tridaID uint) ([]string, error) {
	var vysledek []string
	var rows *sql.Rows
	var err error

	if pismena == "velká písmena (shift)" || pismena == "čísla" || pismena == "interpunkce" {
		var err error
		rows, err = DB.Query(`SELECT slovo FROM slovnik WHERE lekceqwertz_id <= (SELECT id from lekce WHERE pismena = $1) ORDER BY RANDOM() LIMIT $2;`, pismena, pocet)
		if err != nil {
			return vysledek, err
		}
	} else {
		if tridaID == 0 {
			rows, err = DB.Query(`WITH uziv AS ( SELECT klavesnice FROM uzivatel WHERE id = $3 ), qwertz AS ( SELECT slovo FROM slovnik WHERE lekceqwertz_id = ( SELECT id from lekce WHERE pismena = $1 ) AND (SELECT klavesnice FROM uziv) = 'qwertz' ORDER BY RANDOM() LIMIT $2 ), qwerty AS ( SELECT slovo FROM slovnik WHERE lekceqwerty_id = ( SELECT id from lekce WHERE pismena = $1 ) AND (SELECT klavesnice FROM uziv) = 'qwerty' ORDER BY RANDOM() LIMIT $2 ) SELECT * FROM qwertz UNION ALL SELECT * FROM qwerty;`, pismena, pocet, uzivID)
		} else {
			rows, err = DB.Query(`WITH uziv AS ( SELECT klavesnice FROM trida WHERE id = $3 ), qwertz AS ( SELECT slovo FROM slovnik WHERE lekceqwertz_id = ( SELECT id from lekce WHERE pismena = $1 ) AND (SELECT klavesnice FROM uziv) = 'qwertz' ORDER BY RANDOM() LIMIT $2 ), qwerty AS ( SELECT slovo FROM slovnik WHERE lekceqwerty_id = ( SELECT id from lekce WHERE pismena = $1 ) AND (SELECT klavesnice FROM uziv) = 'qwerty' ORDER BY RANDOM() LIMIT $2 ) SELECT * FROM qwertz UNION ALL SELECT * FROM qwerty;`, pismena, pocet, tridaID)
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

	rows, err := DB.Query(`SELECT slovo FROM slovnik_programator ORDER BY RANDOM();`) // tady si taky muzu random dovolit (jen super malo rows)
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

func GetNaucenaPismena(uzivID uint, pismena string, tridaID uint) (string, error) {
	var vysledek strings.Builder

	var rows *sql.Rows
	var err error
	if tridaID == 0 {
		rows, err = DB.Query(`SELECT pismena FROM lekce WHERE id <= (SELECT id from lekce WHERE pismena = $1) AND (klavesnice = COALESCE((SELECT klavesnice FROM uzivatel WHERE id = $2), 'qwertz') OR klavesnice = 'oboje');`, pismena, uzivID)
	} else {
		rows, err = DB.Query(`SELECT pismena FROM lekce WHERE id <= (SELECT id from lekce WHERE pismena = $1) AND (klavesnice = COALESCE((SELECT klavesnice FROM trida WHERE id = $2), 'qwertz') OR klavesnice = 'oboje');`, pismena, tridaID)
	}
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
	row, err := DB.Query(`SELECT * FROM overeni WHERE email = $1;`, email)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func DalSpatnyKod(email string) {
	var pokusy int
	DB.QueryRow(`UPDATE overeni SET pokusy = pokusy - 1 WHERE email = $1 RETURNING pokusy;`, email).Scan(&pokusy)
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
	row, err := DB.Query(`SELECT * FROM zmena_hesla WHERE email = $1;`, email)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
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
	_, err := DB.Exec(`INSERT INTO trida (jmeno, ucitel_id, kod) VALUES ($1, (SELECT id FROM ucitel WHERE uziv_id = $2), $3)`, jmeno, ucitelID, kod)
	return err
}

func GetTrida(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.id, t.jmeno, u.uziv_id AS ucitel_id, t.kod, t.zamknuta, t.smazana, t.klavesnice FROM trida t INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE t.id = $1;`, id)
	if err != nil {
		return trida, err
	}
	err = scan.Row(&trida, row)
	return trida, err
}

func GetTridaByStudentID(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.* FROM trida t INNER JOIN student_a_trida sat ON t.id = sat.trida_id WHERE NOT smazana AND sat.student_id = $1 LIMIT 1;`, id)
	if err != nil {
		return trida, err
	}
	err = scan.Row(&trida, row)
	return trida, err
}

type TridaInfo struct {
	ID            uint   `json:"id" db:"id"`
	Jmeno         string `json:"jmeno" db:"jmeno"`
	Kod           string `json:"kod" db:"kod"`
	Zamknuta      bool   `json:"zamknuta" db:"zamknuta"`
	PocetStudentu int    `json:"pocet_studentu" db:"pocet_studentu"`
	PocetPraci    int    `json:"pocet_praci" db:"pocet_praci"`
}

func GetTridy(ucitelID uint) ([]TridaInfo, error) {
	var tridy []TridaInfo = []TridaInfo{}

	rows, err := DB.Query(`SELECT t.id, t.jmeno, t.kod, t.zamknuta, ( SELECT COUNT(*) FROM uzivatel uz INNER JOIN student_a_trida s ON s.student_id = uz.id WHERE s.trida_id = t.id AND NOT uz.smazany ) as pocet_studentu, ( SELECT COUNT(*) FROM prace p WHERE p.trida_id = t.id AND NOT p.smazana ) as pocet_praci FROM trida t INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE NOT t.smazana AND u.uziv_id = $1;`, ucitelID)
	if err != nil {
		return tridy, err
	}

	err = scan.Rows(&tridy, rows)
	return tridy, err
}

type Student struct {
	ID    uint    `json:"id" db:"id"`
	Jmeno string  `json:"jmeno" db:"skolni_jmeno"`
	Email string  `json:"email" db:"email"`
	CPM   float64 `json:"cpm" db:"cpm"`
}

func GetStudentyZeTridy(tridaID uint) ([]Student, error) {
	var zaci []Student = []Student{}

	rows, err := DB.Query(`WITH cpm_data AS ( SELECT datum, delka_textu, neopravene, cas, uziv_id FROM dokoncene WHERE datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT datum, delka_textu, neopravene, cas, uziv_id FROM dokoncene_procvic WHERE datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) ), cpm AS ( SELECT uziv_id, GREATEST( ( (SUM(delka_textu) - 10 * SUM(neopravene)) / SUM(cas)::NUMERIC ) * 60, 0 ) AS cpm FROM cpm_data GROUP BY uziv_id ) SELECT u.id, u.skolni_jmeno, u.email, COALESCE(mc.cpm, 0) as cpm FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id INNER JOIN trida t ON t.id = s.trida_id FULL OUTER JOIN cpm mc ON mc.uziv_id = u.id WHERE s.trida_id = $1 AND t.smazana = FALSE AND NOT u.smazany GROUP BY u.id, u.skolni_jmeno, u.email, mc.cpm;`, tridaID)
	if err != nil {
		return zaci, err
	}

	err = scan.Rows(&zaci, rows)
	return zaci, err
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

func ZmenitKlavesniciTridy(id uint, klavesnice string) error {
	_, err := DB.Exec(`UPDATE trida SET klavesnice = $1 WHERE id = $2;`, klavesnice, id)
	return err
}

func PrejmenovatStudenta(id uint, skolniJmeno string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET skolni_jmeno = $1 WHERE id = $2;`, skolniJmeno, id)
	return err
}

func PresunoutStudenta(id uint, tridaID int) error {
	_, err := DB.Exec(`UPDATE student_a_trida SET trida_id = $1 WHERE student_id = $2;`, tridaID, id)
	return err
}

func OdebratStudenta(id uint) error {
	_, err := DB.Exec(`DELETE FROM student_a_trida WHERE student_id = $1;`, id)
	return err
}

func GetTridaByKod(kod string) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT * FROM trida WHERE kod = $1 AND NOT smazana;`, kod)
	if err != nil {
		return trida, err
	}
	err = scan.Row(&trida, row)
	return trida, err
}

func GetTridaByUziv(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.* FROM trida t INNER JOIN student_a_trida s ON t.id = s.trida_id WHERE student_id = $1 AND NOT smazana;`, id)
	if err != nil {
		return trida, err
	}
	err = scan.Row(&trida, row)
	return trida, err
}

func ZapsatStudenta(kod string, studentID uint, jmeno string) (int, error) {
	kod = strings.ToUpper(kod)

	var id int
	var smazana bool
	err := DB.QueryRow(`SELECT s.trida_id, t.smazana FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id INNER JOIN trida t ON t.id = s.trida_id WHERE u.id = $1 AND NOT t.smazana;`, studentID).Scan(&id, &smazana)
	if err == nil {
		return 0, errors.New("uz je ve tride")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	if smazana {
		return 0, errors.New("trida je smazana")
	}

	var role int
	err = DB.QueryRow(`SELECT role FROM uzivatel WHERE id = $1`, studentID).Scan(&role)
	if err != nil {
		return 0, err
	}
	if role == 2 {
		return 0, errors.New("jako ucitel nemuzete byt ve tride")
	}

	err = DB.QueryRow(`INSERT INTO student_a_trida (student_id, trida_id) VALUES ($1, (SELECT id FROM trida WHERE kod = $2 AND NOT zamknuta)) ON CONFLICT DO NOTHING RETURNING trida_id;`, studentID, kod).Scan(&id)
	if err != nil {
		return 0, err
	}
	_, err = DB.Exec(`UPDATE uzivatel SET skolni_jmeno = $1 WHERE id = $2;`, jmeno, studentID)
	return id, err
}

func PridatPraci(text string, cas int, tridaID uint) error {
	_, err := DB.Exec(`INSERT INTO prace (trida_id, text, cas) VALUES ($1, $2, $3);`, tridaID, text, cas)
	return err
}

func SmazatPraci(id uint) error {
	_, err := DB.Exec(`UPDATE prace SET smazana = true WHERE id = $1;`, id)
	return err
}

func GetVsechnyPrace(tridaID uint) ([]Prace, error) {
	var prace []Prace = []Prace{}

	rows, err := DB.Query(`WITH soucet_neopravenych AS ( SELECT dp.prace_id, SUM(dp.neopravene) as pocet FROM dokoncena_prace dp INNER JOIN student_a_trida sat ON sat.student_id = dp.student_id INNER JOIN uzivatel u ON u.id = dp.student_id WHERE NOT u.smazany GROUP BY dp.prace_id ), soucet_chyb AS ( SELECT dp.prace_id, SUM((value::VARCHAR(5))::NUMERIC) AS pocet FROM dokoncena_prace dp INNER JOIN student_a_trida sat ON sat.student_id = dp.student_id, jsonb_each_text(dp.chyby_pismenka) WHERE NOT ( SELECT smazany FROM uzivatel WHERE id = dp.student_id ) GROUP BY dp.prace_id ), delky_textu AS ( SELECT dp.prace_id, SUM(dp.delka_textu) AS soucet_delek FROM dokoncena_prace dp INNER JOIN student_a_trida sat ON sat.student_id = dp.student_id INNER JOIN uzivatel u ON u.id = dp.student_id WHERE NOT u.smazany GROUP BY dp.prace_id ), rychlosti AS ( WITH cpmka AS ( SELECT ( ( ( dp.delka_textu - 10 * dp.neopravene ) / dp.cas::FLOAT ) * 60 ) AS cpm, p.id AS prace_id FROM dokoncena_prace dp INNER JOIN prace p ON dp.prace_id = p.id INNER JOIN student_a_trida sat ON sat.student_id = dp.student_id INNER JOIN uzivatel u ON u.id = dp.student_id WHERE NOT u.smazany AND sat.trida_id = $1 ) SELECT prace_id, AVG(GREATEST(cpm, 0)) AS prumerne_cpm FROM cpmka GROUP BY prace_id ), studenti AS ( SELECT dp.prace_id, COUNT(*) as studentu_dokoncilo FROM dokoncena_prace dp INNER JOIN student_a_trida sat ON sat.student_id = dp.student_id INNER JOIN uzivatel u ON u.id = dp.student_id WHERE sat.trida_id = $1 AND NOT u.smazany GROUP BY dp.prace_id ) SELECT p.id, COALESCE(r.prumerne_cpm, -1) AS prumerne_cpm, p.cas, p.datum, p.text, p.trida_id, COALESCE(s.studentu_dokoncilo, 0) AS studentu_dokoncilo, COALESCE( ( ( dt.soucet_delek - COALESCE(sc.pocet, 0) - COALESCE(sn.pocet, 0) ) / dt.soucet_delek::FLOAT ) * 100, -1 ) AS prumerna_presnost FROM soucet_neopravenych sn FULL OUTER JOIN soucet_chyb sc USING (prace_id) FULL OUTER JOIN delky_textu dt USING (prace_id) FULL OUTER JOIN rychlosti r USING (prace_id) FULL OUTER JOIN studenti s USING (prace_id) FULL OUTER JOIN prace p ON p.id = r.prace_id WHERE p.trida_id = $1 AND NOT p.smazana;`, tridaID)
	if err != nil {
		return prace, err
	}
	defer rows.Close()

	err = scan.Rows(&prace, rows)
	return prace, err
}

/* cpmka, presnost, error */
func GetDokoncenePrace(studentID uint) (map[uint]float64, map[uint]float64, error) {
	var cpmka map[uint]float64 = make(map[uint]float64)
	var presnost map[uint]float64 = make(map[uint]float64)

	rows, err := DB.Query(`WITH soucet_chyb AS ( SELECT prace_id, SUM((value::VARCHAR(5))::NUMERIC) AS pocet FROM dokoncena_prace, jsonb_each_text(chyby_pismenka) WHERE student_id = $1 GROUP BY prace_id ) SELECT dp.prace_id, GREATEST( ( ( dp.delka_textu - 10 * dp.neopravene )::FLOAT / dp.cas ) * 60, 0 ) AS cpm, COALESCE( ( ( dp.delka_textu - COALESCE(sc.pocet, 0) - COALESCE(dp.neopravene, 0) ) / dp.delka_textu::FLOAT ) * 100, -1 ) AS prumerna_presnost FROM soucet_chyb sc FULL OUTER JOIN dokoncena_prace dp USING (prace_id) WHERE student_id = $1;`, studentID)
	if err != nil {
		return cpmka, presnost, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var cpm, pres float64
		err := rows.Scan(&id, &cpm, &pres)
		if err != nil {
			return cpmka, presnost, err
		}

		presnost[id] = pres
		cpmka[id] = cpm
	}

	return cpmka, presnost, nil
}

func GetPrace(praceID, studentID uint) (string, int, error) {
	var text string
	var cas int
	err := DB.QueryRow(`SELECT p.text, p.cas FROM prace p INNER JOIN student_a_trida s ON p.trida_id = s.trida_id AND s.student_id = $1 AND p.id = $2 AND NOT p.smazana;`, studentID, praceID).Scan(&text, &cas)
	if err == sql.ErrNoRows {
		return text, cas, errors.New("asi nepatris do teto tridy")
	}
	return text, cas, err
}

func DokoncitPraci(praceID, studentID uint, neopravene int, cas int, delkaTextu int, chybyPismenka map[string]int) (int, error) {
	chybyPismenkaJSON, err := json.Marshal(chybyPismenka)
	if err != nil {
		return 0, err
	}

	var tridaID int
	err = DB.QueryRow(`INSERT INTO dokoncena_prace (prace_id, student_id, neopravene, cas, delka_textu, chyby_pismenka) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ( SELECT trida_id FROM prace WHERE id = $1 );`, praceID, studentID, neopravene, cas, delkaTextu, chybyPismenkaJSON).Scan(&tridaID)
	return tridaID, err
}

type Cviceni2 struct {
	Id      uint   `json:"id"`
	LekceID uint   `json:"lekce_id"`
	Pismena string `json:"pismena"`
}

func GetTypyCviceni(tridaID uint) (map[string][]Cviceni2, error) {
	var mapa = make(map[string][]Cviceni2)
	rows, err := DB.Query(`SELECT typ, pismena, lekce_id, MIN(cislo) FROM ( SELECT c.typ, l.pismena, c.lekce_id, ROW_NUMBER() OVER ( PARTITION BY l.pismena, c.lekce_id ORDER BY c.id ) as cislo FROM cviceni c INNER JOIN lekce l ON c.lekce_id = l.id WHERE l.klavesnice = 'oboje' OR l.klavesnice = ( SELECT klavesnice FROM trida WHERE id = $1 ) ) GROUP BY pismena, typ, lekce_id ORDER BY lekce_id;`, tridaID)
	if err != nil {
		return mapa, err
	}
	defer rows.Close()

	for rows.Next() {
		var typ, pismena string
		var id, lekceID uint
		err := rows.Scan(&typ, &pismena, &lekceID, &id)
		if err != nil {
			return mapa, err
		}

		mapa[typ] = append(mapa[typ], Cviceni2{id, lekceID, pismena})
	}

	return mapa, nil
}

func GetTextZLekce(typ, lekcePismena string) (string, error) {
	rows, _ := DB.Query(`SELECT typ, string_agg( pismena::VARCHAR, ', ' ORDER BY id ) FROM ( SELECT DISTINCT c.typ, l.pismena, l.id FROM cviceni c INNER JOIN lekce l ON l.id = c.lekce_id ) GROUP BY typ;`)
	defer rows.Close()
	return "", nil
}

func CreateSkola(jmeno, email, telefon string) (uint, error) {
	var skolaID uint
	err := DB.QueryRow(`INSERT INTO skola (jmeno, kontaktni_email, kontaktni_telefon) VALUES ($1, $2, $3);`, jmeno, email, telefon).Scan(&skolaID)
	return skolaID, err
}

func GetSkolaByUcitel(uzivID uint) (Skola, error) {
	var skola Skola
	rows, err := DB.Query(`SELECT s.* FROM skola s INNER JOIN ucitel u ON u.skola_id = s.id INNER JOIN uzivatel uz ON uz.id = u.uziv_id WHERE uz.id = $1;`, uzivID)
	if err != nil {
		return skola, err
	}
	err = scan.Row(&skola, rows)
	return skola, err
}

func GetUcitele(skolaID uint) ([]Ucitel, error) {
	var ucitele []Ucitel
	rows, err := DB.Query(`WITH studenti AS ( SELECT * FROM student_a_trida sat INNER JOIN uzivatel uz ON sat.student_id = uz.id WHERE NOT uz.smazany ), studenti_ucitele AS ( SELECT u.id, COUNT(s.student_id) AS pocet_studentu FROM trida t LEFT JOIN studenti s ON s.trida_id = t.id INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE NOT t.smazana GROUP BY u.id ), staty AS ( SELECT ucitel_id, COUNT(*) AS pocet_trid FROM trida WHERE NOT smazana GROUP BY ucitel_id ) SELECT uz.email, pocet_trid, u.id, uz.id AS uziv_id, u.skola_id, SUM(stu.pocet_studentu) AS pocet_zaku FROM ucitel u INNER JOIN uzivatel uz ON uz.id = u.uziv_id LEFT JOIN staty s ON s.ucitel_id = u.id LEFT JOIN studenti_ucitele stu ON stu.id = u.id WHERE u.skola_id = $1 GROUP BY uz.email, pocet_trid, u.id, uz.id, u.skola_id;`, skolaID)
	if err != nil {
		return ucitele, err
	}
	err = scan.Rows(&ucitele, rows)
	return ucitele, err
}

func CreateUcitel(skolaID, uzivID uint) error {
	_, err := DB.Exec(`INSERT INTO ucitel (uziv_id, skola_id) VALUES ($1, $2);`, uzivID, skolaID)
	return err
}

func RemoveUcitelByID(ucitelID uint) error {
	_, err := DB.Exec(`DELETE FROM ucitel WHERE id = $1;`, ucitelID)
	return err
}
