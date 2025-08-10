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

  const isActiveUrl = (url: string) => window.location.hash.startsWith(`#/${url}`) ? 'active': '';
</script>

<nav class="nav">
  <div>
    <a href={routes.MemberOverview.link} use:link class={isActiveUrl('member')}>Mitglieder</a>&nbsp;|&nbsp;
  </div>
  <div>
    <a href={`${routes.EquipmentType.link}${EquipmentType.Jacket}`} use:link class={isActiveUrl('equipment')}>
      Ausrüstung
    </a>&nbsp;|&nbsp;
  </div>
  <div>
    <a href={routes.Orders.link} use:link class={isActiveUrl('orders')}>Bestellungen</a>&nbsp;|&nbsp;
  </div>
  <div>
    <a href={routes.Users.link} use:link class={isActiveUrl('users')}>User</a>
  </div>
  <Button size="xs" color="red" class="ml-5" on:click={logout}>
    <svg class="w-4 h-4 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H8m12 0-4 4m4-4-4-4M9 4H7a3 3 0 0 0-3 3v10a3 3 0 0 0 3 3h2"/>
    </svg>
  </Button>
</nav>

<style>
  .nav {
    margin-left: -24px;
    margin-right: -24px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .nav .active {
    background-color: rgb(4 108 78);
    color: white;
    padding-right: 4px;
    padding-left: 4px;
    border-radius: 6px;
  }
</style>