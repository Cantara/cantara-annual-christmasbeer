<script>
  import Box from "../components/Box.svelte";
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import Select from "../components/Select.svelte";
  import Slider from "../components/Slider.svelte";
  import { onMount, tick } from 'svelte';
  import {bearer} from "../stores/token";

  let beers = [];
  let socket;
  let board;

  const scrollToBottom = async () => {
    await tick();
    board.scrollTop = board.scrollHeight;
  };

  onMount(() => {
    let prot = 'wss'

    if (location.hostname === "localhost" && location.protocol === "http:") {
      prot = 'ws'
    }
    socket = new WebSocket(prot + '://'+ location.host + '/beer');

    socket.onopen = () => {
      console.log('socket connected');
    };

    socket.onmessage = (event) => {
      if(!event.data) { return }
      let b = JSON.parse(event.data);
      beers = [...beers, {name: b.name, extras: [b.brand, b.brew_year, b.abv+"%"], id: b.brand + "_" + b.name + "_" + b.brew_year}];
      //scrollToBottom();
      //socket.close();
      console.log(beers)
    };
  });

  function register() {
    fetch('/score/' + new Date().getFullYear() + "/" + encodeURI(beer.id), {
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
                return
              }
              console.log(data)
            })
            .catch((error) => {
              console.log(error);
            });
  }

  let beer = {}
  let body = {
    rating: 0,
    comment: "",
  }
  let valid_beer = false;
  let valid_rating = false;
  let valid_comment = false;

  let disabled = false;

  $: disabled = !(valid_beer && valid_rating && valid_comment)
</script>

<h2>Register beer score</h2>
<p>If your beer is missing, please ask this years newbie to help you</p>
<form on:submit|preventDefault={() => {}}>
  <Select required label="Beer" values={beers} bind:value={beer} bind:valid={valid_beer}/>
  <Slider required min=1 max=6 label="Rating" bind:value={body.rating} bind:valid={valid_rating} number/>
  <Input required multiline autogrow label="Comment" bind:value={body.comment} bind:valid={valid_comment}/>
  <Button click={register} bind:disabled>Submit</Button>
</form>
