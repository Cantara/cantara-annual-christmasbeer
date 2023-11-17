const summary = document.getElementById("summary");

var prot = prot || 'wss';
if (location.protocol === "http:") {
    prot = 'ws'
}
const scoreSocket = new WebSocket(prot + '://'+ location.host + '/score/sumary');

scoreSocket.onopen = () => {
    console.log('socket connected');
};

scoreSocket.onmessage = (event) => {
    if(!event.data) { return }
    console.log(JSON.parse(event.data))
    //summary.outerHTML = JSON.parse(event.data);
    summary.innerHTML = JSON.parse(event.data);
};


