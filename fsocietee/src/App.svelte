<script>
    import Navbar from "../layout/Navbar.svelte";
    import CreateRoom from "../pages/CreateRoom.svelte";
    import Home from "../pages/Home.svelte"
    import Login from "../pages/Login.svelte";
    import Room from "../pages/Room.svelte";
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
      <Route path="/createRoom" component={CreateRoom} />
      <Route path="/room" component={Room} />
  </div>
</main>
</Router>

<style>
</style>
