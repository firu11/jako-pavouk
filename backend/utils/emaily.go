package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// pošle email s kódem přes emailového poskytovatele
func PoslatOverovaciEmail(email string, kod string) error {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Panic("konverze portu na int se rozbila")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%v>", os.Getenv("EMAIL_FROM")))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikace")
	m.SetBody("text/html", fmt.Sprintf("Váš ověřovací kód na Jako Pavouk je: <b>%v</b>", kod))

	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_FROM"), os.Getenv("EMAIL_HESLO"))
	if err := d.DialAndSend(m); err != nil {
		log.Print("NEFUNGUJE MAIL GG WOOWOO", err)
		MobilNotifikace("NEFUNGUJE MAIL")
		return err
	}
	log.Println("Posláno -", email)
	return nil
}
