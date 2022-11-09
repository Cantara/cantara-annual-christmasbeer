<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";

  import {bearer} from '../stores/token.js';

  function register() {
    fetch('/account/admin/' +  body.username, {
      method: 'PUT',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'omit',
      body: JSON.stringify(body),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': bearer(),
      },
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        console.log(data.error);
        //message.error(data.error);
        return
      }
    })
    .catch((error) => {
      console.log(error);
      //message.error(error);
    });
  }

  let body = {
    username: "",
  }

  let valid_username = false;

  let disabled = false;

$: disabled = !(valid_username)// && body.gdpr)
</script>

<h2>Register</h2>
<form on:submit|preventDefault={() => {}}>
  <Input required label="Username" bind:value={body.username} bind:valid={valid_username} min=3 />
  <Button click={register} bind:disabled>Submit</Button>
</form>
