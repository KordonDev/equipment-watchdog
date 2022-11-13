<script lang="ts">
  import Router from "svelte-spa-router";
  import Login from "./views/security/Login.svelte";
  import { routes } from "./routes";
  import Register from "./views/security/Register.svelte";
  import MemberOverview from "./views/member/MemberOverview.svelte";
  import { logout } from "./views/security/security.service";

  console.log(window.location.pathname)
  if (window.location.pathname === "/") {
    window.location.assign(`/#${routes.MemberOverview}`)
  }
</script>

<main>
  <nav>
    <a href={`/#${routes.Login}`}>Einloggen</a>
    <a href={`/#${routes.Register}`}>Registrieren</a>
    <a href={`/#${routes.MemberOverview}`}>Übersicht</a>
    <button on:click={logout}>Ausloggen</button>
  </nav>

  <Router
    routes={{
      [routes.Register]: Register,
      [routes.Login]: Login,
      [routes.MemberOverview]: MemberOverview,
    }}
  />
</main>

<style>
</style>
