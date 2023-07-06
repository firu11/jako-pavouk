export function formatovanyPismena(pismena: string | string[] | undefined) { // TODO predelat na computed
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