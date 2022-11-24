<script lang="ts">
  import { getMembers } from "./member.service";
  import MemberCard from "./MemberCard.svelte";

  let membersPromise = getMembers();
</script>

<h1>mitglieder</h1>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members as member}
      <MemberCard {member} />
    {/each}
  </div>
{/await}
