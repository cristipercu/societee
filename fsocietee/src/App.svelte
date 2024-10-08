<script>
    import Navbar from "../layout/Navbar.svelte";
    import Home from "../pages/Home.svelte"
    import { Router, Link, Route } from "svelte-routing";
</script>

<Router>
<main>
  <Navbar />
  <div class="container mt-3 d-flex justify-content-center">
  <Route path="/" component={Home} />
  </div>
</main>
</Router>

<style>
</style>
