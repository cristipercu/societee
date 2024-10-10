<script>
import { API_BASE_URL, ROOM_ID_LENGTH } from '../scripts/config.js';
import { makeid } from '../scripts/utils.js'
import { Link } from "svelte-routing";

function handleClick() {
  roomid = makeid(ROOM_ID_LENGTH)
  console.log(roomid)
}

function handleRegister() {

}

</script>

<h1>Have fun!</h1>

<!-- <button type="button" class="btn btn-dark" on:click={handleClick}>create room</button> -->
<Link to="/register">
<button type="button" class="btn btn-dark" on:click={handleRegister}>Register</button>
</Link>
