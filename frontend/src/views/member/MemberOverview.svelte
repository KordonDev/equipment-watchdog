<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { getMembers } from "./member.service";
  import MemberCard from "./MemberCard.svelte";

  let membersPromise = getMembers();
</script>

<Navigation />
<h1>Mitglieder</h1>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members as member}
      <MemberCard {member} />
    {/each}
  </div>
{/await}
