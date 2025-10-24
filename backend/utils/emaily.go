package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type dataDosadit struct {
	Kod string
}

var (
	htmlEmail *template.Template
	dialer    *gomail.Dialer
)

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
	dialer.SSL = true

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

	if err := sendVicPokusu(m); err != nil {
		log.Println("NEFUNGUJE MAIL GG WOOHOO: ", email, err)
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
	m.SetBody("text/plain", fmt.Sprintf("Někdo se zapsal se školou! \n\n %s\n%s\n%s", jmenoSkoly, kontaktniEmail, kontaktniTelefon))

	if err := sendVicPokusu(m); err != nil {
		log.Println("NEFUNGUJE NOVÁ ŠKOLA MAIL GG WOOHOO: ", err)
		MobilNotifikace("NEFUNGUJE MAIL " + err.Error())
		return err
	}
	return nil
}

func sendVicPokusu(m *gomail.Message) error {
	var err error
	for i := range 3 {
		err = dialer.DialAndSend(m)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("pokus %d: EOF u posílání mailu\n", i)
				time.Sleep(time.Second)
				continue // shit happens, try again
			}
			return err
		}
		return nil
	}
	return err
}
