package main

import (
	"backend/utils"
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var minDelkaTextu int = 180

type druh struct {
	ID        int
	Kategorie string
	Obtiznost int
}

var druhyTextu = map[string]druh{
	"Zeměpis":               {1, "Naučné", 2},
	"Dějepis":               {2, "Naučné", 3},
	"Obtížnější pohádky":    {3, "Zábavné", 3},
	"Pohádky":               {4, "Zábavné", 1},
	"Robinson Crusoe":       {5, "Knihy", 1},
	"Gulliverovy cesty":     {6, "Knihy", 1},
	"Ostře sledované vlaky": {7, "Knihy", 2},
}

type lekcos struct {
	Id      uint   `json:"id"`
	Pismena string `json:"pismena"`
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Nenašel jsem soubor .env v /backend.")
	}

	fmt.Printf("Připojuješ se na %s (.env)\n", os.Getenv("DB_JMENO"))

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"), os.Getenv("DB_JMENO"))
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Databaze se pokazila" + err.Error())
	}

	for {
		fmt.Print("(s)lovnik / (p)ohadky věty / (t)exty / (q)uit: ")
		var input string
		fmt.Scan(&input)
		fmt.Println()
		if input == "s" {
			PushSlovnik()
		} else if input == "p" {
			PushRandomVety()
		} else if input == "t" {
			PushTexty()
		} else if input == "q" {
			break
		} else if input == "x" {
			//fmt.Println(pushKnihy())
			break
		}
	}

	fmt.Println("\nHotovo!")
}

func pushKnihy() string {
	fmt.Println("Jdem na Knihy")
	fmt.Println("-----------------------------")

	soubory, err := os.ReadDir("./knihy/")
	if err != nil {
		panic(err)
	}

	var st string = `INSERT INTO texty (cislo, jmeno, typ, txt, delka) VALUES `
	var pocet int
	for _, v := range soubory {
		if v.IsDir() {
			continue
		}
		f, err := os.Open("./knihy/" + v.Name())
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)

		scanner.Scan()
		jmeno := strings.TrimSpace(scanner.Text())
		idKnihy := druhyTextu[jmeno].ID
		if idKnihy == 0 {
			panic("nejaka divna kniha, nemuzu najit id")
		}

		var i int = 1
		for scanner.Scan() {
			if scanner.Text() == "" || scanner.Text()[0] == '%' {
				continue
			}

			text := strings.TrimSpace(scanner.Text())
			st += fmt.Sprintf(`(%d, '%s', %d, '%s', %d), `, i, fmt.Sprintf("Část %v", i), idKnihy, utils.UpravaTextu(text), utf8.RuneCountInString(text))

			i++
			pocet++
		}

		f.Close()
	}

	st = st[:len(st)-2]
	st += ";"

	fmt.Printf("%v textů z knih jde do DB", pocet)

	return st
}

func PushSlovnik() {
	fmt.Println("\nJdem na slovník")
	fmt.Println("-----------------------------")

	rowsZ, err1 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwertz' OR klavesnice = 'oboje' ORDER BY id ASC;`)
	rowsY, err2 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwerty' OR klavesnice = 'oboje' ORDER BY id ASC;`)

	if err1 != nil || err2 != nil {
		fmt.Println("Připojení k databázi se nezdařilo...")
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

	fmt.Println("Lekce z DB načteny")

	f, err := os.Open("texty/slovnik.txt")
	if err != nil {
		fmt.Println("spatna cesta k slovniku")
		return
	}
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	f.Close()

	fmt.Println("Slova načteny")

	st := `INSERT INTO slovnik (slovo, lekceQWERTZ_id, lekceQWERTY_id) VALUES `

	var pismenkaZ string
	var pismenkaY string
	var indexZ int
	var indexY int
	for _, v := range records {
		pismenkaZ = ""
		indexZ = -1
		for _, p := range lekceZ {
			if p.Pismena == "zbylá diakritika" {
				pismenkaZ += "óďťň"
			} else if p.Pismena != "velká písmena (Shift)" && p.Pismena != "závorky" && p.Pismena != "operátory" && p.Pismena != "čísla" && p.Pismena != "interpunkce" {
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
			if p.Pismena == "zbylá diakritika" {
				pismenkaY += "óďťň"
			} else if p.Pismena != "velká písmena (shift)" && p.Pismena != "závorky" && p.Pismena != "operátory" && p.Pismena != "čísla" && p.Pismena != "interpunkce" {
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

	_, err = DB.Exec(`
		DROP TABLE IF EXISTS slovnik;
		CREATE TABLE
    		IF NOT EXISTS slovnik (
        		id SERIAL PRIMARY KEY,
        		slovo VARCHAR(50),
        		lekceQWERTZ_id INT,
        		lekceQWERTY_id INT,
				nahodnost NUMERIC(18, 17) NOT NULL DEFAULT RANDOM()
    	);
		CREATE INDEX idx_random ON slovnik (nahodnost);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v slov jde do DB\n", len(records))

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

func PushRandomVety() {
	fmt.Println("\nJdem na vety")
	fmt.Println("-----------------------------")

	f, err := os.Open("texty/vety.txt")
	if err != nil {
		fmt.Println("spatna cesta k souboru")
		return
	}
	defer f.Close()

	fmt.Println("Pohádky načteny")

	var st string = `INSERT INTO vety (veta, delka) VALUES `

	scanner := bufio.NewScanner(f)
	var pohadky []string
	for scanner.Scan() {
		pohadky = append(pohadky, scanner.Text())
	}

	sort.Slice(pohadky, func(i, j int) bool {
		return utf8.RuneCountInString(pohadky[i]) < utf8.RuneCountInString(pohadky[j])
	})

	var pocet int
	delky := make(map[int]int)
	for _, poh := range pohadky {
		delka := utf8.RuneCountInString(poh)
		_, ok := delky[delka]
		if ok {
			delky[utf8.RuneCountInString(poh)]++
		} else {
			delky[utf8.RuneCountInString(poh)] = 1
		}
		if delka <= 80 {
			pocet++
			st += fmt.Sprintf(`('%s', %v), `, poh, delka)
		}
	}

	st = st[:len(st)-2]
	st += ";"

	_, err = DB.Exec(`
		DROP TABLE IF EXISTS vety;
		CREATE TABLE
    		IF NOT EXISTS vety (
        		id SERIAL PRIMARY KEY,
        		veta TEXT NOT NULL,
				delka INT NOT NULL
    	);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v random vět jde do DB", pocet)

	if _, err := DB.Exec(st); err != nil {
		panic(err)
	}
}

func PushTexty() {
	fmt.Print("\n")
	dejepisQuery := pushDejepis()
	fmt.Print("\n\n")
	zemepisQuery := pushZemepis()
	fmt.Print("\n\n")
	pohadkyQuery := pushPohadky()
	fmt.Print("\n\n")
	knihyQuery := pushKnihy()
	fmt.Print("\n\n")

	var prepare string = `
		DROP TABLE IF EXISTS texty;
		CREATE TABLE texty (
    		id SERIAL PRIMARY KEY,
    		cislo INTEGER,
    		jmeno VARCHAR(255),
    		typ INTEGER REFERENCES druhy_textu (id),
    		txt TEXT,
    		delka INTEGER
		);
		INSERT INTO druhy_textu (id, jmeno, kategorie, obtiznost) VALUES 
	`

	for jmeno, v := range druhyTextu {
		prepare += fmt.Sprintf(`(%d, '%s', '%s', %d), `, v.ID, jmeno, v.Kategorie, v.Obtiznost)
	}
	prepare = prepare[:len(prepare)-2]
	prepare += ` ON CONFLICT DO NOTHING;`

	_, err := DB.Exec(prepare)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(dejepisQuery)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(zemepisQuery)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(pohadkyQuery)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(knihyQuery)
	if err != nil {
		panic(err)
	}
}

func pushDejepis() string {
	fmt.Println("Jdem na dějepis")
	fmt.Println("-----------------------------")

	f, err := os.Open("./texty/dejepis.txt")
	if err != nil {
		panic("spatna cesta k souboru")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var texty map[string]string = make(map[string]string)
	var jmeno string
	for scanner.Scan() {
		if jmeno == "" {
			jmeno = strings.TrimSpace(scanner.Text())
		} else if scanner.Text() == "" {
			texty[jmeno] = strings.TrimSpace(texty[jmeno])
			jmeno = ""
		} else {
			texty[jmeno] += strings.TrimSpace(scanner.Text())
			texty[jmeno] += " "
		}
	}
	texty[jmeno] = strings.TrimSpace(texty[jmeno])

	var pocet int = 0
	var st string = `INSERT INTO texty (cislo, jmeno, typ, txt, delka) VALUES `

	for k, v := range texty {
		delka := utf8.RuneCountInString(v)
		if delka < minDelkaTextu {
			continue
		}
		pocet++
		st += fmt.Sprintf(`(%d, '%s', %d, '%s', %d), `, pocet, k, 2, utils.UpravaTextu(v), delka)
	}

	st = st[:len(st)-2]
	st += ";"

	fmt.Printf("%v dějepisných textů jde do DB", pocet)

	return st
}

func pushZemepis() string {
	fmt.Println("Jdem na zeměpis")
	fmt.Println("-----------------------------")

	f, err := os.Open("./texty/zemepis.json")
	if err != nil {
		panic("spatna cesta k souboru")
	}
	defer f.Close()

	souborByte, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	type Zeme struct {
		Jmeno string `json:"name"`
		Obsah string `json:"content"`
	}
	var vsechnyZeme map[string]Zeme

	err = json.Unmarshal(souborByte, &vsechnyZeme)
	if err != nil {
		panic(err)
	}

	var st string = `INSERT INTO texty (cislo, jmeno, typ, txt, delka) VALUES `
	var pocet int
	for _, v := range vsechnyZeme {

		delka := utf8.RuneCountInString(v.Obsah)
		if delka < minDelkaTextu {
			continue
		}
		pocet++
		st += fmt.Sprintf(`(%d, '%s', %d, '%s', %d), `, pocet, strings.TrimSpace(v.Jmeno), 1, strings.TrimSpace(utils.UpravaTextu(v.Obsah)), utf8.RuneCountInString(v.Obsah))
	}

	st = st[:len(st)-2]
	st += ";"

	fmt.Printf("%v zeměpisných textů jde do DB", pocet)

	return st
}

func pushPohadky() string {
	fmt.Println("Jdem na Pohádky")
	fmt.Println("-----------------------------")

	var soubory []string = []string{"pohadkozem.txt", "pohadky_org.txt"}
	var texty map[string]string = make(map[string]string)

	for _, v := range soubory {
		f, err := os.Open("./texty/" + v)
		if err != nil {
			panic("spatna cesta k souboru")
		}
		scanner := bufio.NewScanner(f)

		var jmeno string
		for scanner.Scan() {
			if jmeno == "" {
				jmeno = strings.TrimSpace(scanner.Text())
			} else if scanner.Text() == "" {
				texty[jmeno] = strings.TrimSpace(texty[jmeno])
				jmeno = ""
			} else {
				texty[jmeno] += strings.TrimSpace(scanner.Text())
				texty[jmeno] += " "
			}
		}
		texty[jmeno] += strings.TrimSpace(texty[jmeno])

		f.Close()
	}

	var st string = `INSERT INTO texty (cislo, jmeno, typ, txt, delka) VALUES `
	r := regexp.MustCompile(`["\(\)]`)

	var easy, hard = 0, 0

	for k, v := range texty {
		delka := utf8.RuneCountInString(v)

		if delka < minDelkaTextu {
			continue
		}

		kolikSus := len(r.FindAllString(v, -1))

		if kolikSus != 0 {
			hard++
			st += fmt.Sprintf(`(%d, '%s', %d, '%s', %d), `, hard, k, 3, v, delka)
		} else {
			easy++
			st += fmt.Sprintf(`(%d, '%s', %d, '%s', %d), `, easy, k, 4, v, delka)
		}
	}

	fmt.Println(easy, hard)

	st = st[:len(st)-2]
	st += ";"

	fmt.Printf("%v pohádek jde do DB", easy+hard)

	return st
}
