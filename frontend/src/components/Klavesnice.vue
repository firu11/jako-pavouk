<template>
    <div id="klavesnice" class="pruhledne">
        <div class="radek" v-for="radek in schema">
            <div v-for="tlacitko in radek"
                 class="klavesa"
                 :class="{oznacenaKlavesa: oznacene(tlacitko) || (tlacitko === 'Shift' && this.shiftSviti)}"
                 :style="{backgroundColor: barva(tlacitko), flexGrow: delkaTlacitka(tlacitko)}">

                <div v-if="tlacitko !== '^v'" :style="{color: '#000'}">
                    {{ tlacPismeno(0, tlacitko) }} <br>
                    {{ tlacPismeno(1, tlacitko) }}
                </div>
                <div v-else id="sipky">
                    <div class="klavesa" style="height: 18px"></div>
                    <div class="klavesa" style="height: 18px"></div>
                </div>
            </div>

        </div>
    </div>
</template>

<script>
export default {
    name: "Klavesnice",
    props: ["aktivniPismeno"],
    data() {
        let barvy = ["#6ada56", "#81bffc", "#fa5ca1", "#ff8800", "#6f86f7"]
        return {
            schema: [
                ["°;", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⟵"],
                ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "'¨"],
                ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "Enter ↵"],
                ["Shift", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift"],
                ["  ", "", "", "", "______", "", "", "", "^v", ""]
            ], delkaKlaves: {"⟵": 3, "Shift": 1, "Enter ↵": 1, "CapsLock": 1, "TAB": 1, "______": 24, "  ": 2},
            prstoklad: {
                P_Ukaz: [barvy[0], "J", "H", "U", "Z", "N", "M", "7ý", "6ž"],
                L_Ukaz: [barvy[1], "G", "T", "R", "F", "V", "B", "5ř", "4č"],
                P_Pros: [barvy[2], "K", "I", "?,", "8á"],
                L_Pros: [barvy[2], "D", "E", "C", "3š"],
                P_Prs: [barvy[3], "O", "L", ":.", "9í"],
                L_Prs: [barvy[3], "X", "S", "W", "2ě"],
                P_Mali: [barvy[4], '"ů', "P", "_-", "0é", '%=', 'ˇ´', '⟵', '()', '/ú', "'¨", '!§', 'Enter ↵', 'Shift'],
                L_Mali: [barvy[4], "Shift", "A", "Q", "Y", "1+", "°;", "TAB", "CapsLock", "Ctrl"],
                Palce: ["#bc73ff", "______"]
            },
            shiftSviti: false
        }
    },
    methods: {
        tlacPismeno(cislo, tlacitko) {
            if (tlacitko.length === 2) return tlacitko.at(cislo)
            else if (tlacitko.length === 1 && cislo === 0) return tlacitko.at(0)
            else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
        },
        oznacene(tlacitko) {
            let pismeno = this.aktivniPismeno[0]
            if ((pismeno === ' ' && tlacitko === '______') || (tlacitko.length === 1 && tlacitko.toLowerCase() === pismeno.toLowerCase()) || (tlacitko.length === 2 && tlacitko.toLowerCase().includes(pismeno.toLowerCase()) && tlacitko !== '  ')) {
                return true
            } else {
                return false
            }
        },
        barva(tlacitko) {
            let pismeno = this.aktivniPismeno[0]
            if (tlacitko === 'Shift') {
                if (this.potrebujeShift(pismeno) && pismeno !== " ") {
                    this.shiftSviti = true
                } else this.shiftSviti = false
                return this.prstoklad['P_Mali'][0]

            } else if (tlacitko.length === 2 && tlacitko.toLowerCase().includes(pismeno.toLowerCase())) {
                for (let prst in this.prstoklad) {
                    for (let tla in this.prstoklad[prst]) {
                        if (this.prstoklad[prst][tla].toLowerCase().charAt(0) === tlacitko.toLowerCase().charAt(0)) {
                            return this.prstoklad[prst][0]
                        } else if (this.prstoklad[prst][tla].toLowerCase().charAt(1) === tlacitko.toLowerCase().charAt(1)) {
                            return this.prstoklad[prst][0]
                        }
                    }
                }
            } else {
                for (let prst in this.prstoklad) {
                    for (let tla in this.prstoklad[prst]) {
                        if (this.prstoklad[prst][tla].toLowerCase() === tlacitko.toLowerCase()) {
                            return this.prstoklad[prst][0]
                        }
                    }
                }
            }
        },
        delkaTlacitka(tlacitko) {
            if (this.delkaKlaves[tlacitko] === undefined) {
                return 0
            }
            return this.delkaKlaves[tlacitko]
        },
        potrebujeShift(pismeno) {
            if (['"', '/', '?', ':', '_', '!', '(', '%', 'ˇ', '°'].includes(pismeno)) {
                return true
            } else if (/^\d$/.test(pismeno)) {
                return true
            }

            return (pismeno === pismeno.toUpperCase() && pismeno !== 'Ů' && pismeno !== 'Ú')  // Ů a Ú nejdou psat se shiftem

        }
    }
}
</script>

<style scoped>
.klavesa {
    width: 40px;
    height: 40px;
    background-color: var(--fialova);
    border-radius: 5px;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: background-color 0.2s;
}

.oznacenaKlavesa {
    border: #fff solid 5px;
    filter: brightness(150%);
}

.oznacenaKlavesa div {
    font-weight: 700 !important;
}

.klavesa div {
    font-weight: 500;
}

.klavesa:has(#sipky) {
    background-color: var(--tmave-fialova);
}

#klavesnice {
    display: flex;
    flex-direction: column;
    gap: 4px;
    background-color: var(--tmave-fialova);
    padding: 10px;
    border-radius: 10px;
    font-size: 0.8em;
    line-height: 1.3em;
    width: 675px;
    margin-top: 50px;
}

.radek {
    display: flex;
    gap: 5px;
}

#sipky {
    display: flex;
    flex-direction: column;
    gap: 4px;
}
</style>