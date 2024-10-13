<script>
  import { API_URL } from '../scripts/config.js';
  import { navigate } from "svelte-routing";

  let username = '';
  let email = '';
  let password = '';
  let message = '';

async function handleSubmit() {
    try {
      const response = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password }),
      });

      if (response.ok) {
        setTimeout(() => {
          navigate("/", {replace: true});
        }, 1000);
      } else if (response.status === 409){
        message = "email or user already exists"
      }
    } catch (error) {
      message = 'Somethin went wrong in the back!';
    }
  }
</script>
<div class="row">
  <h1>register</h1>
  {#if message}
    <p>{message}</p>
  {/if}
</div>

<div class="row">
<form on:submit|preventDefault={handleSubmit}>
  <div class="mb-3">
    <label for="username" class="form-label">username</label>
    <input type="text" id="username" class="form-control" bind:value={username} required />
  </div>
  <div class="mb-3">
    <label for="email" class="form-label">email</label>
    <input type="email" id="email" class="form-control" bind:value={email} required />
  </div>
  <div class="mb-3">
    <label for="password" class="form-label">password</label>
    <input type="password" id="password" class="form-control" bind:value={password} required />
  </div>
  <button type="submit" class="btn btn-dark">register</button>
</form>
</div>
