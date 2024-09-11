package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"

	"gopkg.in/gomail.v2"
)

type fill struct {
	Kod string
}

// pošle email s kódem přes emailového poskytovatele
func PoslatOverovaciEmail(email string, kod string) error {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Panic("konverze portu na int se rozbila")
	}

	t, err := template.ParseFiles("./overovaci_email.html")
	if err != nil {
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
	}
	data := fill{
		Kod: kod,
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%v>", os.Getenv("EMAIL_FROM")))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikace")
	m.AddAlternative("text/plain", fmt.Sprintf("Tvůj ověřovací kód je: %s", kod))
	m.SetBody("text/html", buf.String())
	m.Embed("./pavoucekDoEmailu.png")

	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_FROM"), os.Getenv("EMAIL_HESLO"))
	if err := d.DialAndSend(m); err != nil {
		log.Print("NEFUNGUJE MAIL GG WOOHOO", err)
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
		return err
	}
	log.Println("Posláno -", email)
	return nil
}
