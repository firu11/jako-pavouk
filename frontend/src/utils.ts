import { ref } from "vue";
import { cislaProcvicJmeno, tokenJmeno } from "./stores";

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

export const oznameni = ref([] as { text: String }[])

export function pridatOznameni(text: string = "Něco se pokazilo", cas: number = 4000) {
    let obj = { text: text }
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
        if (oznameni.value.length == 0) {
            pridatOznameni("Dej si čajík a vydýchej se...")
        }
        return true
    }
    return false
}

export class Oznacene {
    index = ref(0)
    max: number = 3
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
    ] as RegExp[]).map((r) => r.source).join("|"), "i");     // BUILD REGEXP + "i" FLAG

    return robots.test(userAgent);
};
