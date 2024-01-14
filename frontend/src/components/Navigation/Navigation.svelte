<script lang="ts">
  import { link } from "svelte-spa-router";
  import { logout } from "../../views/security/security.service";
  import { routes } from "../../routes";
  import { EquipmentType } from "../../views/equipment/equipment.service";
  import { Button } from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import type { User } from "src/views/user/user.service";

  let me: User = null;
  currentUser.subscribe((u) => (me = u));
</script>

<nav>
  <a href={routes.MemberOverview.link} use:link>Mitglieder</a> |
  <a href={`${routes.EquipmentType.link}${EquipmentType.Jacket}`} use:link>
    Ausrüstung
  </a>
  |
  <a href={routes.Orders.link} use:link>Bestellungen</a>
  {#if me && me.isAdmin}
    |
    <a href={routes.Users.link} use:link>User</a>
  {/if}
  <Button size="xs" color="purple" class="ml-5" on:click={logout}>Ausloggen</Button>
</nav>

