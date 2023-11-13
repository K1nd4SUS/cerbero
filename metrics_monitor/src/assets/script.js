async function fetchData() { // calling Cerbero API and then call the fillData procedure
    fetch("http://127.0.0.1:8082/metrics").then(async (answ) => {
        if (answ.ok) {
            data = await answ.json()
            fillData(data)
        }
    })
}

function fillData(data) {

    let serviceCard = document.getElementsByClassName("serviceCard")[0].cloneNode(true) // single service card
    let generalInfo = document.getElementById("generalInfo").cloneNode(true) // generic services info

    let mainContent = document.getElementById("mainContent") // service card container
    mainContent.innerHTML = ""

    generalInfo.getElementsByTagName("h6")[0].innerHTML = `File edits: <span class="highlight">${data.FileEdits}</span>`
    generalInfo.getElementsByTagName("h6")[1].innerHTML = `Registered services: <span class="highlight">${data.ServiceAccess.length}</span>`

    mainContent.append(generalInfo) // setting general info

    data.ServiceAccess.forEach(element => { // creating and then appending a single service card for each service 
        newServiceCard = serviceCard.cloneNode(true)
        newServiceCard.getElementsByTagName("h5")[0].innerHTML = `${element.Service.name}:<span class="highlight">${element.Service.dport}</span>`
        newServiceCard.getElementsByTagName("h6")[0].innerHTML = `${element.Service.protocol}`
        newServiceCard.getElementsByTagName("h6")[1].innerHTML = `Nfq ID: ${element.Service.Nfq}`

        let blackList = newServiceCard.getElementsByClassName("rulesList")[0].getElementsByTagName("p")[0] // setting blacklist and whitelist ruleslist in the card
        blackList.innerHTML = ""
        element.Service.rulesList.blacklist.forEach(element => {
            blackList.innerHTML += `<span class="highlight">${element.type}</span>:</br>${element.filters.join("<br>")}`
        })

        let whitelist = newServiceCard.getElementsByClassName("rulesList")[0].getElementsByTagName("p")[1]
        whitelist.innerHTML = ""
        element.Service.rulesList.whitelist.forEach(element => {
            whitelist.innerHTML += `<span class="highlight">${element.type}</span>:</br>${element.filters.join("<br>")}`
        })

        hitsCard = newServiceCard.getElementsByClassName("hits")[0] // setting the hits card, divided based on the request method
        if (element.Hits != null) {
            let getsCol = hitsCard.getElementsByClassName("gets")[0]
            let postsCol = hitsCard.getElementsByClassName("posts")[0]
            let putsCol = hitsCard.getElementsByClassName("puts")[0]
            let othersCol = hitsCard.getElementsByClassName("others")[0]

            getsCol.innerHTML = ""
            postsCol.innerHTML = ""
            putsCol.innerHTML = ""
            othersCol.innerHTML = ""

            element.Hits.forEach(element => {
                if (element.Method == "GET") {
                    getsCol.innerHTML += `<span class="getText">${element.Method}</span> - ${element.Resource}: ${element.Counter}, <span class="blacklist">${element.Blocked}</span><br>`
                }
                else if (element.Method == "POST") {
                    postsCol.innerHTML += `<span class="postText">${element.Method}</span> - ${element.Resource}: ${element.Counter}, <span class="blacklist">${element.Blocked}</span><br>`
                }
                else if (element.Method == "PUT") {
                    putsCol.innerHTML += `<span class="putText">${element.Method}</span> - ${element.Resource}: ${element.Counter}, <span class="blacklist">${element.Blocked}</span><br>`
                }
                else othersCol.innerHTML += `<span class="othersText">${element.Method}</span> - ${element.Resource}: ${element.Counter}, <span class="blacklist">${element.Blocked}</span><br>`
            })
        }
        else hitsCard.innerHTML = "No hits"

        newServiceCard.append(hitsCard)


        newServiceCard.classList.remove("d-none")
        mainContent.append(newServiceCard)
    });
}

let refreshInterval = null
let textRefreshInterval = null
let nextRefreshCounter
let refreshSeconds = 0

function updateRefresh(seconds) { // updating the refresh time
    let button
    if (seconds === 0) {
        button = document.getElementsByTagName("button")[4]
        document.getElementById("nextRefresh").classList.add("d-none")
        clearInterval(refreshInterval) // if the refresh was previously set, clearing the interval
        clearInterval(textRefreshInterval)
    }
    else {
        document.getElementById("nextRefresh").classList.remove("d-none")
        if(seconds === 5){
            button = document.getElementsByTagName("button")[0]
        }
        else if (seconds === 20) {
            button = document.getElementsByTagName("button")[1]
        }
        else if (seconds === 60) {
            button = document.getElementsByTagName("button")[2]
        }
        else {
            button = document.getElementsByTagName("button")[3]
        }
        refreshSeconds = seconds - 1 // tracking the refresh seconds also with this global variable: useful for the refresh text

        clearInterval(refreshInterval) // if the refresh was previously set, clearing the interval
        clearInterval(textRefreshInterval)
        refreshInterval = setInterval(fetchData, seconds * 1000) // setting the new refresh time
        nextRefreshCounter = refreshSeconds
        textRefreshInterval = setInterval(showRefreshText, 1000)
    }

    Array.from(document.getElementsByTagName("button")).forEach(element => {
        element.style.backgroundColor = ""
    })
    button.style.backgroundColor = "#1cf1fb"


}

function showRefreshText() {
    let refreshText = document.getElementById("nextRefresh")
    refreshText.innerHTML = `Next in ${nextRefreshCounter}s`
    nextRefreshCounter--
    if (nextRefreshCounter == 0) nextRefreshCounter = refreshSeconds
}



fetchData()
updateRefresh(0) // setting no refresh as default