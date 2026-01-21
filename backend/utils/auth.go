package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	privatniKlic  []byte = []byte(os.Getenv("KLIC"))
	TokenLifetime time.Duration
	RefreshWindow = 24 * time.Hour
)

// obsah tokenu
type Data struct {
	jwt.RegisteredClaims
	Email string
	Id    uint
}

// hashování hesla
func HashPassword(heslo string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(heslo), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// porovnání dvou hesel
func CheckPassword(hesloRequest string, hesloDB string) error {
	if hesloDB == "google" {
		return errors.New("ucet je pres google")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hesloDB), []byte(hesloRequest))
	if err != nil {
		return err
	}
	return nil
}

// generace tokenu obsahujícího id, email a doba platnosti
func GenerovatToken(email string, id uint) (string, error) {
	data := Data{
		Email: email,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jakopavouk.cz",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	s, err := token.SignedString(privatniKlic)
	return s, err
}

// validace celého tokenu
func ValidovatToken(tokenString string) (bool, uint, error) {
	data := Data{}
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return privatniKlic, nil
	},
		jwt.WithLeeway(10*time.Second),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("jakopavouk.cz"),
	)
	if err != nil {
		return false, 0, err
	}
	return token.Valid, data.Id, err
}

// validace pouze toho zda je potřeba ho vyměnit
func ValidovatExpTokenu(tokenString string) (bool, error) {
	data := Data{}
	_, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (any, error) {
		return privatniKlic, nil
	},
		jwt.WithLeeway(10*time.Second),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("jakopavouk.cz"),
	)
	if err != nil {
		return false, err
	}

	if data.ExpiresAt == nil {
		return true, nil // consider expired
	}
	// pokud token vyprší během příštích 24 hodin, hodilo by se ho vyměnit už teď, abychom předešli tomu že vyprší ve špatnou chvíli
	timeLeft := time.Until(data.ExpiresAt.Time)
	return timeLeft <= RefreshWindow, nil
}
