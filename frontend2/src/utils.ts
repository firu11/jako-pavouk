export function formatovany_pismena(pismena: string | string[] | undefined) {
    if (pismena === "..." || pismena === undefined) return pismena
    let vratit = "";
    for (let i = 0; i < pismena.length; i++) {
        vratit += i < pismena.length - 1 ? pismena.at(i) + ", " : pismena.at(i);
    }
    return vratit;
}

export function get_token() {
    return localStorage.getItem("pavouk_token")
}