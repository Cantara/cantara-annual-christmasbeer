<script>
  import Login from "../forms/Login.svelte";
  import Register from "../forms/Register.svelte";
  import RegisterBeer from "../forms/RegisterBeer.svelte";
  import RegisterAdmin from "../forms/RegisterPrivilege.svelte";
  import Score from "../forms/Score.svelte";
  import Scores from "../forms/Scores.svelte";

  import {bearer, token} from '../stores/token.js';

  let loggedInn = false;
  let isAdmin = false;
  const unsubscribe = token.subscribe(value => {
    loggedInn = value != null;
    if (loggedInn) {

      fetch('/account/admin', {
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
              .then(response => {
                if (response.ok) {
                  isAdmin = true;
                }
              });
    }
  });
</script>

<svelte:head>
	<title>Dashboard</title>
</svelte:head>

<div class="content flex">
  <h1>Cantara`s annual christmas beer tasting</h1>
  {#if (!loggedInn)}
  <div class="item">
    <h4 style="text-align: left;">Cantara holds a beer tasting every year to determine what beers are the best Norwegian christmas beer.</h4>
    <h4 style="text-align: left;">Time: Lørdag 12.11-2022, kl. 18:00 -></h4>
    <h4 style="text-align: left;">Place: 17 etg i Rebel, Universitetsgata 2</h4>
  </div>
  <div class="item">
    <h4 style="text-align: left;">NOTE!!! The physical tasting has completed, but we'll keep the voting open for remote votes until 20.12</h4>
    <p>There are a set of rules for the tasting and points are weighted based on if it is your first time or if you do certain things during the event.</p>
    <p>If you have brought a Norwegian beer that is missing from the list. Please ask this years selected newbie and have them add it.</p>
    <p><a href="https://wiki.cantara.no/display/puben/Puben+Home">Links to earlier years events</a></p>
  </div>
  <div class="new_line"/>
    <div class="item">
      <Login />
    </div>
    <div class="item">
      <Register />
    </div>
  {:else}
    <div class="item">
      <Score />
    </div>
    <div class="item">
      <RegisterBeer />
    </div>
    {#if (isAdmin)}
      <div class="item">
        <RegisterAdmin />
      </div>
    {/if}
  {/if}
  <div class="new_line"/>
  <div class="item">
    <Scores/>
  </div>
  <div class="new_line" style="padding-top: 1.5em;"/>
</div>

<style>
	h1 {
		color: var(--primary);
		font-size: 3em;
		font-weight: 350;
		margin: .5rem 0;
		font-style: italic;
		font-family: YouthTouch;
	}
  p {
    text-align: left;
  }
  hr {
    width: 100%;
    max-width: 100%;
    height: 0;
    max-height: 0;
    border: solid;
  	display: block;
		margin-top: .5em;
		margin-bottom: .5em;
		margin-left: auto;
		margin-right: auto;
		border-style: inset;
		border-width: 1px;
		color: rgba(0,0,0,.12);
	}
  .inline_content {
    display: flex;
    justify-content: center;
    align-content: center;
  }
  .content {
		position: relative;
		max-width: 1270px;
		margin-left: auto;
		margin-right: auto;
    padding-top: 1em;
    background: #fff;
  }
  .flex {
    display: flex;
    flex-flow: row wrap;
    justify-content: space-around;
    align-items: flex-start;
    align-content: space-around;
  }

  .item {
    flex: 0 0 45%;
  }
  @media screen and (max-width: 720px) {
    .item {
      flex: 0 0 80%;
    }
  }
  .min_item {
    flex: 0 0 20%;
  }
  .large_item {
    flex: 0 0 100%;
    width: 100%;
  }
  .new_line {
    flex: 0 0 100%;
  }
  .item_org {
    flex-basis: auto;
    flex-grow: 1;
    flex-shrink: 1;
  }
  .center {
    align-self: center;
  }
  .data {
    text-align:center;
    padding:4px .5em;
  }
</style>
