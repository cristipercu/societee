<script>
import { Link } from "svelte-routing";
import {isTokenExpired, clearGameLogicFromLocalStorage} from "../scripts/utils"
    import JoinRoom from "../layout/JoinRoom.svelte";

clearGameLogicFromLocalStorage()
let displayName = localStorage.getItem('displayName')
let token = localStorage.getItem('token')
let tokenExpired = isTokenExpired(token)

</script>


{#if tokenExpired}
<Link to="/register">
<button type="button" class="btn btn-dark mt-3">register</button>
</Link>

<Link to="/login">
<button type="button" class="btn btn-dark mt-3">login</button>
</Link>
{:else}
<h1>have fun {displayName}!</h1>


<Link to="/createRoom">
  <button type="button" class="btn btn-dark mt-3">create room</button>
</Link>

<JoinRoom />

{/if}
