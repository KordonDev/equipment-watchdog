<script lang="ts">
  import { link } from "svelte-spa-router";
  import { logout } from "../../views/security/security.service";
  import { routes } from "../../routes";
  import { EquipmentType } from "../../views/equipment/equipment.service";
</script>

<nav>
  <a href={routes.MemberOverview.link} use:link>Übersicht</a>
  <a href={routes.AddMember.link} use:link>Mitglied hinzufügen</a>
  <a href={routes.Users.link} use:link>User</a>
  <a href={`${routes.EquipmentType.link}${EquipmentType.Jacket}`} use:link>
    Ausrüstung
  </a>
  <button on:click={logout}>Ausloggen</button>
</nav>
