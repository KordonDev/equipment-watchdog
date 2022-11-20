<script lang="ts">
  import { getMembers } from "./member.service";
  import { Card } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";

  let membersPromise = getMembers();
</script>

<h1>mitglieder</h1>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members as member}
      <Card class="m-4">
        <a href={`${routes.MemberDetail.link}${member.id}`} use:link>
          {`${member.id} ${member.name}`}
        </a>
      </Card>
    {/each}
  </div>
{/await}
