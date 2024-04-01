= Psaní všemi deseti
  == Rozložení klávesnice
  Klávesnice, kterou používáme každý den na počítači nebo telefonu, není vůbec uzpůsobena na rychlé ani pohodlné psaní. Proč to tak ale je? Proč používáme zrovna rozložení QWERTZ a QWERTY? Proč ne třeba ABCDEF?

  Všechny cesty vedou až do roku 1878 k *psacímu stroji* Remington @psaci-stroj, který se objevil s nám známou QWERTY a o deset let později bylo toto rozložení dokonce přijato za standard. Konstruktéři se při vytváření klávesnice k prvním psacím strojům totiž potýkali s *problémem zasekávajících se kladívek* stisknutých rychle za sebou. Docházelo k tomu ještě častěji, pokud byly tyto dvě klávesy vedle sebe. @psaci-stroj2 Proto bylo nutné klávesy, které se často píší za sebou sobě, "rozházet" dál od sebe a vznikla tak nám známá QWERTY.
  
  V Česku nebo například v Německu se převážně používá lehce modifikovaná QWERTZ, jednoduše proto, že v *němčině* se písmeno Z vyskytuje daleko častěji než Y a tak si tyto klávesy Němci prohodili, aby měli Z blíže po ruce. V té době se na našem území psalo převážně německy, a tak jsme tuto úpravu přijali. @psaci-stroj2

  Přestože dnes už takový problém s kladívky nemáme, zvyk je železná košile.
  
  #figure(
    image("../obrazky/psaci-stroj.jpg", width: 80%),
    caption: [Psací stroj @psaci-stroj-ilustrace],
  ) <stroj>

  #block(breakable: false)[ // aby se to nezalomilo na další stranu
  == Prstoklad
  Když už víme, jak vznikla naše klávesnice, pojďme se podívat, jak rozmístit prsty, abychom ji celou ovládli. Záchranným bodem pro nás jsou klávesy F a J ležící téměř ve středu. Ty na drtivé většině klávesnic mají *malé výstupky*, hmatné i poslepu. Na tyto dvě hlavní umístíme ukazováčky a ostatní prsty na klávesy ve stejné řadě. Vycházet tedy budeme z tlačítek ASDF a JKLŮ. Každý prst si poté hlídá svůj *pomyslný sloupeček*, jak můžete vidět na @klavesnice[obrázku]. Výjimkou jsou palce. Ty se při psaní starají pouze o mezerník.
  
  #v(0.6em)
  #figure(
    image("../obrazky/klavesniceSPavoukem.png", width: 70%),
    caption: "Prstoklad",
  ) <klavesnice>
  ]

  == Implementace
  Některé starší programy na výuku psaní všemi deseti fungovaly na bázi opisu textu o řádek níže. Poté porovnávají zadání s přepisem od uživatele. Tento způsob ale vede k častému dezorientování a k námaze očí z ustavičného koukání tam a zpět.
  
  #v(-4pt)
  ```custom
  ffjj jjff fjfj jfjf fffj jfff     <- zadání
   
  ffjj jjff fj|                     <- přepis
  ```
  #v(-4pt)
  
  Jako Pavouk k psaní přistupuje trošku jinak. Uživatel text neopisuje, ale vše se děje přímo před jeho očima. Přesněji se uživatel orientuje podle podtržení, které označuje vždy další písmeno, které má stisknout. Text již napsaný je upozaděn šedým odstínem písma, aby se uživatel soustředil pouze na slova nadcházející, která jsou bílá. Špatné písmeno je značeno červenou barvou a je na něj upozorněno zvukem.

  Moje aplikace mimo jiné monitoruje nejen rychlost psaní, přesnost nebo čas, ale i písmena, ve kterých uživatel nejčastěji chybuje. Tato statistika je mu pak samozřejmě k nahlédnutí.