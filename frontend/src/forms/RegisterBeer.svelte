<script>
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import {token} from "../stores/token";

  function register() {
    fetch('/beer/' +  body.brand + "_" + body.brew_year + "_" + body.name, {
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
                return
              }
              console.log(data)
            })
            .catch((error) => {
              console.log(error);
            });
  }


  let body = {
    name:"",
    brand:"",
    brew_year: new Date().getFullYear(),
    abv: 0,
  }

  let valid_name = false;
  let valid_brand = false;
  let valid_brew_year = false;
  let valid_abv = false;

  let disabled;

$: disabled = !(valid_name && valid_brand && valid_brew_year && valid_abv)
</script>

<h2>Register beer</h2>
<p>If you think you should have access to register your own beer, please contact this years designated newbie.</p>
<form on:submit|preventDefault={() => {}}>
  <Input required label="Name" bind:value={body.name} bind:valid={valid_name}/>
  <Input required label="Brand" bind:value={body.brand} bind:valid={valid_brand}/>
  <Input required number min=1980 max={new Date().getFullYear()}  label="Brew year" bind:value={body.brew_year} bind:valid={valid_brew_year}/>
  <Input required float min=0 max=98 label="ABV%" bind:value={body.abv} bind:valid={valid_abv}/>
  <Button click={register} bind:disabled>Submit</Button>
</form>