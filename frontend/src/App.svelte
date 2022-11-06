<script lang="ts">
  import Router from "svelte-spa-router";
  import Login from "./views/security/Login.svelte";
  import { routes } from "./routes";
  import Register from "./views/security/Register.svelte";
  import Home from "./views/Home.svelte";
  import MemberOverview from "./views/member/MemberOverview.svelte";
</script>

<main>
  <nav>
    <a href={routes.Home}>Home</a>
    <a href={`/#${routes.Login}`}>Einloggen</a>
    <a href={`/#${routes.Register}`}>Registrieren</a>
    <a href={`/#${routes.MemberOverview}`}>Übersicht</a>
  </nav>

  <Router
    routes={{
      [routes.Register]: Register,
      [routes.Login]: Login,
      [routes.Home]: Home,
      [routes.MemberOverview]: MemberOverview,
    }}
  />
</main>

<style>
</style>
