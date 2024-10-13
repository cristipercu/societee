<script>
    import Teams from "../components/Teams.svelte";
import { API_URL, WS_URL } from "../scripts/config";
import { AmITheAdmin, connectToWebsocket, sendMessage } from "../scripts/utils.js"
import { onMount, onDestroy } from 'svelte';

let roomCode = localStorage.getItem('roomCode');
let roomID = localStorage.getItem('roomID');
let adminFlag = AmITheAdmin();
let displayName = localStorage.getItem('displayName');
let number_of_words_per_player = 10;
let time_per_player = 60;
let url = `${API_URL}/games/settings`;
let ws_url = `${WS_URL}/ws/rooms/${roomID}`;
let html_message;
let ws;
let players = [];


function number_of_words_on_imput(event) {
  number_of_words_per_player = event.target.value
};

function time_per_player_on_imput(event) {
  time_per_player = event.target.value
};

onMount(() => {
function onOpen(ws) {
  sendMessage(ws, displayName, 'join_room', 'sunt aici')
}
const alertPlaceholder = document.getElementById('liveAlertPlaceholder')
const appendAlert = (message, type) => {
  const wrapper = document.createElement('div')
  wrapper.innerHTML = [
          `<div class="alert alert-${type} alert-dismissible" role="alert">`,
          `   <div>${message}</div>`,
          '   <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>',
          '</div>'
        ].join('')
  alertPlaceholder.append(wrapper)
}



function onMessage(message) {
  let data = JSON.parse(message)
  switch (data.type) {
    case 'join_room':
      if (data.user != displayName) {
        html_message = `${data.user} joined the room`;
        players = [...players, data.user];
        appendAlert(html_message, 'dark');
      }
      break
    case 'leave':
      appendAlert(`${data.user} left the building`, 'danger')
  }
}

ws = connectToWebsocket(ws_url, onOpen = onOpen, onMessage=onMessage);
})

onDestroy(() => {
    console.log('Componenta va fi demontată!');
    if (ws) {
      sendMessage(ws, displayName, 'leave', 'ciao');
      ws.close()
  }
});



</script>

<h3>Welcome to {roomCode}</h3>
{#if adminFlag}
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">game settings</h5>
      <div class="mt-3">
        <label for="number_of_words" class="form-label"> number of words per player: {number_of_words_per_player} </label>
        <input type="range" min="8" max="20" step="1" class="form-range" id="number_of_words" bind:value={number_of_words_per_player} on:input={number_of_words_on_imput}>
      </div>
      <div class="mt-3">
        <label for="time_per_player" class="form-label"> time per player: {time_per_player} seconds</label>
        <input type="range" min="40" max="120" step="5" class="form-range" id="time_per_player" bind:value={time_per_player} on:input={time_per_player_on_imput}>
      </div>
    </div>
  </div>

<Teams/>

{/if}



<div id="liveAlertPlaceholder" class="mt-3"></div>

