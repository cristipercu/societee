<script>
import { navigate } from "svelte-routing";
import { API_URL } from "../scripts/config"

let room_code = '';
let password = '';
let message;
let showPassword = false
let roomID;
let displayName = localStorage.getItem('displayName');

async function handleSubmit() {
  let url = `${API_URL}/roomCode`
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: getInternalApiHeaders(true),
      body: JSON.stringify({room_code, password}),
    });

    const data = await response.json();
    if (!response.ok) {
      console.log('response not ok')
      room_code = '';
      password = '';
      message = data.error;
    }

    if (response.ok) {
      roomID = data.id;
      localStorage.setItem('roomCode', data.room_code);
      localStorage.setItem('roomID', data.id);
      localStorage.setItem('roomAdmin', data.owner_name);
      addMeToRoom();
      navigate('/room');
    }
  } catch (error) {
    message = 'error connecting to server, call CP';
  }
}; 

async function addMeToRoom() {
  let url = `${API_URL}/rooms/addme`;
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: getInternalApiHeaders(true),
      body: JSON.stringify({"room_id": roomID, "player_name": displayName})
    });
    
    const data = await response.json();
    if (!response.ok) {
      console.log('adding user to room not ok')
    };

    if (response.ok) {
      console.log(data);
    };

  } catch (error) {
  message = 'error connecting to server, call CP';
  }
}

function togglePasswordVisibility() {
    showPassword = !showPassword;
    password = '';
  }


</script>


<div class="card mt-3">
  <div class="card-body">
  <h3> or join room </h3>
  <form on:submit|preventDefault={handleSubmit}>
    <div class="mb-3">
      <label class="form-label" for="roomCode">room code</label>
      <input class="form-control" type="text" id="roomCode" bind:value={room_code} required />
    </div>
    {#if showPassword}
      <label class="form-label" for="password">password</label>
      <input class="form-control" type="text" id="password" bind:value={password} required />
    {/if}
    <button type="button" class="btn btn-dark mt-3" on:click={togglePasswordVisibility}>
    {showPassword ? 'remove password' : 'add password'}
    </button>
    <div class="mt-3">
      <button class="btn btn-dark" type="submit">join</button>
    </div>
    {#if message}
    <p>{message}</p> 
    {/if}
  </form>
</div>
</div>
