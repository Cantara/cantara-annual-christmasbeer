const scores = new Map();

let prot = 'wss'

if (location.protocol === "http:") {
    prot = 'ws'
}
const socket = new WebSocket(prot + '://'+ location.host + '/score');

socket.onopen = () => {
    console.log('socket connected');
};

socket.onmessage = (event) => {
    if(!event.data) { return }
    const s = JSON.parse(event.data);
    lastScores = [s, ...lastScores];
    if (lastScores.length > 5) {
        lastScores = lastScores.slice(0, 5);
    }	
    const beer_id = s.beer.brand + "_" + s.beer.name + "_" + s.beer.brew_year;
    if (scores.get(beer_id)) {
        scores.set(beer_id, [...scores.get(beer_id), s]);
    } else {
        scores.set(beer_id, [s])
    }
    let hr = [];
    let mr = [];
    scores.forEach(function (v, key, map) {
        console.log("v:", v[0])
        if (hr.length < 5) {
            hr = [...hr, {beer: v[0].beer, avg: v.reduce((p,c) => p + c.rating, 0) / v.length}]
            mr = [...mr, {beer: v[0].beer, num: v.length}]
        }
    });
    highestRated = hr.sort((a,b)=> {
        if (a.avg < b.avg) {
            return 1;
        } else if (a.avg > b.avg) {
            return -1;
        }
        return 0;
    }).slice(0, 5);
    mostRated = mr.sort((a,b)=> {
        if (a.num < b.num) {
            return 1;
        } else if (a.num > b.num) {
            return -1;
        }
        return 0;
    }).slice(0, 5);
    //scrollToBottom();
    //socket.close();
};


