<script lang="ts">
  import { getMembers } from "./member.service";

  let membersPromise = getMembers();
</script>

<h1>mitglieder</h1>
{#await membersPromise then members}
  <ul>
    {#each members as member}
      <li>
        {`${member.id} ${member.name}`}
      </li>
    {/each}
  </ul>
{/await}
