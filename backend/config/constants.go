package config

import (
	"math"
	"regexp"
	"time"
)

const (
	PocetZnaku         = 1500
	PocetPismenVeSlovu = 4
	TokenLifetime      = time.Hour * 24 * 15
	CifraCislaZaJmenem = 4
)

var (
	RegexJmeno      = regexp.MustCompile(`^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$`)
	MaxCisloZaJmeno = int(math.Pow(10, float64(CifraCislaZaJmenem)))
)
