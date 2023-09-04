package databaze

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
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
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Jmeno       string    `json:"jmeno"`
	Heslo       string    `json:"heslo"`
	DayStreak   int32     `json:"daystreak"`
	PosledniDen time.Time `json:"posledniden"`
	Klavesnice  string    `json:"klavesnice"`
}

type Slovnik struct {
	ID    uint   `json:"id"`
	Slovo string `json:"slovo"`
}

type Dokoncene struct {
	ID        uint    `json:"id"`
	UzivID    uint    `json:"uziv_id"`
	CviceniID uint    `json:"cviceni_id"`
	CPM       float32 `json:"cpm"`
	Preklepy  uint    `json:"preklepy"`
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
	rows, err := DB.Queryx(`SELECT a.lekce_id FROM (SELECT lekce_id, COUNT(lekce_id) as pocet_doko FROM dokoncene d INNER JOIN cviceni c ON d.cviceni_id = c.id WHERE uziv_id = $1 GROUP BY lekce_id) a INNER JOIN (SELECT lekce_id, COUNT(lekce_id) as pocet_cvic FROM cviceni GROUP BY lekce_id) b ON a.lekce_id = b.lekce_id WHERE pocet_doko = pocet_cvic;`, uzivID)
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

func GetPreklepyACPM(uzivID uint) ([]float32, []float32, error) {
	var preklepy []float32
	var cpm []float32

	var poslednich int = 10
	rows, err := DB.Queryx(`SELECT preklepy, cpm FROM dokoncene WHERE uziv_id = $1 ORDER BY id DESC LIMIT $2;`, uzivID, poslednich)
	if err != nil {
		return preklepy, cpm, err
	}
	defer rows.Close()

	for rows.Next() {
		var preklep float32
		var cpmko float32
		err := rows.Scan(&preklep, &cpmko)
		if err != nil {
			return preklepy, cpm, err
		}

		preklepy = append(preklepy, float32(preklep))
		cpm = append(cpm, cpmko)
	}
	return preklepy, cpm, nil
}

func DokonceneProcento(uzivID uint) (float32, error) { // TODO predelat na jeden sql
	var pocet int32
	err := DB.QueryRowx(`SELECT COUNT(*) FROM dokoncene WHERE uziv_id = $1;`, uzivID).Scan(&pocet)
	if err != nil {
		return 0, err
	}

	var pocet2 int32
	err = DB.QueryRowx(`SELECT COUNT(*) FROM cviceni;`).Scan(&pocet2)
	if err != nil {
		return 0, err
	}

	return float32(pocet) / float32(pocet2) * 100, nil
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	var uzivID uint
	err := DB.QueryRowx(`INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3);`, email, jmeno, hesloHash).Scan(&uzivID)
	if err != nil {
		return 0, err
	}
	return uzivID, nil
}

func PridatDokonceneCvic(cvicID uint, uzivID uint, cpm float32, preklepy int, cas float32) error {
	if _, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy, cas) VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT dokoncene_pkey DO UPDATE SET cpm = EXCLUDED.cpm, preklepy = EXCLUDED.preklepy, cas = EXCLUDED.cas;`, uzivID, cvicID, cpm, preklepy, cas); err != nil {
		return err
	}
	uziv, err := GetUzivByID(uzivID)
	if err != nil {
		return err
	}
	if uziv.PosledniDen.Format(time.DateOnly) == time.Now().Add(-24*time.Hour).Format(time.DateOnly) {
		if _, err := DB.Exec(`UPDATE uzivatel SET posledniden = $1, daystreak = daystreak + 1 WHERE id = $2 AND posledniden != $1;`, time.Now().Format(time.DateOnly), uzivID); err != nil {
			return err
		}
	} else {
		if _, err := DB.Exec(`UPDATE uzivatel SET posledniden = $1, daystreak = 1 WHERE id = $2;`, time.Now().Format(time.DateOnly), uzivID); err != nil {
			return err
		}
	}
	return nil
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
