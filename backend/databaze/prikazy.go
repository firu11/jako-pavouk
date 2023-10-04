package databaze

import (
	"errors"
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

type Slovnik struct {
	ID    uint   `json:"id"`
	Slovo string `json:"slovo"`
}

type Dokoncene struct {
	ID        uint      `json:"id"`
	UzivID    uint      `json:"uziv_id"`
	CviceniID uint      `json:"cviceni_id"`
	CPM       float32   `json:"cpm"`
	Preklepy  uint      `json:"preklepy"`
	Cas       int       `json:"cas"`
	Datum     date.Date `json:"datum"`
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

func GetDokonceneCvicVLekci(uzivID uint, lekceID uint, pismena string) ([]int32, error) {
	var cviceniIDs []int32 = []int32{}
	var rows *sqlx.Rows
	var err error

	if pismena != "" {
		lekceID, err = GetLekceIDbyPismena(pismena)
		if err != nil {
			return cviceniIDs, err
		}
	}
	rows, err = DB.Queryx(`SELECT d.cviceni_id FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1 AND uziv_id = $2;`, lekceID, uzivID)
	if err != nil {
		return cviceniIDs, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int32
		if err = rows.Scan(&id); err != nil {
			return cviceniIDs, err
		}
		cviceniIDs = append(cviceniIDs, id)
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

/* preklepy, cpm, daystreak, cas */
func GetUdaje(uzivID uint) ([]float32, []float32, int, float32, error) {
	var preklepy []float32
	var cpm []float32
	var daystreak int = 0
	var celkovyCas float32 = 0

	var poslednich int = 10
	rows, err := DB.Queryx(`SELECT preklepy, cpm FROM dokoncene WHERE uziv_id = $1 ORDER BY den DESC LIMIT $2;`, uzivID, poslednich)
	if err != nil {
		return preklepy, cpm, daystreak, celkovyCas, err
	}
	defer rows.Close()

	for rows.Next() {
		var preklep float32
		var cpmko float32
		err := rows.Scan(&preklep, &cpmko)
		if err != nil {
			return preklepy, cpm, daystreak, celkovyCas, err
		}

		preklepy = append(preklepy, float32(preklep))
		cpm = append(cpm, cpmko)
	}

	rows, err = DB.Queryx(`SELECT den, cas FROM dokoncene WHERE uziv_id = $1 ORDER BY den DESC;`, uzivID)
	if err != nil {
		return preklepy, cpm, daystreak, celkovyCas, err
	}
	defer rows.Close()

	var dny []date.Date
	for rows.Next() {
		var c float32
		var d date.Date
		if err := rows.Scan(&d, &c); err != nil {
			return preklepy, cpm, daystreak, celkovyCas, err
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

	return preklepy, cpm, daystreak, celkovyCas, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var x float32
	err := DB.QueryRowx(`SELECT cast((SELECT COUNT(*) FROM dokoncene WHERE uziv_id = $1) as float) / (SELECT COUNT(*) FROM cviceni) as x;`, uzivID).Scan(&x)
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

func PridatDokonceneCvic(cvicID uint, uzivID uint, cpm float32, preklepy int, cas float32) error {
	_, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy, cas) VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT unikatni DO UPDATE SET cpm = EXCLUDED.cpm, preklepy = EXCLUDED.preklepy, cas = EXCLUDED.cas;`, uzivID, cvicID, cpm, preklepy, cas)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}

func GetSlovaProLekci(uzivID uint, pismena string) ([]string, error) {
	var vysledek []string

	var k string
	err := DB.QueryRowx(`SELECT klavesnice FROM uzivatel WHERE id = $1;`, uzivID).Scan(&k)
	if err != nil {
		return vysledek, err
	}

	var rows *sqlx.Rows
	if k == "qwertz" {
		rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwertz_id = (SELECT id from lekce WHERE pismena = $1);`, pismena)
	} else {
		rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwerty_id = (SELECT id from lekce WHERE pismena = $1);`, pismena)
	}
	if err != nil {
		return vysledek, err
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
	var vysledek string
	rows, err := DB.Queryx(`SELECT pismena FROM lekce WHERE id <= (SELECT id from lekce WHERE pismena = $1) AND (klavesnice = COALESCE((SELECT klavesnice FROM uzivatel WHERE id = $2), 'qwertz') OR klavesnice = 'oboje');`, pismena, uzivID)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	var pismenaJedny string
	for rows.Next() {
		pismenaJedny = ""
		err := rows.Scan(&pismenaJedny)
		if err != nil {
			return vysledek, err
		}

		vysledek += pismenaJedny
	}

	return vysledek, nil
}

func CreateNeoverenyUziv(email string, hesloHASH string, jmeno string, kod string, cas int64) error {
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

func SmazatNeoverenyPoLimitu() error {
	_, err := DB.Exec(`DELETE FROM overeni WHERE cas < $1;`, time.Now().Unix())
	return err
}
