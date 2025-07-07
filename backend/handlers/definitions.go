package handlers

import (
	"regexp"
	"time"
)

var PocetZnaku float32
var RegexJmeno *regexp.Regexp

// struct body requestu
type (
	bodyDokoncit struct {
		Preklepy         int            `json:"neopravenePreklepy" validate:"min=0"` //sus reqired nebere nulu takze min=0 asi ok
		Cas              int            `json:"cas" validate:"required"`
		DelkaTextu       int            `json:"delkaTextu" validate:"required"`
		NejcastejsiChyby map[string]int `json:"nejcastejsiChyby" validate:"required"`
	}

	bodyPoslatEmail struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
	}

	bodyRegistrace struct {
		Email string `json:"email" validate:"required,email"`
		Jmeno string `json:"jmeno" validate:"required,min=3,max=16"`
		Heslo string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyPrihlaseni struct {
		EmailNeboJmeno string `json:"email" validate:"required"`
		Heslo          string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyUprava struct {
		Zmena   string `json:"zmena"`
		Hodnota string `json:"hodnota"`
	}

	bodyZmenaHesla struct {
		Email string `json:"email" validate:"required,email"`
	}

	bodyOvereniZmenaHesla struct {
		Email string `json:"email" validate:"required,email"`
		Kod   string `json:"kod" validate:"required,len=5"`
		Heslo string `json:"heslo" validate:"required,ascii,min=5,max=72"`
	}

	bodyGoogle struct {
		AccessToken string `json:"access_token"`
	}

	bodyTestPsani struct {
		Typ string `json:"typ" validate:"required"`
	}
)

// skolni
type (
	bodyCreateTrida struct {
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
	bodyZmenaTridy struct {
		TridaID uint   `json:"trida_id" validate:"required"`
		Zmena   string `json:"zmena"`
		Hodnota string `json:"hodnota"`
	}
	bodyZmenaJmenaStudenta struct {
		ID      uint   `json:"id" validate:"required"`
		Jmeno   string `json:"jmeno" validate:"max=30"`
		TridaID int    `json:"trida_id"`
	}
	bodyZapis struct {
		Kod   string `json:"kod" validate:"required"`
		Jmeno string `json:"jmeno" validate:"required,min=1,max=30"`
	}
	bodyPridatPraci struct {
		Text    string `json:"text" validate:"required,min=1"`
		Cas     int    `json:"cas" validate:"required"`
		TridaID uint   `json:"trida_id" validate:"required"`
	}
	bodyGetText struct {
		TridaID uint    `json:"trida_id"`
		Typ     string  `json:"typ" validate:"required"`
		ZLekce  string  `json:"z_lekce"`
		Delka   float32 `json:"delka"`
	}
	bodyZapisSkoly struct {
		JmenoSkoly       string `json:"jmeno_skoly"`
		KontaktniEmail   string `json:"kontaktni_email" validate:"email"`
		KontaktniTelefon string `json:"kontaktni_telefon"`
	}
	bodyUpravaUcitele struct {
		Akce  string `json:"akce" validate:"oneof='smazat' 'pridat'"`
		Email string `json:"email" validate:"required_if=Akce 'pridat',email"`
	}

	praceProStudenta struct {
		ID       uint      `json:"id"`
		TridaID  uint      `json:"-"`
		Text     string    `json:"-"`
		Cas      int       `json:"cas"`
		Datum    time.Time `json:"datum"`
		Cpm      float64   `json:"cpm"`
		Presnost float64   `json:"presnost"`
	}
)
