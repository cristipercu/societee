<script>
    import Navbar from "../layout/Navbar.svelte";
    import Home from "../pages/Home.svelte"
    import Register from "../pages/Register.svelte"
    import { Router, Link, Route } from "svelte-routing";


</script>

<Router>
<main>
  <Navbar />
  <div class="container mt-3 d-flex flex-column align-items-center justify-content-center h-100">
      <Route path="/" component={Home} />
      <Route path="/register" component={Register} />
  </div>
</main>
</Router>

<style>
</style>
