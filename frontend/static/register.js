const registerForm = document.getElementById("register-form");
const registerButton = document.getElementById("register-submit");
const registerErrorMsg = document.getElementById("register-error-msg");

registerButton.addEventListener("click", (e) => {
    e.preventDefault();
    registerErrorMsg.style.opacity = 0;
    const body = {
        username: registerForm.username.value,
        firstname: registerForm.name.value,
        password: registerForm.password.value,
    }
    fetch('/account/' +  body.username, {
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
        if (response.status === 204) {
            registerForm.reset();
            location.reload();
        } else {
            return response.json()
        }
    })
    .then(data => {
        console.log(data);
        if (data?.error) {
            registerErrorMsg.innerText = data.error;
            registerErrorMsg.style.opacity = 1;
            return;
        }
    })
    .catch((error) => {
        registerErrorMsg.innerText = error;
        registerErrorMsg.style.opacity = 1;
        console.log(error);
        //message.error(error);
    });
})
