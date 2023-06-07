package utils

import (
	"errors"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	return err
}

/*
Vraci id, error.
id = 0 znamena ze se bud neco pokazilo nebo je autentizace nepovinna
*/
func Autentizace(c *fiber.Ctx, povinna bool) (uint, error) {
	if c.Get("Token") != "" && len(c.Get("Token")) != 0 {
		spravnej, id, err := ValidovatToken(c.Get("Token"))

		if spravnej && err == nil {
			return id, nil
		} else if !spravnej {
			if povinna {
				return 0, errors.New("spatny token")
			} else {
				return 0, nil
			}
		} else {
			return 0, err
		}
	} else {
		if povinna {
			return 0, errors.New("je potreba autentizace (JWT Token)")
		}
		return 0, nil
	}
}

func Prumer(arr []float32) float32 {
	var soucet float32 = 0
	for _, v := range arr {
		soucet += v
	}
	if len(arr) == 0 {
		return 0
	}
	return soucet / float32(len(arr))
}
