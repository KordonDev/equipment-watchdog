<script lang="ts">
  import { Card } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import type { Member } from "./member.service";

  export let member: Member;
</script>

<Card class="m-4">
  <a href={`${routes.MemberDetail.link}${member.id}`} use:link>
    {`${member.id} ${member.name}`}
  </a>
</Card>
