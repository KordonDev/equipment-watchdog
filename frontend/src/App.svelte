<script lang="ts">
  import Router, { link, replace } from "svelte-spa-router";
  import Login from "./views/security/Login.svelte";
  import { routes } from "./routes";
  import Register from "./views/security/Register.svelte";
  import MemberOverview from "./views/member/MemberOverview.svelte";
  import { logout } from "./views/security/security.service";
  import AddMember from "./views/member/AddMember.svelte";
  import MemberDetail from "./views/member/MemberDetail.svelte";
  import Notification from "./components/Notification/Notification.svelte";

  if (window.location.pathname === "/" && window.location.hash === "") {
    replace(routes.MemberOverview.link);
  }
</script>

<main>
  <nav>
    <a href={routes.Login.link} use:link>Einloggen</a>
    <a href={routes.MemberOverview.link} use:link>Übersicht</a>
    <a href={routes.AddMember.link} use:link>Mitglied hinzufügen</a>
    <button on:click={logout}>Ausloggen</button>
  </nav>

  <Router
    routes={{
      [routes.Register.path]: Register,
      [routes.Login.path]: Login,
      [routes.MemberOverview.path]: MemberOverview,
      [routes.AddMember.path]: AddMember,
      [routes.MemberDetail.path]: MemberDetail,
    }}
  />
</main>

<Notification />

<style>
</style>
