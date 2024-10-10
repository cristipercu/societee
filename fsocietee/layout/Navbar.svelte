<script>
  import { Link } from "svelte-routing";
</script>

<style>
</style>

<nav class="navbar navbar-expand-lg bg-body-tertiary">
  <div class="container-fluid">
    <Link to="/">
    <span class="navbar-brand">societee</span>
    </Link>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarText" aria-controls="navbarText" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarText">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <Link to="/rules">
          <span class="nav-link active" aria-current="page" >rules</span>
          </Link>
        </li>
      </ul>
      <span class="navbar-text">
        have fun!
      </span>
    </div>
  </div>
</nav>
