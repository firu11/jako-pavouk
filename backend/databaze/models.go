package databaze

import (
	"time"

	"github.com/rickb777/date"
)

type (
	Lekce struct {
		ID      uint   `json:"id" db:"id"`
		Pismena string `json:"pismena" db:"pismena"`
		// Skupina uint        nepouzivame ale je tam
		// Klavesnice string
	}

	Cviceni struct {
		ID  uint   `json:"id" db:"id"`
		Typ string `json:"typ" db:"typ"`
	}

	Uzivatel struct {
		ID              uint      `json:"id" db:"id"`
		Email           string    `json:"email" db:"email"`
		Jmeno           string    `json:"jmeno" db:"jmeno"`
		Heslo           string    `json:"heslo" db:"heslo"`
		Klavesnice      string    `json:"klavesnice" db:"klavesnice"`
		Datum           date.Date `json:"datum" db:"datum"`
		Role            int       `json:"role" db:"role"`
		UcitelVeSkoleID int       `json:"" db:"skola_id"`
		SkolniJmeno     string    `json:"skolni_jmeno" db:"skolni_jmeno"`
	}

	NeoUziv struct {
		Email  string `json:"email" db:"email"`
		Jmeno  string `json:"jmeno" db:"jmeno"`
		Heslo  string `json:"heslo" db:"heslo"`
		Kod    string `json:"kod" db:"kod"`
		Cas    int64  `json:"cas" db:"cas"`
		Pokusy int    `json:"pokusy" db:"pokusy"`
	}

	ZmenaHeslaUziv struct {
		Email string `json:"email" db:"email"`
		Kod   string `json:"kod" db:"kod"`
		Cas   int64  `json:"cas" db:"cas"`
	}

	Slovnik struct {
		ID    uint   `json:"id" db:"id"`
		Slovo string `json:"slovo" db:"slovo"`
	}

	Dokoncene struct {
		ID            uint           `json:"id" db:"id"`
		UzivID        uint           `json:"uziv_id" db:"uziv_id"`
		CviceniID     uint           `json:"cviceni_id" db:"cviceni_id"`
		Neopravene    uint           `json:"neopravene" db:"neopravene"`
		Cas           float32        `json:"cas" db:"cas"`
		DelkaTextu    int            `json:"delka_textu" db:"delka_textu"`
		Datum         date.Date      `json:"datum" db:"datum"`
		ChybyPismenka map[string]int `json:"chyby_pismenka" db:"chyby_pismenka"`
	}

	Trida struct {
		ID         uint   `json:"id" db:"id"`
		Jmeno      string `json:"jmeno" db:"jmeno"`
		UcitelID   uint   `json:"-" db:"ucitel_id"`
		Kod        string `json:"kod" db:"kod"`
		Zamknuta   bool   `json:"zamknuta" db:"zamknuta"`
		Smazana    bool   `json:"smazana" db:"smazana"`
		Klavesnice string `json:"klavesnice" db:"klavesnice"`
	}

	Skola struct {
		ID          uint      `json:"id" db:"id"`
		Jmeno       string    `json:"jmeno" db:"jmeno"`
		DenZalozeni date.Date `json:"den_zalozeni" db:"den_zalozeni"`
		Aktivni     bool      `json:"aktivni" db:"aktivni"`
		// kontaktni_email
		// kontaktni_telefon
	}

	Ucitel struct {
		ID      uint `json:"id" db:"id"`
		SkolaID uint `json:"skola_id" db:"skola_id"`

		UzivID uint   `json:"uziv_id" db:"uziv_id"`
		Email  string `json:"email" db:"email"`

		PocetTrid int `json:"pocet_trid" db:"pocet_trid"`
		PocetZaku int `json:"pocet_zaku" db:"pocet_zaku"`
	}

	Prace struct {
		ID                uint      `json:"id" db:"id"`
		TridaID           uint      `json:"-" db:"trida_id"`
		Text              string    `json:"text" db:"text"`
		Cas               int       `json:"cas" db:"cas"`
		Datum             time.Time `json:"datum" db:"datum"`
		PrumerneCPM       float32   `json:"prumerne_cpm" db:"prumerne_cpm"`
		PrumernaPresnost  float32   `json:"prumerna_presnost" db:"prumerna_presnost"`
		StudentuDokoncilo int       `json:"studentu_dokoncilo" db:"studentu_dokoncilo"`
	}

	Procvic struct {
		ID        uint    `json:"id" db:"id"`
		Jmeno     string  `json:"jmeno" db:"jmeno"`
		Kategorie string  `json:"kategorie" db:"kategorie"`
		CPM       float32 `json:"cpm" db:"-"`
		Obtiznost int     `json:"obtiznost" db:"obtiznost"`
	}
)

// helpers
type (
	Cvic struct {
		Id  int     `json:"id"`
		Cpm float32 `json:"cpm"`
	}

	TridaInfo struct {
		ID            uint   `json:"id" db:"id"`
		Jmeno         string `json:"jmeno" db:"jmeno"`
		Kod           string `json:"kod" db:"kod"`
		Zamknuta      bool   `json:"zamknuta" db:"zamknuta"`
		PocetStudentu int    `json:"pocet_studentu" db:"pocet_studentu"`
		PocetPraci    int    `json:"pocet_praci" db:"pocet_praci"`
	}

	Student struct {
		ID       uint    `json:"id" db:"id"`
		Jmeno    string  `json:"jmeno" db:"skolni_jmeno"`
		Email    string  `json:"email" db:"email"`
		CPM      float64 `json:"cpm" db:"cpm"`
		Presnost float64 `json:"presnost" db:"presnost"`
	}

	Cviceni2 struct {
		Id      uint   `json:"id"`
		LekceID uint   `json:"lekce_id"`
		Pismena string `json:"pismena"`
	}
)
