// na telefonu chci aby šel web přidat na plochu jako link, ne jako standalone aplikace
if (document.body.clientWidth <= 900) {
    document.getElementById("manifest").remove()
}

// detekce starych prohlizecu
// https://ishadeed.com/article/flexbox-gap/
var flex = document.createElement("div")
flex.style.display = "flex"
flex.style.flexDirection = "column"
flex.style.rowGap = "1px"

flex.appendChild(document.createElement("div"))
flex.appendChild(document.createElement("div"))

document.body.appendChild(flex)
var isSupported = flex.scrollHeight === 1
flex.parentNode.removeChild(flex)

if (!isSupported) alert("Váš prohlížeč je pravděpodobně příliš starý, nebo nepodporuje vše na to, aby mohl tuto aplikaci používat. Bohužel nemůže psát Jako Pavouk.")