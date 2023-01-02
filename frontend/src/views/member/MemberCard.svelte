<script lang="ts">
  import { Card } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import type { Member } from "./member.service";
  import Equipments from "../../components/Equipment/Equipments.svelte";

  export let member: Member;
  export let columns: number;
</script>

<Card class="m-4" size="xl">
  <a href={`${routes.MemberDetail.link}${member.id}`} use:link>
    {`${member.id} ${member.name}`}
  </a>
  <Equipments {member} {columns} />
</Card>
