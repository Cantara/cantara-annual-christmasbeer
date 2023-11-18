const beers = document.getElementById("beers");
const beerTable = document.getElementById("beerTable");

var prot = prot || 'wss';
if (location.protocol === "http:") {
    prot = 'ws'
}
const beerSocket = new WebSocket(prot + '://'+ location.host + '/beer');

beerSocket.onopen = () => {
    console.log('beer socket connected');
};

beerSocket.onmessage = (event) => {
    if(!event.data) { return }
    const beer = JSON.parse(event.data)
    console.log(beer)
    beers.appendChild(new Option(beer.brand + " " + beer.name + " " + beer.brew_year, beer.id));
    const row = beerTable.insertRow(1);
    row.insertCell(0).innerText = beer.brand;
    row.insertCell(1).innerText = beer.name;
    row.insertCell(2).innerText = beer.brew_year;
    row.insertCell(3).innerText = beer.abv;
};


