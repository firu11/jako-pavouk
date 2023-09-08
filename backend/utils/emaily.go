package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func PoslatOverovaciEmail(email string, kod string) error {
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_UZIV"), os.Getenv("EMAIL_HESLO"), os.Getenv("EMAIL_HOST"))
	msg := []byte(fmt.Sprintf("To: %v\r\n"+
		"From: Jako Pavouk <%v>\r\n"+
		"Subject: Verifikace\r\n"+
		"\r\n"+
		"Váš ověřovací kód pro Jako Pavouk je: %v\r\n", email, os.Getenv("EMAIL_FROM"), kod))
	adresat := []string{email}

	err := smtp.SendMail(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"), auth, os.Getenv("EMAIL_FROM"), adresat, msg)
	if err != nil {
		return err
	}

	fmt.Println("Posláno -", email)
	return nil
}
