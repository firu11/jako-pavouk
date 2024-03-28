#set page(numbering: none)

#show outline.entry: it => par(first-line-indent: 0em, // krejzy úprava obsahu https://stackoverflow.com/questions/77031078/how-to-remove-numbers-from-outline
  if it.at("label", default: none) == <modified-entry> {
      it // prevent infinite recursion
  } else if it.level == 1 {
    v(8pt)
    strong(text(13.5pt)[#outline.entry(it.level, it.element, it.body, [], it.page) <modified-entry>])
  } else if it.level == 2 {
    v(-4pt)
    h(1.4em)
    it
  } else if it.level == 3 {
    v(-4pt)
    h(2.8em)
    it
  }
)

#align(center)[#text(size: 17pt)[*Obsah*]]
#v(-1em)
#outline(title: "")