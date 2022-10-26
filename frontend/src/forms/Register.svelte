<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import Checkbox from "../components/Checkbox.svelte";

  function putDatabase() {
    fetch('/christmasbeer/register', {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'omit',
      body: JSON.stringify(body),
      headers: {
        'Authorization': 'Basic ' + btoa(user.name + ":" + user.password),
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
    })
    .catch((error) => {
      console.log(error);
      //message.error(error);
    });
  }

  let body = {
    username: "",
    password: "",
    gdpr: false,
  }

  let valid_username = false;
  let valid_password = false;

  let disabled = false;

  export let user = {
    name: "",
    password: "",
  }

$: disabled = !(valid_username && valid_password && body.gdpr)
</script>

<h2>Register</h2>
<form on:submit|preventDefault={() => {}}>
  <Input required label="Username" bind:value={body.username} bind:valid={valid_username} min=3/>
  <Input required label="Password" password bind:value={body.password} bind:valid={valid_password} min=6 max=64/>
  <div style="display: inline-flex">
    <Checkbox required bind:checked={body.gdpr}/>
    <p>Check to consent to storage of provided values (GDPRish)</p>
  </div>
  <Button click={putDatabase} bind:disabled>Submit</Button>
</form>
