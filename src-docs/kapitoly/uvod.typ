= Úvod
  == Výběr tématu
  Nápad na vývoj této webové aplikace byl fakt, že já sám *jsem neuměl psát všemi deseti*, a tak jsem se to rozhodl naučit. Se zděšením jsem ale zjistil, že co se týče české klávesnice, není moc možností, přičemž nějaké z nich stojí skoro 1000 Kč.

  Ten fakt, že má můj projekt šanci uspět a nebýt jen další zapomenutá ročníková práce, pro mě byl motivací.
  
  == Použité technologie
  #figure(
    image("../obrazky/technologie.png", width: 90%),
    caption: [Vue.js @vue-logo,#h(5pt) PostgreSQL @postgres-logo,#h(5pt) Go @go-logo],
  )
  
    === Frontend
    Pro vývoj frontendu neboli samotné webové stránky jsem si vybral framework *Vue.js*. Hlavním důvodem bylo, že už jsem s ním měl zkušenosti ze skupinového projektu pro minulý rok.
    
    Místo Javascriptu jsem používal *Typescript*, ale nejsem si jist, jestli ho mohu doporučit. Čas od času jsem měl pocit, že musím psát hodně kódu navíc a že to přináší *více škody než užitku*.

    #pagebreak()
    Co se týče kaskádových stylů, držel jsem se *čistého CSS*, protože si myslím, že s nástroji jako Bootstrap nebo Tailwind ztrácíte jak čitelnost, tak kontrolu nad kódem. Přijde mi, že CSS se nějakým způsobem odděluje od samotného HTML z nějakého důvodu a nechci se vracet k tomuto:

    ```html
    <nav>
      <ul style="display: flex; justify-content: center; align-items: 
      center; max-width: 720px; margin: 0 auto; height: 100vh; 
      list-style: none;">
        <li style="width: 125px; height: 50px; transition: 
        background-position-x 0.9s linear; text-align: center;">
          <a href="/" style="font-size: 22px; color: #777; 
          text-decoration: none; transition: all 0.45s;">Home</a>
        </li>
        <li style="width: 125px; height: 50px; transition: 
        background-position-x 0.9s linear; text-align: center;">
          <a href="/about" style="font-size: 22px; color: #777; 
          text-decoration: none; transition: all 0.45s;">About</a>
        </li>
      </ul>
    </nav>
    ```
  
    === Backend
    Jako jazyk pro backend, program starající se o data, jsem si zvolil *Go*, protože se mi v té době hodně líbil Python a chtěl jsem zkusit něco podobného, ale více low-level. Dnes se mi ale zpátky nechce.
  
    Použil jsem také framework Fiber, i když teď cítím, že bych se v klidu obešel bez něj, pouze za pomoci *standardní knihovny*. 
    
    Ještě jsem totiž nevěděl, jak moc dobrou standardní knihovnu Go má a že mindset Go programátorů je úplným opakem JS developerů. Když chce člověk něco udělat v Javascriptu, většinou už existuje asi 287 knihoven, které všechny dělají to samé trošku jinak. Nějaká lépe, nějaká hůř. Vybereme si tedy jednu z nich a začneme bezhlavě importovat cizí programy, kterým nerozumíme a ani nechceme rozumět.
  
    V Go to tak není. Cokoliv, co potřebujete od regulárních výrazů, přes unit testování, datum a čas až po komunikaci se sítí už je v standardní knihovně, psané lidmi, kteří rozumí nejen jazyku samotnému, ale i danému problému.
      
    === Databáze
    Pro uchovávání dat o uživatelích, slovech pro psaní a podobně jsem zvolil relační databázi *PostgreSQL*, kvůli její popularitě a výkonu.
  
    Ve výběrovém konkurzu byly i alternativy jako MongoDB nebo SQLite. Mongo jsem ale zavrhl, protože jsem chtěl spíše relační SQL databázi a SQLite mi zase v produkci přišlo hloupé.
  
    === Dokumentace
    Pro psaní této dokumentace jsem využil *Typst* _"A markup-based typesetting system that is designed to be as powerful as_ LaTeX _ while being much easier to learn and use."_ #h(3pt) @typst Chtěl jsem se vyhnout stylování v MS Word, a LaTeX mi přišel dost zkostnatělý. Chvilku jsem tedy hledal a nakonec zkusil moderní, také open-source, alternativu.