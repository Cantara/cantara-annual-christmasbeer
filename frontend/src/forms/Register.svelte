<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";

  import { token } from '../stores/token.js';

  function register() {
    fetch('/account/' +  body.username, {
      method: 'PUT',
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
      if (data.error) {
        console.log(data.error);
        //message.error(data.error);
        return
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
    //gdpr: false,
  }

  let valid_username = false;
  let valid_password = false;

  let disabled = false;

$: disabled = !(valid_username && valid_password)// && body.gdpr)
</script>

<h2>Register</h2>
<form on:submit|preventDefault={() => {}}>
  <Input required label="Username" bind:value={body.username} bind:valid={valid_username} min=3 />
  <Input required label="Password" password bind:value={body.password} bind:valid={valid_password} min=6 max=64 />
  <!--
  <div style="display: inline-flex">
    <Checkbox required bind:checked={body.gdpr}/>
    <p>Check to consent to storage of provided values (GDPRish)</p>
  </div>
  -->
  <Button click={register} bind:disabled>Submit</Button>
</form>
