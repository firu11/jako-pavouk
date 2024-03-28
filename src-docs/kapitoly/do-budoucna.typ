= Do budoucna
  == Systém pro školy
  Napadají mě dva způsoby jak aplikaci udělat *profitabilní*. Jednou možností je za pár let, až bude aplikace více známá, udělat po vzoru konkurence kurz placený. 
  
  Druhá, podle mého názoru lepší varianta, by byla přidat *systém pro učitele*, pomocí kterého by mohli například vytvářet třídy, přidávat žáky, sledovat jejich pokroky a zadávat domácí úkoly.

  Tato funkcionalita pro školy by byla placená, ať už ve formě roční licence, jednorázové platby, nebo menší částky za každého studenta.
  
  == Věty pro lekce
  Trénování na celých větách je zatím možné pouze pro uživatele, kteří si již celou klávesnici osvojili. Důvod je jednoduchý. Nemám k dispozici žádný soubor vět, obsahujících pouze *omezený výběr písmen*.

  Zkoušel jsem využít sílu, v současné době velmi populárních, LLMs (large language model), ale bohužel moc neporozuměly mému požadavku. Myslím, že tyto modely úplně nerozumí vztahu mezi písmeny a slovy (@dotaz[dotaz]). Zatím jediným způsobem, jak získat takový soubor, je napsat si věty sám.

  #figure(caption: [(březen 2024) #h(1pt) @gemini @chatgpt], supplement: "Dotaz", gap: 1em)[
    ```
    Napiš mi 10 vět obsahující pouze písmenka f,j,g,h,d,k,s,l,a,ů, q,ú,w,p,e,o,r,i,t,u,z.
  
    Odpověď Google Gemini:
    1. Fialový drak s hnědýma očima klidně seděl na louce.
    2. Hasič s houbou uhasí hořící dům.
    3. Krásná princezna s úsměvem tančila na plese.
    4. Děti si hrály na schovávanou v parku.

    Odpověď ChatGPT:
    1. Jak úžasné je, že žluté květy rostou v dáli.
    2. Vůně jasmínu zdobí chodbu hotelu.
    3. Klidně si odpočiň pod stromem.
    4. Fénix je symbol obnovy a síly.
    ```
  ]<dotaz>
  