package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type dataDosadit struct {
	Kod string
}

var htmlEmail *template.Template
var dialer *gomail.Dialer

func SetupEmaily() error {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Panic("konverze portu na int se rozbila")
	}

	htmlEmail, err = template.ParseFiles("./overovaci_email.html")
	if err != nil {
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
	}

	dialer = gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_FROM"), os.Getenv("EMAIL_HESLO"))

	return nil
}

func PoslatOverovaciEmail(email string, kod string) error {
	data := dataDosadit{
		Kod: kod,
	}
	buf := new(bytes.Buffer)
	if err := htmlEmail.Execute(buf, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%v>", os.Getenv("EMAIL_FROM")))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikace")
	m.AddAlternative("text/plain", fmt.Sprintf("Tvůj ověřovací kód je: %s", kod))
	m.SetBody("text/html", buf.String())
	m.Embed("./pavoucekDoEmailu.png")

	if err := dialer.DialAndSend(m); err != nil {
		log.Print("NEFUNGUJE MAIL GG WOOHOO", err)
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
		return err
	}
	log.Println("Posláno -", email)
	return nil
}

func PoslatInterniEmail(jmenoSkoly string, kontaktniEmail string, kontaktniTelefon string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%v>", os.Getenv("EMAIL_FROM")))
	m.SetHeader("To", os.Getenv("EMAIL_MUJ"))
	m.SetHeader("Subject", "Nová škola")
	m.SetBody("text/plain", fmt.Sprintf("Někdo se chce zapsat se školou! \n\n %s\n%s\n%s", jmenoSkoly, kontaktniEmail, kontaktniTelefon))

	if err := dialer.DialAndSend(m); err != nil {
		log.Print("NEFUNGUJE MAIL GG WOOHOO", err)
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
		return err
	}
	return nil
}
