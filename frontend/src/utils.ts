import { ref } from "vue";
import { cislaProcvicJmeno, levelyPresnosti, levelyRychlosti, nastaveniJmeno, tokenJmeno } from "./stores";

export function formatovanyPismena(pismena: string | string[] | undefined) {
    if (pismena === "..." || pismena === undefined) return pismena
    let vratit = "";
    for (let i = 0; i < pismena.length; i++) {
        vratit += i < pismena.length - 1 ? pismena.at(i) + ", " : pismena.at(i);
    }
    return vratit;
}

export function format(p: string) {
    if (p === "zbylá diakritika") return "Zbylá diakritika"
    else if (p === "velká písmena (shift)") return "Velká písmena (Shift)"
    else if (p === "závorky") return "Závorky"
    else if (p === "operátory") return "Operátory"
    else if (p === "čísla") return "Číslovky"
    else if (p === "interpunkce") return "Interpunkce"
    return formatovanyPismena(p)
}

export function getToken() {
    return localStorage.getItem(tokenJmeno)
}

export function getCisloProcvic(id: string) {
    let cislo = localStorage.getItem(cislaProcvicJmeno + id)
    if (cislo === null) {
        localStorage.setItem(cislaProcvicJmeno + id, "2")
        return "1"
    }
    if (cislo == "10") {
        localStorage.setItem(cislaProcvicJmeno + id, "1")
        return cislo
    }
    localStorage.setItem(cislaProcvicJmeno + id, String(Number(cislo) + 1))
    return cislo
}

export const oznameni = ref([] as { text: string, typ: string }[])

export function pridatOznameni(text: string = "Něco se pokazilo", cas: number = 4000, typ: string = "vykricnik") {
    let obj = { text: text, typ: typ }
    oznameni.value.push(obj)
    setTimeout(() => {
        oznameni.value.splice(oznameni.value.indexOf(obj), 1);
    }, cas)
}

export function napovedaKNavigaci() {
    pridatOznameni("Pro nápovědu k navigaci se podívej do záložky Jak psát.")
}

export function checkTeapot(e: any) {
    if (e.response && e.response.status == 418) {
        if (oznameni.value.length < 3) {
            pridatOznameni("Dej si čajík a vydýchej se...")
        }
        return true
    }
    return false
}

export class Oznacene {
    index = ref(0)
    max: number = 4
    bezOznaceni: boolean = false
    mensi() {
        if (this.index.value > 1) {
            this.index.value--
        }
    }
    vetsi() {
        if (this.index.value < this.max) {
            this.index.value++
        }
    }
    setMax(max: number) {
        this.max = max
    }
    is(n: number) {
        if (n < 6 && n == this.index.value) return true
        else if (n >= 6 && 14 > n && this.index.value + 1 == n) return true
        else if (n >= 14 && this.index.value + 2 == n) return true
        return false
    }
}

// https://stackoverflow.com/questions/20084513/detect-search-crawlers-via-javascript
export function jeToRobot(userAgent: string): boolean {
    const robots = new RegExp(([
        /bot/, /spider/, /crawl/,                               // GENERAL TERMS
        /APIs-Google/, /AdsBot/, /Googlebot/,                   // GOOGLE ROBOTS
        /mediapartners/, /Google Favicon/,
        /FeedFetcher/, /Google-Read-Aloud/,
        /DuplexWeb-Google/, /googleweblight/,
        /bing/, /yandex/, /baidu/, /duckduck/, /yahoo/,           // OTHER ENGINES
        /ecosia/, /ia_archiver/,
        /facebook/, /instagram/, /pinterest/, /reddit/,          // SOCIAL MEDIA
        /slack/, /twitter/, /whatsapp/, /youtube/,
        /semrush/,                                            // OTHER
    ] as RegExp[]).map((r) => r.source).join("|"), "i")     // BUILD REGEXP + "i" FLAG

    return robots.test(userAgent)
}


export class MojeMapa extends Map<string, number> {
    async put(znak: string) {
        znak = znak.toLocaleLowerCase()

        let pocet = this.get(znak)
        if (pocet === undefined) {
            this.set(znak, 1)
        } else {
            this.set(znak, +pocet + 1)
        }
    }
    top(n: number) {
        let nejvetsi = new Map<string, number>();
        for (let i = 0; i < n; i++) {
            let nej: any[] = [undefined, 0]
            this.forEach((pocet, znak) => {
                if (pocet > nej[1] && nejvetsi.get(znak) == undefined) {
                    nej[0] = znak
                    nej[1] = pocet
                }
            })
            if (nej[0] != undefined) nejvetsi.set(nej[0], nej[1])
        }
        return nejvetsi
    }
}

export function getCisloPochvaly(rychlost: number, presnost: number) {
    if (rychlost >= levelyRychlosti[2] && presnost >= levelyPresnosti[1]) { // paradni
        return 0
    } else if (rychlost >= levelyRychlosti[1] && rychlost < levelyRychlosti[2] && presnost >= levelyPresnosti[1]) { // rychlost muze byt lepsi
        return 1
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[2]) { // presnost muze byt lepsi
        return 2
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[1] && rychlost < levelyRychlosti[2]) { // oboje muze byt lepsi
        return 3
    } else if (rychlost < levelyRychlosti[1] && presnost < levelyPresnosti[0]) { // oboje bad
        return 6
    } else if (rychlost < levelyRychlosti[1]) { // rychlost bad
        return 4
    } else if (presnost < levelyPresnosti[0]) { // presnost bad
        return 5
    }
    return 0 // nestane se
}

export function clone(obj: any) { // kvůli starším prohlížečům (koukám na tebe safari <14.0)
    let x: any
    try {
        x = structuredClone(obj)
    } catch {
        x = JSON.parse(JSON.stringify(obj))
    }
    return x
}

export function saveNastaveni(diakritika: boolean, velkaPismena: boolean, vetySlova: boolean, delka: number, klavesnice: boolean) {
    localStorage.setItem(nastaveniJmeno, JSON.stringify({ "diakritika": diakritika, "velkaPismena": velkaPismena, "vetySlova": vetySlova, "delka": delka, "klavesnice": klavesnice }))
}