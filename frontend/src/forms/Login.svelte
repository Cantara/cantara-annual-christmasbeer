<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";

  import { token } from '../stores/token.js';

  function login() {
    fetch('/account/login', {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'omit',
      body: JSON.stringify(body),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
    })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      if (data.error) {
        console.log(data.error);
        //message.error(data.error);
        return;
      }
      token.set(data)
    })
    .catch((error) => {
      console.log(error);
      //message.error(error);
    });
  }

  let body = {
    username: "",
    password: "",
  }

  let valid_username = false;
  let valid_password = false;

  let disabled = false;

$: disabled = !(valid_username && valid_password)
</script>

<h2>Login</h2>
<form on:submit|preventDefault={() => {}}>
  <Input required label="Username" bind:value={body.username} bind:valid={valid_username}/>
  <Input required label="Password" password bind:value={body.password} bind:valid={valid_password}/>
  <Button click={login} bind:disabled>Submit</Button>
</form>
