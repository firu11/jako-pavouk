package databaze

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Tahle funkce projde vsechny slovicka ve csv souboru a nasazi do databaze do tabulky slovnik spolu s id lekce ve které je pomocí dosavandne naucenych slovicek možné ho napsat
*/
func PushSlovnik() {
	fmt.Println("jdem na to")

	_, err := DB.Exec(`
		DROP TABLE IF EXISTS slovnik;
		CREATE TABLE
    		IF NOT EXISTS slovnik (
        		id SERIAL PRIMARY KEY,
        		slovo VARCHAR(50),
        		lekceqwertz_id INT,
        		lekceqwerty_id INT
    	);
	`)
	if err != nil {
		log.Panic(err)
	}

	rowsZ, err1 := DB.Query(`SELECT pismena FROM lekceQWERTZ;`)
	rowsY, err2 := DB.Query(`SELECT pismena FROM lekceQWERTY;`)

	if err1 != nil || err2 != nil {
		return
	}
	defer rowsZ.Close()
	defer rowsY.Close()

	lekceZ := []string{}
	for rowsZ.Next() {
		var pis string
		rowsZ.Scan(&pis)
		lekceZ = append(lekceZ, pis)
	}

	lekceY := []string{}
	for rowsY.Next() {
		var pis string
		rowsY.Scan(&pis)
		lekceY = append(lekceY, pis)
	}

	f, _ := os.Open("C:/Users/Firu/Downloads/lekce.csv")
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	f.Close()

	st := `INSERT INTO slovnik (slovo, lekceQWERTZ_id, lekceQWERTY_id) VALUES `

	var pismenkaZ string
	var pismenkaY string
	var indexZ int
	var indexY int
	for _, v := range records {
		log.Println(v)
		pismenkaZ = ""
		pismenkaY = ""
		indexZ = -1
		indexY = -1
		for i, p := range lekceZ {
			pismenkaZ += p
			if obsahujeJenOKPismena(v[0], pismenkaZ) {
				indexZ = i
			}
			if indexZ != -1 {
				break
			}
		}

		for i, p := range lekceY {
			pismenkaY += p
			if obsahujeJenOKPismena(v[0], pismenkaY) {
				indexY = i
			}
			if indexY != -1 {
				break
			}
		}

		st += fmt.Sprintf(`('%s', %v, %v), `, v[0], indexZ+1, indexY+1)
	}
	st = st[:len(st)-2]
	st += ";"
	fmt.Println(st)

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
