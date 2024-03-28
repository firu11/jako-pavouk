#set page(numbering: none)

#[
  #set page(margin: 0in)
  #set align(center)
  
  #v(7em)

  #stack(dir: ltr, spacing: 1em)[
    #image("obrazky/logogyarab.png", width: 6em)
  ][
    #v(1.5em)
    #set align(left)
    #stack(dir: ttb, spacing: 12pt)[
      #text(20pt)[*Gymnázium Arabská, Praha 6, Arabská 14*]
    ][
      #text(16pt)[Obor programování, vyučující Mgr. Jan Lána]
    ]
  ]
  
  #v(8em)

  #box(inset: (left: 8%),
    image("obrazky/pavoukStudent.svg", width: 50%)
  )
  
  #v(5em)
  
  #text(20pt)[*Jako Pavouk*]\
  #text(15pt)[Psaní všemi deseti]
  
  #v(4em)
  
  #text(18pt)[Filip Růžička]
  
  #v(1.5em)
  
  Květen, 2024
  
  #v(1fr)

  #set page(margin: 1in)
  
  #v(1fr)
  #set align(left)
  
  Prohlašuji, že jsem jediným autorem tohoto projektu, všechny citace jsou řádně označené a všechna použitá literatura a další zdroje jsou v práci uvedené. Tímto dle zákona 121/2000 Sb. (tzv.~Autorský zákon) ve znění pozdějších předpisů uděluji bezúplatně škole Gymnázium, Praha 6, Arabská 14 oprávnění k výkonu práva na rozmnožování díla (§ 13) a práva na sdělování díla veřejnosti (§ 18) na dobu časově neomezenou a bez omezení územního rozsahu.

  #v(4em)

  #set align(right)
  
  V #box(width: 7em, repeat[.]) dne #box(width: 8em, repeat[.]) #h(1fr) Filip Růžička #box(width: 8em, repeat[.])
  #v(8em)
]

// ---------------- abstract -----------------

#set align(center)
#text(size: 17pt)[*Anotace*]\
#text(size: 12pt)[Zadání projektu]
#v(0.7em)

#pad(x: 2em)[
Webová aplikace na výuku psaní všemi deseti. Bude se skládat z lekcí přidávajících vždy 2-3 písmena na klávesnici. Tímto způsobem si člověk postupně osvojí celou klávesnici. V každé lekci budou cvičení s náhodnými písmenky, ale i se slovy obsahující pouze již naučená písmena. 

Aplikace bude taktéž obsahovat přihlašovací systém, díky kterému bude mít každý uživatel k dispozici profil pro zobrazení svých zlepšení a statistik jako např. průměrnou rychlost, chybovost, procento dokončených lekcí atd. Co se týče technologií, použiji programovací jazyk Go pro backend, framework Vue.js pro frontend a databázi PostgreSQL pro ukládání dat o uživatelích a pro přístup k slovníku spisovné češtiny.
]

