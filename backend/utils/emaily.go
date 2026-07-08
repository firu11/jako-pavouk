package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/firu11/jako-pavouk/backend/config"
	"gopkg.in/gomail.v2"
)

type dataDosadit struct {
	Kod string
}

var (
	htmlEmail   *template.Template
	dialer      *gomail.Dialer
	emailConfig config.EmailConfig
)

func SetupEmaily(cfg config.EmailConfig) {
	emailConfig = cfg
	if !cfg.Enabled() {
		log.Println("email is not configured; email sending disabled")
		return
	}

	var err error
	htmlEmail, err = template.ParseFiles("./overovaci_email.html")
	if err != nil {
		log.Fatalf("failed to load email template: %v", err)
	}

	dialer = gomail.NewDialer(cfg.Host, cfg.Port, cfg.From, cfg.Password)
	dialer.SSL = true
}

func PoslatOverovaciEmail(email string, kod string) error {
	if htmlEmail == nil || dialer == nil {
		return errors.New("email is not configured")
	}

	buf := new(bytes.Buffer)
	if err := htmlEmail.Execute(buf, dataDosadit{Kod: kod}); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%s>", emailConfig.From))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikace")
	m.AddAlternative("text/plain", fmt.Sprintf("Tvůj ověřovací kód je: %s", kod))
	m.SetBody("text/html", buf.String())
	m.Embed("./pavoucekDoEmailu.png")

	if err := sendVicPokusu(m); err != nil {
		log.Println("NEFUNGUJE MAIL GG WOOHOO: ", email, err)
		notifyMobile("NEFUNGUJE MAIL " + err.Error())
		return err
	}

	log.Println("Posláno -", email)
	return nil
}

func PoslatInterniEmail(jmenoSkoly string, kontaktniEmail string, kontaktniTelefon string) error {
	if emailConfig.NotificationURL == "abc" || dialer == nil {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Jako Pavouk <%s>", emailConfig.From))
	m.SetHeader("To", emailConfig.To)
	m.SetHeader("Subject", "Nová škola")
	m.SetBody("text/plain", fmt.Sprintf("Někdo se zapsal se školou! \n\n %s\n%s\n%s", jmenoSkoly, kontaktniEmail, kontaktniTelefon))

	if err := sendVicPokusu(m); err != nil {
		log.Println("NEFUNGUJE NOVÁ ŠKOLA MAIL GG WOOHOO: ", err)
		notifyMobile("NEFUNGUJE MAIL " + err.Error())
		return err
	}

	return nil
}

func sendVicPokusu(m *gomail.Message) error {
	if dialer == nil {
		return errors.New("email is not configured")
	}

	var err error
	for i := range 3 {
		err = dialer.DialAndSend(m)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("pokus %d: EOF u posílání mailu\n", i)
				time.Sleep(time.Second)
				continue
			}
			return err
		}
		return nil
	}

	return err
}

func notifyMobile(message string) {
	if emailConfig.NotificationURL == "" {
		return
	}

	if _, err := http.Post(emailConfig.NotificationURL, "text/plain", strings.NewReader(message)); err != nil {
		log.Println("failed to send mobile notification:", err)
	}
}
