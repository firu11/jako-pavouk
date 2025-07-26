package databaze

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	mathRand "math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/blockloop/scan/v2"
	"github.com/rickb777/date"
)

var RegexJmeno *regexp.Regexp
var MaxCisloZaJmeno int // 10_000
var poslednich int = 15

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
	row, err := DB.Query(`SELECT uz.*, CASE WHEN u.smazany THEN 0 ELSE COALESCE(u.skola_id, 0) END AS skola_id FROM uzivatel uz LEFT JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.id = $1 AND NOT uz.smazany;`, uzivID)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv Uzivatel
	row, err := DB.Query(`SELECT uz.*, CASE WHEN u.smazany THEN 0 ELSE COALESCE(u.skola_id, 0) END AS skola_id FROM uzivatel uz LEFT JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.email = $1 AND NOT uz.smazany;`, email)
	if err != nil {
		return uziv, err
	}
	err = scan.Row(&uziv, row)
	return uziv, err
}

func GetUzivByJmeno(jmeno string) (Uzivatel, error) {
	var uziv Uzivatel
	row, err := DB.Query(`SELECT uz.*, CASE WHEN u.smazany THEN 0 ELSE COALESCE(u.skola_id, 0) END AS skola_id FROM uzivatel uz LEFT JOIN ucitel u ON u.uziv_id = uz.id WHERE uz.jmeno = $1 AND NOT uz.smazany;`, jmeno)
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

	err := DB.QueryRow(`WITH default_hodnoty AS ( SELECT 0 AS neopravene, 0 AS delka_textu, 0 AS cas, CURRENT_DATE AS datum, 0 AS opravene, '{}'::jsonb AS chyby_pismenka ), vsechny_za_dva_tydny AS ( SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene, chyby_pismenka FROM dokoncene WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT neopravene, delka_textu, cas, datum, ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ) AS opravene, chyby_pismenka FROM dokoncene_procvic WHERE uziv_id = $1 AND datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT * FROM default_hodnoty ), soucty_pismenek AS ( SELECT key AS pismeno, SUM(value::NUMERIC) AS soucet FROM vsechny_za_dva_tydny, jsonb_each_text(chyby_pismenka) GROUP BY key UNION ALL SELECT '' AS pismeno, 0 AS soucet WHERE NOT EXISTS ( SELECT 1 FROM vsechny_za_dva_tydny, jsonb_each_text(chyby_pismenka) ) ), vsechny AS ( SELECT delka_textu, cas, datum::date FROM dokoncene WHERE uziv_id = $1 UNION ALL SELECT delka_textu, cas, datum::date FROM dokoncene_procvic WHERE uziv_id = $1 ), soucty_dnes AS ( SELECT SUM(cas) AS cas_dnes, SUM(delka_textu) AS napsanych_pismen_dnes FROM vsechny WHERE datum = CURRENT_DATE ), soucty_dva_tydny AS ( SELECT SUM(cas) AS cas_dva_tydny, SUM(delka_textu) AS napsanych_pismen_dva_tydny FROM vsechny WHERE datum > CURRENT_DATE - MAKE_INTERVAL(days => 14) ), soucty_celkem AS ( SELECT SUM(cas) AS cas_celkem, SUM(delka_textu) AS napsanych_pismen_celkem FROM vsechny ) SELECT GREATEST( ( (SUM(delka_textu) - 10 * SUM(neopravene)) / GREATEST(SUM(cas)::NUMERIC, 1) ) * 60, 0 ) AS rychlost, COALESCE( ( ( SUM(delka_textu) - SUM(neopravene) - COALESCE(SUM(opravene), 0) ) / GREATEST(SUM(delka_textu)::NUMERIC, 1) ) * 100, -1 ) AS presnost, COALESCE(jsonb_object_agg(pismeno, soucet), '{}'::jsonb) AS chyby_pismenka, COALESCE(max(cas_dnes), 0) AS cas_dnes, COALESCE(max(cas_dva_tydny), 0) AS cas_dva_tydny, COALESCE(max(cas_celkem), 0) AS cas_celkem, COALESCE(max(napsanych_pismen_dnes), 0) AS napsanych_pismen_dnes, COALESCE(max(napsanych_pismen_dva_tydny), 0) AS napsanych_pismen_dva_tydny, COALESCE(max(napsanych_pismen_celkem), 0) AS napsanych_pismen_celkem FROM soucty_pismenek, vsechny_za_dva_tydny, soucty_dnes, soucty_dva_tydny, soucty_celkem;`, uzivID).Scan(&rychlost, &presnost, &chybyPismenkaJsonb, &cas[0], &cas[1], &cas[2], &napsanychPismen[0], &napsanychPismen[1], &napsanychPismen[2])
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
	// ignoruje duplicitní inserty které jsou méně než 10s od sebe. po 10s už povoluje doplicity
	_, err = DB.Exec(`WITH recent AS (SELECT 1 FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2 AND neopravene = $3 AND cas = $4 AND delka_textu = $5 AND now() - datum <= interval '10 seconds' LIMIT 1) INSERT INTO dokoncene (uziv_id, cviceni_id, neopravene, cas, delka_textu, chyby_pismenka) SELECT $1, $2, $3, $4, $5, $6 WHERE NOT EXISTS (SELECT 1 FROM recent);`, uzivID, cvicID, neopravene, cas, delkaTextu, chybyPismenkaJSON)
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
	_, err = DB.Exec(`WITH recent AS (SELECT 1 FROM dokoncene_procvic WHERE uziv_id = $1 AND typ_textu = $2 AND neopravene = $3 AND cas = $4 AND delka_textu = $5 AND now() - datum <= interval '10 seconds' LIMIT 1) INSERT INTO dokoncene_procvic (uziv_id, typ_textu, neopravene, cas, delka_textu, chyby_pismenka) SELECT $1, $2, $3, $4, $5, $6 WHERE NOT EXISTS (SELECT 1 FROM recent);`, id, procvicCislo, neopravene, cas, delkaTextu, chybyPismenkaJSON)
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
