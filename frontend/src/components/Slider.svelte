<script>
  export let min_height = "4.4em";
  export let label = ""
  export let value;
  export let required = false;
  export let downarrow = false;
  export let min = 0;
  export let max = 100;
  export let valid = false;
  let visible = false;

$: {
  if (required && value < min) {
    valid = false;
  } else if (max !== -1 && value.length > max) {
    valid = false;
  } else {
    valid = true;
  }
}
</script>
<svelte:head>
  <script src="https://kit.fontawesome.com/673f197acb.js" crossorigin="anonymous"></script>
</svelte:head>

{#if label != ""}
  <fieldset class:invalid={!valid}>
    {#if downarrow}
      <legend><i class='fas fa-angle-down' style="font-size:1.25em"/>{(required ? ' * ' : '') + label}</legend>
    {:else}
      <legend>{(required ? '* ' : '') + label}</legend>
    {/if}
    <input type="range" bind:value={value} min={min} max={max} />
  </fieldset>
{/if}

<style>
div {
  border-style: solid !important;
  border-width: thin;
  border-color: gray;
  border-radius: 4px;
  height: 100%;
  overflow: auto;
  min-height: 1.3em;
}

input[type=”range”] {
  -webkit-appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #ff4500;
  cursor: ew-resize;
  box-shadow: 0 0 2px 0 #555;
  transition: background .3s ease-in-out;
}

input {
  width: 100%;
  margin: .5em 0 .75em 12px;
  padding: 0;
}

input:focus, textarea:focus {
  outline:none;
}

fieldset, div {
  text-align: left;
}

fieldset {
  margin-top: 1em;
  margin-bottom: .5em;
  border: 1px solid #ccc !important;
  border-style: solid;
  border-radius: .5em;
  padding: 0.01em 1em 0.01em 0em;
  display: flex;
  justify-content: center;
}
legend {
  white-space: nowrap;
  width: auto;
  margin-left: .75em;
}

input:disabled {
	color: #ccc;
}
.invalid {
  border-color: var(--error) !important;
}
.fas {
  font-family: "Font Awesome 5 Free";
}

input[type="range"] {
  -webkit-appearance: none;
  height: .5em;
  background: rgba(255, 255, 255, 0.6);
  background-image: linear-gradient(var(--slider-primary), var(--slider-primary));
  background-size: 50% 100%;
  background-repeat: no-repeat;
  border: 1px solid #ccc !important;
  border-radius: .75em;
}

/* Input Thumb */
input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  height: 1em;
  width: 1em;
  border-radius: 50%;
  background: var(--slider-primary);
  cursor: ew-resize;
  box-shadow: 0 0 2px 0 #555;
  transition: background .1s ease-in-out;
}

input[type="range"]::-moz-range-thumb {
  -webkit-appearance: none;
  height: 1em;
  width: 1em;
  border-radius: 50%;
  background: var(--slider-primary);
  cursor: ew-resize;
  box-shadow: 0 0 2px 0 #555;
  transition: background .1s ease-in-out;
}

input[type="range"]::-ms-thumb {
  -webkit-appearance: none;
  height: 1em;
  width: 1em;
  border-radius: 50%;
  background: var(--slider-primary);
  cursor: ew-resize;
  box-shadow: 0 0 2px 0 #555;
  transition: background .1s ease-in-out;
}

input[type="range"]::-webkit-slider-thumb:hover {
  background: var(--slider-secoundary);
}

input[type="range"]::-moz-range-thumb:hover {
  background: var(--slider-secoundary);
}

input[type="range"]::-ms-thumb:hover {
  background: var(--slider-secoundary);
}

/* Input Track */
input[type=range]::-webkit-slider-runnable-track  {
  -webkit-appearance: none;
  box-shadow: none;
  border: none;
  background: transparent;
}

input[type=range]::-moz-range-track {
  -webkit-appearance: none;
  box-shadow: none;
  border: none;
  background: transparent;
}

input[type="range"]::-ms-track {
  -webkit-appearance: none;
  box-shadow: none;
  border: none;
  background: transparent;
}
</style>
