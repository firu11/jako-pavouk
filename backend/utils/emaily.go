package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/badoux/checkmail"
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
		log.Print("NEFUNGUJE MAIL GG WOOHOO", err)
		MobilNotifikace("NEFUNGUJE MAIL")
		return err
	}
	log.Println("Posláno -", email)
	return nil
}

func ValidaceEmailu(email string) error {
	err := checkmail.ValidateHost(email)
	if err != nil {
		return err
	}
	err = checkmail.ValidateHostAndUser(os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_FROM"), email)
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		return errors.New(smtpErr.Error())
	}
	return nil
}
