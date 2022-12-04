<script lang="ts">
  import { link } from "svelte-spa-router";
  import { logout } from "../../views/security/security.service";
  import { routes } from "../../routes";
</script>

<nav>
  <a href={routes.Login.link} use:link>Einloggen</a>
  <a href={routes.MemberOverview.link} use:link>Übersicht</a>
  <a href={routes.AddMember.link} use:link>Mitglied hinzufügen</a>
  <a href={routes.Users.link} use:link>User</a>
  <button on:click={logout}>Ausloggen</button>
</nav>
