package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"sync"
)

type UzivCekajici struct {
	Jmeno      string
	Email      string
	HesloHash  string
	Kod        string
	DobaTrvani int64
}

var UzivCekajiciNaOvereni []UzivCekajici

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

	fmt.Println("Email Sent Successfully!")
	return nil
}

func VycistitSeznam(deleteIndexy []int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for _, indexDel := range deleteIndexy {
			UzivCekajiciNaOvereni[indexDel] = UzivCekajiciNaOvereni[len(UzivCekajiciNaOvereni)-1]
			UzivCekajiciNaOvereni = UzivCekajiciNaOvereni[:len(UzivCekajiciNaOvereni)-1]
		}
	}()
}
