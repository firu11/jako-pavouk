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