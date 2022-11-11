<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import Checkbox from "../components/Checkbox.svelte";

  import {bearer} from '../stores/token.js';
  import Select from "../components/Select.svelte";

  let accounts = []
  function getAccounts() {
    fetch('/account/accounts', {
      method: 'GET',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'omit',
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
              accounts = data.map(x => { return {name: x.firstname, id: x.id}})
            })
            .catch((error) => {
              console.log(error);
              //message.error(error);
            });
  }

  function register() {
    fetch('/account/privilege/' +  account.id, {
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

  let account = {}
  let body = {
    newbie: false,
    admin: false,
  }

  let valid_account = false;

  let disabled = false;

$: disabled = !(valid_account)// && body.gdpr)
</script>

<h2>Register</h2>
<Button click={getAccounts}>Refresh accounts</Button>
<form on:submit|preventDefault={() => {}}>
  <Select required label="Account" values={accounts} bind:value={account} bind:valid={valid_account} />
  <ul>
    <li><h3>Newbie?</h3></li>
    <li><Checkbox name="Newbie?" bind:checked={body.newbie} /></li>
  </ul>
  <ul>
    <li><h3>Admin?</h3></li>
    <li><Checkbox name="Admin?" bind:checked={body.admin} /></li>
  </ul>
  <Button click={register} bind:disabled>Submit</Button>
</form>

<style>
  ul {
    display: inline-flex;
  }
</style>