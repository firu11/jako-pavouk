#set page(
  paper: "a4",
  margin: 3cm,
  numbering: "1"
)
#set text(
  lang: "cs",
  region: "cz",
  size: 12pt,
)
#show heading.where(level: 1): it => text(
  size: 18pt,
  it
)
#show heading.where(level: 2): it => text(
  size: 16pt,
  it
)
#show heading.where(level: 3): it => text(
  size: 14pt,
  it
)

#show raw.where(block: true): it => align(block( // nastavení code bloku
  stroke: rgb("#ddd"),
  fill: rgb("#fafafa"),
  inset: 10pt,
  radius: 0.5em,
  above: 2em,
  below: 2em,
  it
), center)

#set raw(syntaxes: "custom_highlighting.yaml")

#show par: set block(spacing: 0.65em)
#set par(justify: true, first-line-indent: 2em) // hezky rádoby do bloku

#show figure: it => {it; v(0.4em)} // trošku větsí mezera za obrázkem
#set heading(numbering: "1.")

// ------------- Blbosti -------------

#let LaTeX = {
  let A = (
    offset: (
      x: -0.33em,
      y: -0.3em,
    ),
    size: 0.7em,
  )
  let T = (
    x_offset: -0.12em    
  )
  let E = (
    x_offset: -0.2em,
    y_offset: 0.23em,
    size: 1em
  )
  let X = (
    x_offset: -0.1em
  )
  [L#h(A.offset.x)#text(size: A.size, baseline: A.offset.y)[A]#h(T.x_offset)T#h(E.x_offset)#text(size: E.size, baseline: E.y_offset)[E]#h(X.x_offset)X]
}
#show "LaTeX": name => LaTeX // hahahaha

// ------------- Dokument -------------

#include "titulnistrany.typ"
#include "obsah.typ"
#counter(page).update(1)
#include "kapitoly/uvod.typ"
#pagebreak()
#include "kapitoly/psani.typ"
#include "kapitoly/problemy.typ"
#include "kapitoly/design.typ"
#pagebreak()
#include "kapitoly/do-budoucna.typ"
#include "kapitoly/zaver.typ"

#pagebreak()
#block()[
  #set par(justify: false) // citace aby se neroztahovaly
  #bibliography(
    "citace.yaml", 
    title: "Odkazy", 
    style: "the-lancet") // 1 = "the-lancet", [1] = "angewandte-chemie", i pod čarou = "gb-7714-2015-note"
]
