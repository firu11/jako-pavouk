package databaze

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/blockloop/scan/v2"
)

func CreateTrida(jmeno string, ucitelID uint, kod string) error {
	_, err := DB.Exec(`INSERT INTO trida (jmeno, ucitel_id, kod) VALUES ($1, (SELECT id FROM ucitel WHERE uziv_id = $2), $3)`, jmeno, ucitelID, kod)
	return err
}

func GetTrida(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.id, t.jmeno, u.uziv_id AS ucitel_id, t.kod, t.zamknuta, t.smazana, t.klavesnice FROM trida t INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE t.id = $1 AND NOT t.smazana AND NOT u.smazany;`, id)
	if err != nil {
		return trida, err
	}
	defer row.Close()
	err = scan.Row(&trida, row)
	return trida, err
}

func GetTridaByStudentID(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.* FROM trida t INNER JOIN student_a_trida sat ON t.id = sat.trida_id INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE NOT t.smazana AND NOT u.smazany AND sat.student_id = $1 LIMIT 1;`, id)
	if err != nil {
		return trida, err
	}
	defer row.Close()
	err = scan.Row(&trida, row)
	return trida, err
}

func GetTridy(ucitelID uint) ([]TridaInfo, error) {
	var tridy []TridaInfo = []TridaInfo{}

	rows, err := DB.Query(`SELECT t.id, t.jmeno, t.kod, t.zamknuta, ( SELECT COUNT(*) FROM uzivatel uz INNER JOIN student_a_trida s ON s.student_id = uz.id WHERE s.trida_id = t.id AND NOT uz.smazany ) as pocet_studentu, ( SELECT COUNT(*) FROM prace p WHERE p.trida_id = t.id AND NOT p.smazana ) as pocet_praci FROM trida t INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE NOT t.smazana AND NOT u.smazany AND u.uziv_id = $1;`, ucitelID)
	if err != nil {
		return tridy, err
	}
	defer rows.Close()

	err = scan.Rows(&tridy, rows)
	return tridy, err
}

func GetStudentyZeTridy(tridaID uint) ([]Student, error) {
	var zaci []Student = []Student{}

	rows, err := DB.Query(`WITH cpm_data AS ( SELECT datum, delka_textu, neopravene, cas, uziv_id FROM dokoncene WHERE datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) UNION ALL SELECT datum, delka_textu, neopravene, cas, uziv_id FROM dokoncene_procvic WHERE datum::date > CURRENT_DATE - MAKE_INTERVAL(days => 13) ), cpm AS ( SELECT uziv_id, GREATEST( ( (SUM(delka_textu) - 10 * SUM(neopravene)) / SUM(cas)::NUMERIC ) * 60, 0 ) AS cpm FROM cpm_data GROUP BY uziv_id ) SELECT u.id, u.skolni_jmeno, u.email, COALESCE(mc.cpm, 0) as cpm FROM uzivatel u INNER JOIN student_a_trida s ON s.student_id = u.id INNER JOIN trida t ON t.id = s.trida_id FULL OUTER JOIN cpm mc ON mc.uziv_id = u.id WHERE s.trida_id = $1 AND t.smazana = FALSE AND NOT u.smazany GROUP BY u.id, u.skolni_jmeno, u.email, mc.cpm;`, tridaID)
	if err != nil {
		return zaci, err
	}
	defer rows.Close()

	err = scan.Rows(&zaci, rows)
	return zaci, err
}

func GetStudentyZPrace(praceID uint) ([]Student, error) {
	var zaci []Student = []Student{}

	rows, err := DB.Query(`WITH studenti_ze_tridy AS ( SELECT u.id, u.skolni_jmeno, u.email FROM student_a_trida sat JOIN uzivatel u ON u.id = sat.student_id WHERE sat.trida_id = ( SELECT trida_id FROM prace WHERE id = $1 ) AND NOT u.smazany ), vysledky_teto_prace AS ( SELECT student_id as id, GREATEST( ( (delka_textu - 10 * neopravene) / cas::NUMERIC ) * 60, 0 ) AS cpm, COALESCE( ( ( delka_textu - neopravene - COALESCE( ( SELECT SUM(value::NUMERIC) FROM jsonb_each_text(chyby_pismenka) ), 0 ) ) / GREATEST(delka_textu::NUMERIC, 1) ) * 100, -1 ) AS presnost FROM dokoncena_prace WHERE prace_id = $1 ) SELECT id, skolni_jmeno, email, COALESCE(cpm, -1) as cpm, COALESCE(presnost, -1) as presnost FROM studenti_ze_tridy LEFT JOIN vysledky_teto_prace USING(id);`, praceID)
	if err != nil {
		return zaci, err
	}
	defer rows.Close()

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
	defer row.Close()
	err = scan.Row(&trida, row)
	return trida, err
}

func GetTridaByUziv(id uint) (Trida, error) {
	var trida Trida
	row, err := DB.Query(`SELECT t.* FROM trida t INNER JOIN student_a_trida s ON t.id = s.trida_id WHERE student_id = $1 AND NOT smazana;`, id)
	if err != nil {
		return trida, err
	}
	defer row.Close()
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

func CreateSkola(jmeno, email, telefon string) (uint, error) {
	var skolaID uint
	err := DB.QueryRow(`INSERT INTO skola (jmeno, kontaktni_email, kontaktni_telefon) VALUES ($1, $2, $3) RETURNING id;`, jmeno, email, telefon).Scan(&skolaID)
	return skolaID, err
}

func GetSkolaByUcitel(uzivID uint) (Skola, error) {
	var skola Skola
	rows, err := DB.Query(`SELECT s.* FROM skola s INNER JOIN ucitel u ON u.skola_id = s.id INNER JOIN uzivatel uz ON uz.id = u.uziv_id WHERE uz.id = $1 AND NOT u.smazany AND NOT uz.smazany;`, uzivID)
	if err != nil {
		return skola, err
	}
	defer rows.Close()

	err = scan.Row(&skola, rows)
	return skola, err
}

func GetUcitele(skolaID uint) ([]Ucitel, error) {
	var ucitele []Ucitel
	rows, err := DB.Query(`WITH studenti AS ( SELECT * FROM student_a_trida sat INNER JOIN uzivatel uz ON sat.student_id = uz.id WHERE NOT uz.smazany ), studenti_ucitele AS ( SELECT u.id, COUNT(s.student_id) AS pocet_studentu FROM trida t LEFT JOIN studenti s ON s.trida_id = t.id INNER JOIN ucitel u ON u.id = t.ucitel_id WHERE NOT t.smazana GROUP BY u.id ), staty AS ( SELECT ucitel_id, COUNT(*) AS pocet_trid FROM trida WHERE NOT smazana GROUP BY ucitel_id ) SELECT uz.email, COALESCE(pocet_trid, 0) AS pocet_trid, u.id, uz.id AS uziv_id, COALESCE(SUM(stu.pocet_studentu), 0) AS pocet_zaku FROM ucitel u INNER JOIN uzivatel uz ON uz.id = u.uziv_id LEFT JOIN staty s ON s.ucitel_id = u.id LEFT JOIN studenti_ucitele stu ON stu.id = u.id WHERE u.skola_id = $1 AND NOT u.smazany GROUP BY uz.email, pocet_trid, u.id, uz.id ORDER BY u.id;`, skolaID)
	if err != nil {
		return ucitele, err
	}
	defer rows.Close()

	err = scan.Rows(&ucitele, rows)
	return ucitele, err
}

func CreateUcitel(skolaID, uzivID uint) error {
	_, err := DB.Exec(`INSERT INTO ucitel (uziv_id, skola_id) VALUES ($1, $2) ON CONFLICT (uziv_id) DO UPDATE SET smazany = false, skola_id = EXCLUDED.skola_id;`, uzivID, skolaID)
	return err
}

func RemoveUcitelByEmail(email string) error {
	_, err := DB.Exec(`UPDATE ucitel SET smazany = true WHERE uziv_id = ( SELECT id FROM uzivatel WHERE email = $1 );`, email)
	return err
}
