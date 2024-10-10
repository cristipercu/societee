<script>
    import Navbar from "../layout/Navbar.svelte";
    import Home from "../pages/Home.svelte"
    import Login from "../pages/Login.svelte";
    import Register from "../pages/Register.svelte"
    import { Router, Link, Route } from "svelte-routing";



</script>

<Router>
<main>
  <Navbar />
  <div class="container mt-3 d-flex flex-column align-items-center justify-content-center h-100">
      <Route path="/" component={Home} />
      <Route path="/register" component={Register} />
      <Route path="/login" component={Login} />
  </div>
</main>
</Router>

<style>
</style>
