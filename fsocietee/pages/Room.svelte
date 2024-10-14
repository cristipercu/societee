<script>
import Teams from "../components/Teams.svelte";
import { API_URL, WS_URL } from "../scripts/config";
import { AmITheAdmin, connectToWebsocket, getInternalApiHeaders, sendMessage } from "../scripts/utils.js"
import { onMount, onDestroy } from 'svelte';

let roomCode = localStorage.getItem('roomCode');
let roomID = localStorage.getItem('roomID');
let adminFlag = AmITheAdmin();
let displayName = localStorage.getItem('displayName');
let number_of_words_per_player = 10;
let time_per_player = 60;
let ws_url = `${WS_URL}/ws/rooms/${roomID}?token=${localStorage.getItem('token')}`;
let html_message;
let ws;
let players = [];


function number_of_words_on_imput(event) {
  number_of_words_per_player = event.target.value
};

function time_per_player_on_imput(event) {
  time_per_player = event.target.value
};

async function getRoomMembers() {
  let url = `${API_URL}/rooms/members/${roomID}`;
  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: getInternalApiHeaders(true)
    });
    
    const data = await response.json()
    if (!response.ok) {
      html_message = data.error;    
    }

    if (response.ok) {
      console.log(data);
      return data
    }
  } catch (error) {
    html_message = error; 
  }
} 

onMount(() => {
function onOpen(ws) {
  sendMessage(ws, displayName, 'join_room', 'here i am')
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
      players = getRoomMembers();
      if (data.user != displayName) {
        html_message = `${data.user} joined the room`;
        appendAlert(html_message, 'dark');
      }
      console.log(players);
      break
    case 'leave':
      appendAlert(`${data.user} left the building`, 'danger')
  }
}

ws = connectToWebsocket(ws_url, onOpen = onOpen, onMessage=onMessage);
})

async function removeMeFromRoom() {
  let url = `${API_URL}/rooms/leave`;
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: getInternalApiHeaders(true),
      body: JSON.stringify({"room_id": roomID, "player_name": displayName})
    });

    const data = await response.json();
    if (!response.ok) {
      html_message = data.error
    };
    if (response.ok) {
      console.log(data)
    };
  } catch (error) {
    console.log(error)
  }
}

onDestroy(() => {
  if (ws) {
    removeMeFromRoom()
    sendMessage(ws, displayName, 'leave', 'ciao');
    ws.close()
    const apiUrl =`${API_URL}/rooms/leave`;
    const data =  {"room_id": roomID, "player_name": displayName}

    fetch(apiUrl, {
      method: 'POST',
      headers: getInternalApiHeaders(true), 
      body: JSON.stringify(data)
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Eroare  la trimiterea requestului');
      }
      return response.json();
    })
    .then(data => {
      console.log('Răspuns de la API:', data);
    })
    .catch(error => {
      console.error('Eroare:', error);
    });
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

