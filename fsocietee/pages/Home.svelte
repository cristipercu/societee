<script>
import { Link } from "svelte-routing";
import {isTokenExpired} from "../scripts/utils"

let displayName = localStorage.getItem('displayName')
let token = localStorage.getItem('token')
let tokenExpired = isTokenExpired(token)

console.log(token);
console.log(tokenExpired);

</script>


{#if tokenExpired}
<Link to="/register">
<button type="button" class="btn btn-dark">register</button>
</Link>

<Link to="/login">
<button type="button" class="btn btn-dark">login</button>
</Link>
{:else}
<h1>have fun {displayName}!</h1>

{/if}
