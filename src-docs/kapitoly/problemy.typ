= Problematika
  == Slovník
  Jednou z prvních komplikací bylo najít nějaký dobrý slovník (respektive seznam českých slov, ne jejich významy). Jediný takový seznam, který jsem byl schopen dohledat, byl Český národní korpus SYN2015 o cca 70 000 slovech, který jsem si zatím zredukoval asi na *61 000* a stále čas od času narazím na nějaké cizí nebo zvláštní slovo. Ručně projít tisíce slov totiž není žádný med.

  Věty objevující se v testu psaní jsem posbíral z webu Pohádkozem @pohadkozem. Jejich pohádky jsou ze starších knih, jejichž autoři už jsou více než 70 let po smrti.
  
  == Fungování na mobilních zařízeních
  Fungování aplikace na výuku psaní všemi deseti na mobilních zařízeních je docela *paradox*. Přestože samotné psaní nedává smysl na malém displeji, musí být stránka přístupná na jakémkoli stroji. V roce 2023 asi 3/4 návštěvnosti webu totiž pocházelo z telefonů. @pc-mobile

  Jako Pavouk je tedy responzivní a plně funkční na malých displejích, co se obsahu, přihlašovacího systému, statistik a dalšího týče. Když se ale uživatel pokusí vstoupit do samotného přepisování textů, setká se s hláškou _"Psaní na telefonech zatím neučíme…"_
  
  Aplikace detekuje mobilní zařízení jednoduše podle šířky HTML dokumentu. Druhou největsí modifikací webu je asi menu, které se schovává do strany a je přístupné přes kulaté tlačítko v rohu obrazovky.
  
  == Dvě různá rozložení klávesnice
  V České republice jsou běžně používány dvě různá rozložení klávesnice a to QWERTZ a QWERTY. Jedním z problémů, které nevycházely z technického provedení, ale z čistého konceptu, byla implementace těchto dvou rozložení. Nemyslím tím samotnou grafickou klávesnici na frontendu, kde stačilo prohodit Z a Y, ale problém se slovy pro každou lekci zvlášť.
  
  Moje první *nedomyšlené řešení*, které jsem psal ještě v době, kdy stránka disponovala jen texty z náhodně složených písmen, spočívalo v záměně písmenek po vygenerování textu. Tento způsob však logicky fungoval jen pro texty složené z náhodných písmen („ttzz zttt ztzz“) a záměna těchto písmen, ve cvičeních se slovy nebyla ideální („kůylata yůstat důkay“).
  
  #block(breakable: false)[
  V databázi máme relaci „slovnik“, která jak název napovídá, obsahuje všechna slova používaná k trénování. V lekci je ale nutné používat pouze slova složená z *již naučených* písmen. Proto jsem si lekce rozdělil na skupiny podle klávesnice, pro kterou jsou vytvořeny a každému slovu přiřadil id lekce, ve které už je možné ho napsat. Níže v ukázce dat vlevo z databáze můžeme vidět jednotlivé lekce a typ klávesnice, ke kterému patří. Potom v druhé tabulce vpravo vidíme slova ze slovníku a „id“ lekcí, do kterých patří v závislosti na rozložení klávesnice.

  ```asm
                lekce               |              slovnik
               ¯¯¯¯¯¯¯              |             ¯¯¯¯¯¯¯¯¯
   id |    pismena    | klavesnice  |     slovo | qwertz_id | qwerty_id
  ----+---------------+------------ | ----------+-----------+-----------
                 ...                |
    4 | aů            | oboje       |                ...
    5 | tz            | qwertz      | autorství |        17 | 17
    6 | ty            | qwerty      |    autorů |         9 | 9
    7 | ru            | oboje       |      auty |        13 | 7
                 ...                |    autory |        13 | 9
   12 | cn            | oboje       |    autoři |        16 | 16
   13 | yxm           | qwertz      |                ...
   14 | zxm           | qwerty      |
   15 | žý            | oboje       |
                 ...                |
  ```
  Jako příklad se můžeme podívat na slovo `auty`, které do množiny již použitelných slov nepřichází ve stejné lekci. U varianty QWERTY totiž už umíme všechna potřebná písmena u 7. lekce. Na druhou stranu u QWERTZ, kde u 7. lekce ještě neumíme písmeno Y, se slovo `auty` objeví až ve 13. lekci, kdy se Y naučíme.
  ]

  #block(breakable: false)[ // aby se to nezalomilo na další stranu
  == Optimalizace
  *SEO* (Search Engine Optimization) je něco, čemu by se většina vývojářů raději vyhnula. Bohužel se ale bez toho vaše aplikace k nikomu nedostane. Proto jsem se i na tuto stránku vývoje zaměřil a optimalizoval všechny aspekty aplikace tak, aby byly všechny vyhledávače spokojené a aby hladce běžela na jakkoliv starém hardwaru (testováno na Macbooku, 2008, Intel Core 2 Duo, 4GB RAM, linux).

  PageSpeed Insights (běžící na open-source Lighthouse)@lighthouse je *nástroj od Google*, který poskytuje statistiku o výkonu a celkové použitelnosti jakékoliv webové stránky. Nejdůležitější jsou ale doporučení a způsoby, jak veškeré problémy se stránkou vyřešit. I to tento nástroj poskytuje. Poté co jsem prošel každý jeden problém a obětoval pár hodin svého času, jsem dosáhl výsledku níže.

  #figure(
    image("../obrazky/seo.png", width: 90%),
    caption: [PageSpeed Insights screenshot],
  ) <seo>
  ]

  Trošku nešťastné je *jméno "Jako Pavouk"*, které sice hezky vystihuje aplikaci a umožňuje mi používat pavouka jako maskota, ale překrývá se se stejnojmenným filmem, a tak boj o první příčky ve výsledcích hledání není jednoduchý.
  
  Co se týče samotných vyhledávačů, nejpopulárnější Google má v tuto chvíli paradoxně moji aplikaci rád nejméně. Alternativy jako DuckDuckGo, Qwant nebo Bing zobrazují stránku mnohem častěji a jako jednu z prvních. Naštěstí se ale i pozice na Google vyhledávači postupně zlepšuje.

  == Výběr jmen uživatelů
  Jelikož Jako Pavouk disponuje i přihlašováním pomocí účtů Google, musel jsem se potýkat s otázkou, jaké *unikátní jméno* uživateli přiřadit pro naši databázi, která nemá stejné požadavky na jeho tvar. Od googlu po přihlášení získám kromě e-mailu i jméno a příjmení uživatele, ze kterého mohu vytvořit jméno pro naši aplikaci.

  První implementace vypadala tak, že jsem vytáhl jména všech uživatelů z databáze, v lineárním čase je přidal do hashovací tabulky a poté v konstantním čase zkoušel různé kombinace jména a příjmení. Datovou strukturu jsem sice použil, abych se nemusel zbytečně dotazovat do databáze, ale když jsem se podruhé zamyslel, nebylo to nejoptimálnější řešení, jelikož šance, že jméno už bude používáno, je velmi malá. Proto bylo často zbytečné vytvářet hash tabulku jen kvůli jedné jediné kontrole.

  #block(breakable: false)[
  Způsob, kterým jsem nakonec první řešení nahradil, spočívá v *přidání náhodného čísla za jméno* uživatele. V tomto případě je velmi malá pravděpodobnost kolize v databázi, a ještě menší, že tato kolize nastane vícekrát. Díky tomuto vylepšení není potřeba načítat všechna jména uživatelů a ani na databázi nebude příliš velký nápor.
  ]