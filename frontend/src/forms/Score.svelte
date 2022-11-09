<script>
  import Box from "../components/Box.svelte";
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import Select from "../components/Select.svelte";
  import Slider from "../components/Slider.svelte";
  import { onMount, tick } from 'svelte';

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
      beers = [...beers, {name: b.name, extras: [b.brand, b.brew_year, b.abv+"%"]}];
      //scrollToBottom();
      //socket.close();
      console.log(beers)
    };
  });

  function putService() {
    let bodyT = body
    bodyT.service.health.service_type = bodyT.service.health.service_type.name
    fetch('/christmasbeer/score/'+body.beerid, {
      method: 'PUT',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'omit',
      body: JSON.stringify(bodyT),
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
    beerid:"",
    service: {
      elb_listener_arn: "",
      elb_securitygroup_id: "",
      port: 0,
      path: "",
      artifact_id: "",
      health: {
        service_name: "",
        service_tag: "",
        service_type: {
          name: "",
        },
      },
      local_override_properties: "",
      semantic_update_service_properties: "",
    },
    key: "",
  }
  let scope = -1;

  let valid_key = false;
  let valid_scope = false;
  let valid_server = false;
  let valid_loadbalancer = false;
  let valid_port = false;
  let valid_path = false;
  let valid_artifact = false;
  let valid_health_name = false;
  let valid_health_tag = false;
  let valid_health_type = false;
  let valid_local = false;
  let valid_semantic = false;

  let disabled = false;

  export let user = {
    name: "",
    password: "",
  }
  export let loadbalancers = [];
  let loadbalancer = {};
  let loadbalancersDropdown = [];
  let health_service_types = [
    {
      name: "A2A",
      extras: [""],
    }
  ];

$: {
  loadbalancersDropdown = loadbalancers.map(value => ({
    name: value.dns_name,
    extras: value.paths,
    arn: value.arn,
    listener_arn: value.listener_arn,
    security_group: value.security_group
  }))
}
$: body.service.elb_listener_arn = loadbalancer.listener_arn
$: body.service.elb_securitygroup_id = loadbalancer.security_group

$: disabled = !(valid_scope && valid_server && (valid_loadbalancer || loadbalancer != {}) && valid_key && valid_port && valid_path && valid_artifact && valid_health_name && valid_health_tag && (valid_health_type || body.service.health.service_type.name != "") && valid_local && valid_semantic && valid_key)
</script>

<h2>Register beer score</h2>
<p>If your beer is missing, please ask this years newbie to help you</p>
<form on:submit|preventDefault={() => {}}>
  <Select required label="Beer" values={beers} bind:value={body.service.health.service_type} bind:valid={valid_health_type}/>
  <Slider required label="Rating" bind:value={scope} bind:valid={valid_scope} number/>
  <Input required multiline autogrow label="Comment" bind:value={body.key} bind:valid={valid_key}/>
  <Button click={putService} bind:disabled>Submit</Button>
</form>
