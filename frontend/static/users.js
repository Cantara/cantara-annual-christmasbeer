//const accounts = document.getElementById("accounts");

var prot = prot || 'wss';
if (location.protocol === "http:") {
    prot = 'ws'
}
const accountSocket = new WebSocket(prot + '://'+ location.host + '/account');

accountSocket.onopen = () => {
    console.log('account socket connected');
};

accountSocket.onmessage = (event) => {
    if(!event.data) { return }
    const account = JSON.parse(event.data)
    console.log(account)
    //accounts.appendChild(new Option(account, account));
};


