package databaze

import (
	"log"
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
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Jmeno string `json:"jmeno"`
	Heslo string `json:"heslo"`
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

func GetLekce() ([][]Lekce, error) {
	rows, err := DB.Query(`SELECT * FROM lekce;`)
	if err != nil {
		return [][]Lekce{}, nil
	}
	defer rows.Close()

	var lekce [][]Lekce
	var skupina []Lekce

	var cisloSkupiny uint = 1
	for rows.Next() {
		jednaLekce := Lekce{}
		var skup uint
		err := rows.Scan(&jednaLekce.ID, &jednaLekce.Pismena, &skup)
		if err != nil {
			log.Fatal(err)
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
	lekce, err := GetLekce()
	if err != nil {
		return []int32{}, err
	}
	var lekce_ids []int32 = []int32{}
	for _, skupina := range lekce {
		for _, lekce := range skupina {
			// zjistim kolik ma kazda lekce cviceni
			var pocet int
			err := DB.QueryRow(`SELECT COUNT(*) FROM cviceni WHERE lekce_id = $1;`, lekce.ID).Scan(&pocet)
			if err != nil {
				return []int32{}, err
			}

			cvicVLekci, err := GetDokonceneCvicVLekci(lekce.ID)
			if err != nil {
				return []int32{}, err
			}

			if pocet == len(cvicVLekci) && pocet != 0 {
				lekce_ids = append(lekce_ids, int32(lekce.ID))
			}
		}
	}

	return lekce_ids, nil //ted mam dokonceny cviceni potrebuju zjistit jestli to je vsechno v ty lekce TODO
}

func GetDokonceneCvicVLekci(lekceID uint) ([]int32, error) {
	rows, err := DB.Query(`SELECT d.cviceni_id FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1;`, lekceID)
	if err != nil {
		return []int32{}, err
	}
	defer rows.Close()

	var cviceni_ids []int32 = []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		cviceni_ids = append(cviceni_ids, id)
	}

	return cviceni_ids, nil
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
	rows, err := DB.Query(`SELECT id, typ FROM cviceni WHERE lekce_id = $1;`, lekceID)
	if err != nil {
		return []Cviceni{}, err
	}
	defer rows.Close()

	var cviceni []Cviceni
	for rows.Next() {
		jednoCviceni := Cviceni{}
		err := rows.Scan(&jednoCviceni.ID, &jednoCviceni.Typ)
		if err != nil {
			log.Fatal(err)
		}

		cviceni = append(cviceni, jednoCviceni)
	}

	return cviceni, nil
}

func GetCviceniVLekciByPismena(pismena string) ([]Cviceni, error) {
	cviceni := []Cviceni{}
	id, err := GetLekceIDbyPismena(pismena)
	if err != nil {
		return cviceni, err
	}
	cviceni, err = GetCviceniVLekciByID(id)
	return cviceni, err
}

func GetUzivByID(id uint) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRow(`SELECT id, email, jmeno, heslo FROM uzivatel WHERE id = $1;`, id).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv = Uzivatel{}
	err := DB.QueryRow(`SELECT id, email, jmeno, heslo FROM uzivatel WHERE email = $1;`, email).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo)
	return uziv, err
}

func GetPreklepyACPM(uzivID uint) ([]float32, []float32, error) {
	var posledni int = 10
	rows, err := DB.Query(`SELECT preklepy, cpm FROM dokoncene WHERE uziv_id = $1 ORDER BY id DESC LIMIT $2;`, uzivID, posledni)
	if err != nil {
		return []float32{}, []float32{}, err
	}
	defer rows.Close()

	var preklepy []float32 = []float32{}
	var cpm []float32 = []float32{}

	for rows.Next() {
		var preklep float32
		var cpmko float32
		err := rows.Scan(&preklep, &cpmko)
		if err != nil {
			log.Fatal(err)
		}

		preklepy = append(preklepy, float32(preklep))
		cpm = append(cpm, cpmko)
	}

	return preklepy, cpm, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var pocet int32
	err := DB.QueryRow(`SELECT COUNT(*) FROM dokoncene WHERE uziv_id = $1;`, uzivID).Scan(&pocet)
	if err != nil {
		return 0, err
	}

	var pocet2 int32
	err = DB.QueryRow(`SELECT COUNT(*) FROM cviceni;`).Scan(&pocet2)
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
	_, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy) VALUES ($1, $2, $3, $4);`, uzivID, cvicID, cpm, preklepy)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}
