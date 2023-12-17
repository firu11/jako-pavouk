import { ref } from "vue";

export function formatovanyPismena(pismena: string | string[] | undefined) {
    if (pismena === "..." || pismena === undefined) return pismena
    let vratit = "";
    for (let i = 0; i < pismena.length; i++) {
        vratit += i < pismena.length - 1 ? pismena.at(i) + ", " : pismena.at(i);
    }
    return vratit;
}

export function getToken() {
    return localStorage.getItem("pavouk_token")
}

export const oznameni = ref([] as { text: String }[])

export function pridatOznameni(text: string = "Něco se pokazilo", cas: number = 4000) {
    let obj = { text: text }
    oznameni.value.push(obj)
    setTimeout(() => {
        oznameni.value.splice(oznameni.value.indexOf(obj), 1);
    }, cas)
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