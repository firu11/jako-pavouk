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
	Klavesnice  bool      `json:"klavesnice"` /* TRUE - QWERTZ, FALSE - QWERTY */
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

	var uziv Uzivatel
	var err error
	var rows *sqlx.Rows

	if uzivID != 0 {
		uziv, err = GetUzivByID(uzivID)
		if err != nil {
			return lekce, err
		}
		if uziv.Klavesnice {
			rows, err = DB.Queryx(`SELECT * FROM lekceQWERTZ;`)
			if err != nil {
				return lekce, err
			}
		} else {
			rows, err = DB.Queryx(`SELECT * FROM lekceQWERTY;`)
			if err != nil {
				return lekce, err
			}
		}
	} else {
		rows, err = DB.Queryx(`SELECT * FROM lekceQWERTZ;`)
		if err != nil {
			return lekce, err
		}
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
	rows, err := DB.Queryx(`SELECT id FROM lekceQWERTZ l WHERE 0 = (SELECT COUNT(*) FROM cviceni c WHERE c.lekce_id = l.id) - (SELECT COUNT(*) FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id AND d.uziv_id = $1 AND l.id = c.lekce_id);`, uzivID)
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
	err := DB.QueryRowx(`SELECT id FROM lekceQWERTZ WHERE pismena = $1;`, pismena).Scan(&id)
	if err != nil {
		err2 := DB.QueryRowx(`SELECT id FROM lekceQWERTY WHERE pismena = $1;`, pismena).Scan(&id)
		if err2 != nil {
			return 0, err2
		}
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

	uziv, err := GetUzivByID(uzivID)
	if err != nil {
		return cviceni, err
	}

	var rows *sqlx.Rows
	if uziv.Klavesnice {
		rows, err = DB.Queryx(`SELECT id, typ FROM cviceni WHERE lekce_id = (SELECT id FROM lekceQWERTZ where pismena = $1 LIMIT 1);`, pismena)
		if err != nil {
			return cviceni, err
		}
	} else {
		rows, err = DB.Queryx(`SELECT id, typ FROM cviceni WHERE lekce_id = (SELECT id FROM lekceQWERTY where pismena = $1 LIMIT 1);`, pismena)
		if err != nil {
			return cviceni, err
		}
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
	if err != nil {
		return err
	}
	return nil
}

func ZmenitKlavesnici(id uint, novaKlavesnice string) error {
	var hodnotaDoDB bool
	if novaKlavesnice == "QWERTZ" {
		hodnotaDoDB = true
	} else if novaKlavesnice == "QWERTY" {
		hodnotaDoDB = false
	}
	_, err := DB.Exec(`UPDATE uzivatel SET klavesnice = $1 WHERE id = $2;`, hodnotaDoDB, id)
	if err != nil {
		return err
	}
	return nil
}

func PrejmenovatUziv(id uint, noveJmeno string) error {
	_, err := DB.Exec(`UPDATE uzivatel SET jmeno = $1 WHERE id = $2;`, noveJmeno, id)

	return err
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

func DokonceneProcento(uzivID uint) (float32, error) {
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
	_, err := DB.Exec(`INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3)`, email, jmeno, hesloHash)
	if err != nil {
		return 0, err
	}
	uziv, err := GetUzivByEmail(email)
	if err != nil {
		return 0, err
	}
	return uziv.ID, nil
}

func PridatDokonceneCvic(cvicID uint, uzivID uint, cpm float32, preklepy int) error {
	if _, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy) VALUES ($1, $2, $3, $4);`, uzivID, cvicID, cpm, preklepy); err != nil {
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

	uziv, err := GetUzivByID(uzivID)
	if err != nil {
		return vysledek, err
	}

	var rows *sqlx.Rows
	if uziv.Klavesnice {
		rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwertz_id = (SELECT id from lekceQWERTZ WHERE pismena = $1);`, pismena)
		if err != nil {
			return vysledek, err
		}
	} else {
		rows, err = DB.Queryx(`SELECT slovo FROM slovnik WHERE lekceqwerty_id = (SELECT id from lekceQWERTY WHERE pismena = $1);`, pismena)
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

func GetNaucenaPismena(pismena string) (string, error) {
	var vysledek string
	rows, err := DB.Queryx(`SELECT pismena FROM lekce WHERE id <= (SELECT id from lekce WHERE pismena = $1);`, pismena)
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
