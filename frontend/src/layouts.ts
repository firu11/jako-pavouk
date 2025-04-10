const barvicky = ["#6ada56", "#81bffc", "#fa5ca1", "#ff8800", "#6f86f7", "#bc73ff"]
export const barvy = new Map([["P_Ukaz", barvicky[0]], ["L_Ukaz", barvicky[1]], ["P_Pros", barvicky[2]], ["L_Pros", barvicky[2]], ["P_Prs", barvicky[3]], ["L_Prs", barvicky[3]], ["P_Mali", barvicky[4]], ["L_Mali", barvicky[4]], ["Palce", barvicky[5]]])

// windows
export const schemaWindows = [
    ["°;", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⟵"],
    ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "'¨"],
    ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "Enter ↵"],
    ["Shift", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift"],
    ["Ctrl", "", "", "Alt", "______", "", "", "", "∧∨", ""]
]
export const delkaKlavesWindows: { [id: string]: number } = { "⟵": 3, "Shift": 1, "Enter ↵": 1, "CapsLock": 1, "TAB": 1, "______": 24, "Ctrl": 2, "   ": 2 }
export const prstokladWindows: { [id: string]: string[] } = {
    "P_Ukaz": ["Z", "J", "H", "U", "N", "M", "7ý", "8á"],
    "L_Ukaz": ["G", "T", "R", "F", "V", "B", "5ř", "6ž"],
    "P_Pros": ["K", "I", "?,", "9í"],
    "L_Pros": ["D", "E", "C", "4č"],
    "P_Prs": ["O", "L", ":.", "0é"],
    "L_Prs": ["X", "S", "W", "3š"],
    "P_Mali": ['"ů', "P", "_-", '%=', 'ˇ´', '⟵', '()', '/ú', "'¨", '!§', 'Enter ↵', 'Shift'],
    "L_Mali": ["Y", "Shift", "A", "Q", "1+", "°;", "2ě", "TAB", "CapsLock", "Ctrl"],
    "Palce": ["______", "Alt"]
}
export const specialniZnakyWindowsQWERTZ = new Map<string, string[]>(Object.entries({ "Ů": ["°;", "U"], "@": ["Ctrl", "Alt", "V"], "[": ["Ctrl", "Alt", "F"], "]": ["Ctrl", "Alt", "G"], "{": ["Ctrl", "Alt", "B"], "}": ["Ctrl", "Alt", "N"], "#": ["Ctrl", "Alt", "X"], "$": ["Ctrl", "Alt", "\"ů"], "~": ["Ctrl", "Alt", "1+"], "^": ["Ctrl", "Alt", "3š"], "&": ["Ctrl", "Alt", "C"], "*": ["Ctrl", "Alt", "_-"], "<": ["Ctrl", "Alt", "?,"], ">": ["Ctrl", "Alt", ":."], "`": ["Ctrl", "Alt", "7ý"], "|": ["Ctrl", "Alt", "W"], "\\": ["Ctrl", "Alt", "Q"] }))
export const specialniZnakyWindowsQWERTY = new Map<string, string[]>(Object.entries({ "Ů": ["°;", "U"], "@": ["Ctrl", "Alt", "2ě"], "[": ["Ctrl", "Alt", "/ú"], "]": ["Ctrl", "Alt", "()"], "{": ["Ctrl", "Alt", "Shift", "/ú"], "}": ["Ctrl", "Alt", "Shift", "()"], "#": ["Ctrl", "Alt", "3š"], "$": ["Ctrl", "Alt", "4č"], "~": ["Ctrl", "Alt", "Shift", "°;"], "^": ["Ctrl", "Alt", "6ž"], "&": ["Ctrl", "Alt", "7ý"], "*": ["Ctrl", "Alt", "8á"], "<": ["Ctrl", "Alt", "?,"], ">": ["Ctrl", "Alt", ":."], "`": ["Ctrl", "Alt", "°;"], "|": ["Ctrl", "Alt", "Shift", "'¨"], "\\": ["Ctrl", "Alt", "'¨"] }))

// macos
export const schemaMacOS = [
    ["><", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⌫"],
    ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "↵"],
    ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "`¨", "enter-noha"],
    ["Shift1", "|\\", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift2"],
    ["", "", "⌥", "", "______", "", "⌥", "", "∧∨", ""]
]
export const delkaKlavesMacOS: { [id: string]: number } = { "⌫": 3, "Shift1": 3, "Shift2": 8, "CapsLock": 1, "TAB": 1, "______": 24 }
export const prstokladMacOS: { [id: string]: string[] } = {
    "P_Ukaz": ["Z", "J", "H", "U", "N", "M", "7ý", "8á"],
    "L_Ukaz": ["G", "T", "R", "F", "V", "B", "5ř", "6ž"],
    "P_Pros": ["K", "I", "?,", "9í"],
    "L_Pros": ["D", "E", "C", "4č"],
    "P_Prs": ["O", "L", ":.", "0é"],
    "L_Prs": ["X", "S", "W", "3š"],
    "P_Mali": ["\"ů", "P", "_-", "%=", "ˇ´", "⌫", "()", "/ú", "`¨", "!§", "↵", "enter-noha", "Shift2"],
    "L_Mali": ["Y", "Shift1", "|\\", "A", "Q", "1+", "><", "2ě", "TAB", "CapsLock"],
    "Palce": ["______", "⌥"]
}
export const specialniZnakyMacOS = new Map<string, string[]>(Object.entries({ "Ů": ["⌥", "%=", "U"], "@": ["⌥", "2ě"], "#": ["⌥", "3š"], "$": ["⌥", "4č"], "~": ["⌥", "5ř"], "^": ["⌥", "6ž"], "&": ["⌥", "7ý"], "*": ["⌥", "8á"], "{": ["⌥", "9í"], "}": ["⌥", "0é"], "°": ["⌥", "%="], "[": ["⌥", "/ú"], "]": ["⌥", "()"], ";": ["⌥", "\"ů"], "'": ["⌥", "!§"] }))

// linux
export const schemaLinux = [
    ["°;", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⟵"],
    ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "'¨"],
    ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "Enter ↵"],
    ["Shift", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift"],
    ["", "", "", "", "______", "AltGr", "", "", "∧∨", ""]
]
export const prstokladLinux: { [id: string]: string[] } = {
    "P_Ukaz": ["Z", "J", "H", "U", "N", "M", "7ý", "8á"],
    "L_Ukaz": ["G", "T", "R", "F", "V", "B", "5ř", "6ž"],
    "P_Pros": ["K", "I", "?,", "9í"],
    "L_Pros": ["D", "E", "C", "4č"],
    "P_Prs": ["O", "L", ":.", "0é"],
    "L_Prs": ["X", "S", "W", "3š"],
    "P_Mali": ['"ů', "P", "_-", '%=', 'ˇ´', '⟵', '()', '/ú', "'¨", '!§', 'Enter ↵', 'Shift'],
    "L_Mali": ["Y", "Shift", "A", "Q", "1+", "°;", "2ě", "TAB", "CapsLock", "Ctrl"],
    "Palce": ["______", "AltGr"]
}
export const specialniZnakyLinuxQWERTZ = new Map<string, string[]>(Object.entries({ "Ů": ["°;", "U"], "@": ["AltGr", "V"], "[": ["AltGr", "F"], "]": ["AltGr", "G"], "{": ["AltGr", "B"], "}": ["AltGr", "N"], "#": ["AltGr", "3š"], "$": ["AltGr", "\"ů"], "~": ["AltGr", "A"], "^": ["AltGr", "6ž"], "&": ["AltGr", "7ý"], "*": ["AltGr", "8á"], "<": ["AltGr", "?,"], ">": ["AltGr", ":."], "`": ["AltGr", "°;"], "|": ["AltGr", "W"], "\\": ["AltGr", "Q"] }))
export const specialniZnakyLinuxQWERTY = new Map<string, string[]>(Object.entries({ "Ů": ["°;", "U"], "@": ["AltGr", "2ě"], "[": ["AltGr", "/ú"], "]": ["AltGr", "()"], "{": ["AltGr", "Shift", "/ú"], "}": ["AltGr", "Shift", "()"], "#": ["AltGr", "3š"], "$": ["AltGr", "4č"], "~": ["AltGr", "Shift", "°;"], "^": ["AltGr", "6ž"], "&": ["AltGr", "7ý"], "*": ["AltGr", "8á"], "<": ["AltGr", "?,"], ">": ["AltGr", ":."], "`": ["AltGr", "°;"], "|": ["AltGr", "Shift", "'¨"], "\\": ["AltGr", "'¨"] }))
