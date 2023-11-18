package databaze

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func PushCviceni() {
	m := make(map[int][3]int) /* lekceID - [nova, naucena, slova] */
	m[1] = [3]int{4, 0, 0}
	m[2] = [3]int{3, 2, 0}
	m[3] = [3]int{3, 2, 0}
	m[4] = [3]int{3, 2, 1}
	m[5] = [3]int{2, 1, 2}
	m[6] = [3]int{2, 1, 2}
	m[7] = [3]int{2, 1, 2}
	m[8] = [3]int{2, 1, 2}
	m[9] = [3]int{2, 1, 2}
	m[10] = [3]int{3, 1, 2}
	m[11] = [3]int{2, 1, 2}
	m[12] = [3]int{2, 1, 2}
	m[13] = [3]int{3, 1, 2}
	m[14] = [3]int{3, 1, 2}
	m[15] = [3]int{2, 1, 2}
	m[16] = [3]int{2, 1, 2}
	m[17] = [3]int{2, 1, 2}
	m[18] = [3]int{3, 1, 2}
	m[19] = [3]int{1, 1, 2}
	m[20] = [3]int{2, 1, 0}

	query := `INSERT INTO cviceni (lekce_id, typ) VALUES`

	for cislo, cviceni := range m {
		for typ, cv := range cviceni {
			for pocet := 0; pocet < cv; pocet++ {
				var typString string
				if typ == 0 {
					typString = "nova"
				} else if typ == 1 {
					typString = "naucena"
				} else {
					typString = "slova"
				}
				query += fmt.Sprintf(" (%d, '%s'),", cislo, typString)
			}
		}
	}

	query = query[:len(query)-1]
	query += ";"

	_, err := DB.Exec(`
		DROP TABLE IF EXISTS cviceni;
		CREATE TABLE cviceni (
        	id SERIAL PRIMARY KEY,
        	typ VARCHAR(20) DEFAULT 'nova',
        	lekce_id INT,
        	FOREIGN KEY (lekce_id) REFERENCES lekce(id)
    	);
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = DB.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

type lekcos struct {
	Id      uint   `json:"id"`
	Pismena string `json:"pismena"`
}

/*
Tahle funkce projde vsechny slovicka ve csv souboru a nasazi do databaze do tabulky slovnik spolu s id lekce ve které je pomocí dosavandne naucenych slovicek možné ho napsat
*/
func PushSlovnik() {
	fmt.Println("jdem na to")

	DBConnect()

	_, err := DB.Exec(`
		DROP TABLE IF EXISTS slovnik;
		CREATE TABLE
    		IF NOT EXISTS slovnik (
        		id SERIAL PRIMARY KEY,
        		slovo VARCHAR(50),
        		lekceQWERTZ_id INT,
        		lekceQWERTY_id INT
    	);
	`)
	if err != nil {
		log.Panic(err)
	}

	rowsZ, err1 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwertz' OR klavesnice = 'oboje' ORDER BY id ASC;`)
	rowsY, err2 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwerty' OR klavesnice = 'oboje' ORDER BY id ASC;`)

	if err1 != nil || err2 != nil {
		return
	}
	defer rowsZ.Close()
	defer rowsY.Close()

	lekceZ := []lekcos{}
	for rowsZ.Next() {
		var l lekcos = lekcos{}
		rowsZ.Scan(&l.Id, &l.Pismena)
		lekceZ = append(lekceZ, l)
	}

	lekceY := []lekcos{}
	for rowsY.Next() {
		var l lekcos = lekcos{}
		rowsY.Scan(&l.Id, &l.Pismena)
		lekceY = append(lekceY, l)
	}

	log.Println(lekceZ, lekceY)

	f, err := os.Open("./slovnik.csv")
	if err != nil {
		log.Println("spatna cesta k csvcku")
		return
	}
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	f.Close()

	st := `INSERT INTO slovnik (slovo, lekceQWERTZ_id, lekceQWERTY_id) VALUES `

	var pismenkaZ string
	var pismenkaY string
	var indexZ int
	var indexY int
	for _, v := range records {
		pismenkaZ = ""
		indexZ = -1
		for _, p := range lekceZ {
			if p.Pismena == "Zbylá diakritika" {
				pismenkaZ += "óďťň"
			} else if p.Pismena != "Velká písmena (Shift)" {
				pismenkaZ += p.Pismena
			}
			if obsahujeJenOKPismena(v[0], pismenkaZ) {
				indexZ = int(p.Id)
			}
			if indexZ != -1 {
				break
			}
		}
		pismenkaY = ""
		indexY = -1
		for _, p := range lekceY {
			if p.Pismena == "Zbylá diakritika" {
				pismenkaY += "óďťň"
			} else if p.Pismena != "Velká písmena (Shift)" {
				pismenkaY += p.Pismena
			}
			if obsahujeJenOKPismena(v[0], pismenkaY) {
				indexY = int(p.Id)
			}
			if indexY != -1 {
				break
			}
		}

		st += fmt.Sprintf(`('%s', %v, %v), `, v[0], indexZ, indexY)
	}
	st = st[:len(st)-2]
	st += ";"
	/* fmt.Println(st) */

	if _, err := DB.Exec(st); err != nil {
		panic(err)
	}
}

func obsahujeJenOKPismena(slovo string, pismena string) bool {
	ok := true
	for _, x := range slovo {
		if !strings.Contains(pismena, string(x)) {
			ok = false
			break
		}
	}
	return ok
}
