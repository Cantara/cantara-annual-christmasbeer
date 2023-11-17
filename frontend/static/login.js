const loginForm = document.getElementById("login-form");
const loginButton = document.getElementById("login-submit");
const loginErrorMsg = document.getElementById("login-error-msg");

//const t = localStorage.getItem("token");

loginButton.addEventListener("click", (e) => {
    loginErrorMsg.style.opacity = 0;
    e.preventDefault();
    const body = {
        username: loginForm.username.value,
        password: loginForm.password.value,
    }
    fetch('/account/login', {
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
        if (response.status === 204) {
            location.reload();
        } else {
            return response.json()
        }
    })
    .then(data => {
        console.log(data);
        if (data.error) {
            loginErrorMsg.innerText = data.error;
            loginErrorMsg.style.opacity = 1;
            return;
        }
    })
    //.then(response => response.json())
    //.then(data => {
    //    console.log(data)
    //    if (data.error) {
    //        console.log(data.error);
    //        //message.error(data.error);
    //        return;
    //    }
    //    //token.set(data)
    //    localStorage.setItem("token", data);
    //})
    .catch((error) => {
        loginErrorMsg.innerText = error;
        loginErrorMsg.style.opacity = 1;
        //console.log(error);
        //message.error(error);
    });
})
