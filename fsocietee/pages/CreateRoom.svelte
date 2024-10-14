<script>
import { navigate } from "svelte-routing";
import { API_URL } from '../scripts/config.js';
import { getInternalApiHeaders } from "../scripts/utils.js"
import { makeid } from '../scripts/utils';

let max_members = 4;
let password;
let showPassword = false;
let owner_name = localStorage.getItem('displayName');
let message;
let room_code = makeid(5)
let displayName = localStorage.getItem('displayName');

function togglePasswordVisibility() {
    showPassword = !showPassword;
    password = '';
  }

function handleInput(event) {
    max_members = event.target.value;
}


async function handleSubmit() {
    try {
      const response = await fetch(`${API_URL}/rooms`, {
        method: 'POST',
        headers: getInternalApiHeaders(true),
        body: JSON.stringify({room_code, password, owner_name, max_members }),
      });

      const results = await response.json();
      if (response.ok) {
        console.log(results);
        localStorage.setItem('roomCode', room_code);
        localStorage.setItem('roomID', results['room_id']);
        localStorage.setItem('roomAdmin', displayName);
        navigate("/room", {replace: true});
      }

      if (!response.ok) {
        message = results.error
      }
    } catch (error) {
      message = 'Somethin went wrong in the back!';
      console.log(error.message);
    }
}
</script>

<style>
input[type=range]::-webkit-slider-thumb {
  -webkit-appearance: none; 
  appearance: none;
  width: 15px; 
  height: 15px;
  background: black; 
  cursor: pointer; 
  border-radius: 50%; 
}
</style>

<div class="row">
  <h3>room settings</h3>
  <form on:submit|preventDefault={handleSubmit}>
    <div class="mb-3">
      <label for="players" class="form-label">number of players {max_members}</label>
      <input type="range" min="4" max="10" step="2" class="form-range" id="players" bind:value={max_members} on:input={handleInput}>

      {#if showPassword}
        <label class="form-label" for="password">password</label>
        <input class="form-control" type="text" id="password" bind:value={password} required />
      {/if}
      <button type="button" class="btn btn-dark mt-3" on:click={togglePasswordVisibility}>
      {showPassword ? 'remove password' : 'add password'}
      </button>
    </div>
    {#if message}
      <p>{message}</p>
    {/if}
 
    <button class="btn btn-dark" type="submit">create room</button>
  </form>
</div>
