package utils

import (
	"backend/databaze"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var PocetPismenVeSlovu int

func GenerovatTextCviceni(pismena string, typ string, uzivID uint, tridaID uint, pocetZnaku float32) ([]string, error) {
	var text []string
	var err error

	switch typ {
	case "nova":
		if pismena == "zbylá diakritika" {
			pismena = "óďťň"
		}
		var pismenaRuny []rune = []rune(pismena)
		var pocetPismenek int = len(pismenaRuny)

		for _, v := range pismenaRuny {
			zakladni := PismenaPobliz[v]
			if zakladni != 0 {
				pismenaRuny = append(pismenaRuny, zakladni)
			}
		}

		var slovo strings.Builder
		for i := 0; i < int(pocetZnaku/float32(PocetPismenVeSlovu+1)); i++ {
			for j := 0; j < PocetPismenVeSlovu/2; j++ {
				r := rand.Intn(len(pismenaRuny))
				slovo.WriteRune(pismenaRuny[r])

				if len(pismenaRuny) != pocetPismenek*2 {
					r = rand.Intn(len(pismenaRuny))
					slovo.WriteRune(pismenaRuny[r])
					continue
				}

				if r+1 > pocetPismenek {
					r %= pocetPismenek
				} else {
					r += pocetPismenek
				}

				slovo.WriteRune(pismenaRuny[r])
			}
			slovo.WriteRune(' ')
			text = append(text, slovo.String())
			slovo.Reset()
		}
	case "naucena":
		var naucenaPismena string
		if pismena == "velká písmena (shift)" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťóFJGHDKSLAŮTZRUEIWOQPÚVBCNXMYĚŠČŘŽÝÁÍÉŇĎŤÓ"
		} else if pismena == "zbylá diakritika" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťó"
		} else if pismena == "čísla" {
			naucenaPismena = "fjghdkslaůtzrueiwoqpúvbcnxmyěščřžýáíéňďťó1234567890"
		} else {
			naucenaPismena, err = databaze.GetNaucenaPismena(uzivID, pismena, tridaID)
			if err != nil {
				log.Println(err)
				return text, err
			}
		}

		var slovo strings.Builder
		for i := 0; i < int(pocetZnaku/float32(PocetPismenVeSlovu+1)); i++ {
			for j := 0; j < PocetPismenVeSlovu; j++ {
				r := rand.Intn(utf8.RuneCountInString(naucenaPismena)) // utf-8 jsou sus
				slovo.WriteRune([]rune(naucenaPismena)[r])
			}
			slovo.WriteRune(' ')
			text = append(text, slovo.String())
			slovo.Reset()
		}
	case "slova":
		var slova []string
		slova, err = databaze.GetSlovaProLekci(uzivID, pismena, int(pocetZnaku/7.5+50), tridaID)
		if err != nil {
			log.Println(err)
			return text, err
		}

		var pocetSlovKMani int = len(slova)
		var i int = 0
		var delka int = 0

		for pocetSlovKMani != 0 {
			if float32(delka) >= pocetZnaku {
				break
			}
			text = append(text, slova[i]+" ")
			delka += utf8.RuneCountInString(slova[i]) + 1

			i++
			if i >= pocetSlovKMani {
				i = 0
			}
		}

		if pismena == "velká písmena (shift)" { // dam kazdy prvni pismeno velkym
			for i := 0; i < len(text); i++ {
				r := []rune(text[i])
				text[i] = fmt.Sprintf("%c%s", unicode.ToUpper(r[0]), string(r[1:]))
			}
		} else if pismena == "čísla" {
			for i := 1; i < len(text); i += 2 {
				r := rand.Intn(8999) + 1000
				text[i] = strconv.Itoa(r) + " "
			}
		} else if pismena == "interpunkce" {
			znamenka := []string{"?", "!", ",", "."}
			delka := len(znamenka)
			for i := 0; i < len(text); i++ {
				r := rand.Intn(delka)
				text[i] = strings.Replace(text[i], " ", znamenka[r]+" ", 1)
			}
		}
	case "programator":
		var slova []string
		slova, err = databaze.GetProgramatorSlova()
		if err != nil {
			log.Println(err)
			return text, err
		}
		var pocetSlovKMani int = len(slova)

		if pismena == "závorky" {
			var zavorky []string = []string{"[\u005D", "()", "{}", "<>"}
			var zavorkyLen = len(zavorky)
			var i int = 0
			var zi int = 0
			var delka int = 0

			rand.Shuffle(zavorkyLen, func(i, j int) { zavorky[i], zavorky[j] = zavorky[j], zavorky[i] })
			for {
				if float32(delka) >= pocetZnaku {
					break
				}
				str := fmt.Sprintf("%s%v%s ", string([]rune(zavorky[zi])[0]), slova[i], string([]rune(zavorky[zi])[1]))
				text = append(text, str)
				delka += utf8.RuneCountInString(str)

				i++
				zi++
				if i >= pocetSlovKMani {
					i = 0
				}
				if zi >= zavorkyLen {
					zi = 0
				}
			}
		} else if pismena == "operátory" {
			var oper []string = []string{"=", "==", "!=", "<=", ">=", "<", ">", "+", "-", "*", "/", "%", "+=", "-=", "*=", "/=", "||", "&&"}
			var operLen = len(oper)
			var i int = 0
			var zi int = 0
			var delka int = 0

			rand.Shuffle(operLen, func(i, j int) { oper[i], oper[j] = oper[j], oper[i] })
			text = append(text, slova[pocetSlovKMani-1]+" ")
			for {
				if float32(delka) >= pocetZnaku {
					break
				}
				str := fmt.Sprintf("%s %v ", string([]rune(oper[zi])), slova[i])
				text = append(text, str)
				delka += utf8.RuneCountInString(str)

				i++
				zi++
				if i >= pocetSlovKMani {
					i = 0
				}
				if zi >= operLen {
					zi = 0
				}
			}
		}
	default:
		log.Print("Cviceni ma divnej typ")
		return text, errors.New("cviceni ma divnej typ")
	}

	return text, nil
}
