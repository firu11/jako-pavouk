package databaze

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

/*
Tahle funkce projde vsechny slovicka ve csv souboru a nasazi do databaze do tabulky slovnik spolu s id lekce ve které je pomocí dosavandne naucenych slovicek možné ho napsat
*/
func PushSlovnik() {
	fmt.Println("jdem na to")

	rows, err := DB.Query(`SELECT pismena FROM lekce;`)
	if err != nil {
		return
	}
	defer rows.Close()

	lekce := []string{}
	for rows.Next() {
		var pis string
		rows.Scan(&pis)
		lekce = append(lekce, pis)
	}

	f, _ := os.Open("C:/Users/Firu/Downloads/lekce.csv")
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	f.Close()

	st := `INSERT INTO slovnik (slovo, lekce_id) VALUES `

	var pismenka string
	var berem bool
	for _, v := range records {
		pismenka = ""
		for i, p := range lekce {
			pismenka += p
			berem = true
			for _, x := range v[0] {
				if !strings.Contains(pismenka, string(x)) {
					berem = false
					break
				}
			}
			if berem {
				st += fmt.Sprintf(`('%s', %v), `, v[0], i+1)
				break
			}
		}
	}
	st = st[:len(st)-2]
	st += ";"
	fmt.Println(st)
	if _, err = DB.Exec(`DELETE FROM slovnik WHERE true;`); err != nil {
		panic(err)
	}

	if _, err = DB.Exec(st); err != nil {
		panic(err)
	}
}
