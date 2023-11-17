const scoreForm = document.getElementById("score-form");
const scoreButton = document.getElementById("score-submit");
const scoreErrorMsg = document.getElementById("score-error-msg");

scoreButton.addEventListener("click", (e) => {
    e.preventDefault();
    scoreErrorMsg.style.opacity = 0;
    const body = {
        rating: parseInt(scoreForm.rating.value),
        comment: scoreForm.comment.value,
    }
    fetch('/score/' + new Date().getFullYear() + "/" + scoreForm.beer.value, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        body: JSON.stringify(body),
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        console.log(response)
        if (response.status !== 204) {
            return response.json()
        }
    })
    .then(data => {
        console.log(data);
        if (data.error) {
            scoreErrorMsg.innerText = data.error;
            scoreErrorMsg.style.opacity = 1;
            return;
        }
    })
    .catch((error) => {
        scoreErrorMsg.innerText = error;
        scoreErrorMsg.style.opacity = 1;
        console.log(error);
    });
})
