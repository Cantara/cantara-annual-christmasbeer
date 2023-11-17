const beers = document.getElementById("beers");

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
    beers.appendChild(new Option(beer, beer));
};


