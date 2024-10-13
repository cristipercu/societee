<script>
  import { API_URL } from "../scripts/config";
  import { handleLogin } from "../scripts/utils" 

  let email = '';
  let password = '';
  let message = '';
  let url = `${API_URL}/login`

async function handleSubmit() {
    await handleLogin(email, password, message, url);
}
</script>

<h1>log in</h1>
{#if message}
  <p>{message}</p>
{/if}
<div class="row">
  <form on:submit|preventDefault={handleSubmit}>
    <div class="mb-3">
      <label class="form-label" for="email">email</label>
      <input class="form-control" type="email" id="email" bind:value={email} required />
    </div>
    <div class="mb-3">
      <label class="form-label" for="password">password</label>   
      <input class="form-control" type="password" id="password" bind:value={password} required />
    </div>
    <button class="btn btn-dark" type="submit">login</button>
  </form>
</div>
