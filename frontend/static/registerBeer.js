const beerForm = document.getElementById("beer-form");
const beerButton = document.getElementById("beer-submit");
const beerErrorMsg = document.getElementById("beer-error-msg");

beerButton.addEventListener("click", (e) => {
    e.preventDefault();
    beerErrorMsg.style.opacity = 0;
    const body = {
        brand: beerForm.brand.value,
        name: beerForm.name.value,
        brew_year: parseInt(beerForm.year.value),
        abv: parseFloat(beerForm.abv.value),
    }
    fetch('/beer/' + encodeURI(body.brand + "_" + body.name + "_" + body.brew_year), {
        method: 'PUT',
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
        if (response.status !== 200) {
            return response.json()
        }
        beerForm.reset();
    })
    .then(data => {
        console.log(data);
        if (data?.error) {
            beerErrorMsg.innerText = data.error;
            beerErrorMsg.style.opacity = 1;
            return;
        }
    })
    .catch((error) => {
        beerErrorMsg.innerText = error;
        beerErrorMsg.style.opacity = 1;
        console.log(error);
    });
})
