<script>
  import Box from "../components/Box.svelte";
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import Select from "../components/Select.svelte";
  import Slider from "../components/Slider.svelte";
  import { onMount, tick } from 'svelte';
  import {bearer} from "../stores/token";

  let lastScores = [];
  let mostRated = [];
  let highestRated = [];
  let socket;
  let board;
  const scores = new Map();

  const scrollToBottom = async () => {
    await tick();
    board.scrollTop = board.scrollHeight;
  };

  onMount(() => {
    let prot = 'wss'

    if (location.hostname === "localhost" && location.protocol === "http:") {
      prot = 'ws'
    }
    socket = new WebSocket(prot + '://'+ location.host + '/score');

    socket.onopen = () => {
      console.log('socket connected');
    };

    socket.onmessage = (event) => {
      if(!event.data) { return }
      let s = JSON.parse(event.data);
      lastScores = [s, ...lastScores];
      if (lastScores.length > 5) {
        lastScores = lastScores.slice(0, 5);
      }
      let beer_id = s.beer.brand + "_" + s.beer.name + "_" + s.beer.brew_year;
      //let beer_id = s.year + "_" + s.beer.brand + "_" + s.beer.name + "_" + s.beer.brew_year
      if (scores.get(beer_id)) {
          scores.set(beer_id, [...scores.get(beer_id), s]);
      } else {
          scores.set(beer_id, [s])
      }
      let hr = [];
      let mr = [];
      scores.forEach(function (v, key, map) {
          console.log("v:", v[0])
          if (hr.length < 5) {
              hr = [...hr, {beer: v[0].beer, avg: v.reduce((p,c) => p + c.rating, 0) / v.length}]
              mr = [...mr, {beer: v[0].beer, num: v.length}]
          }
      });
      highestRated = hr.sort((a,b)=> {
          if (a.avg < b.avg) {
              return 1;
          } else if (a.avg > b.avg) {
              return -1;
          }
          return 0;
      }).slice(0, 5);
      mostRated = mr.sort((a,b)=> {
          if (a.num < b.num) {
              return 1;
          } else if (a.num > b.num) {
              return -1;
          }
          return 0;
      }).slice(0, 5);
      //scrollToBottom();
      //socket.close();
      //console.log(highestRated);
      //console.log(mostRated);
      //console.log(lastScores);
      //console.log(scores);
    };
  });
</script>

<h2>Five newest ratings</h2>
{#if (lastScores && Array.isArray(lastScores))}
    <ol>
        {#each lastScores as score}
            <li>{score.beer.brand} {score.beer.name} {score.beer.brew_year}: {parseInt(score.rating)} {score.scorer}<br> {score.comment}</li>
        {/each}
    </ol>
{/if}
<h2>Five beers with highest avg rating</h2>
{#if (highestRated && Array.isArray(highestRated))}
    <ol>
        {#each highestRated as score}
            <li>{score.beer.brand} {score.beer.name} {score.beer.brew_year}: avg {parseInt(score.avg)}</li>
        {/each}
    </ol>
{/if}
<h2>Five most rated beers</h2>
{#if (mostRated && Array.isArray(mostRated))}
    <ol>
        {#each mostRated as score}
            <li>{score.beer.brand} {score.beer.name} {score.beer.brew_year}: num {parseInt(score.num)}</li>
        {/each}
    </ol>
{/if}

<style>
    ol {
        text-align: left;
    }
</style>