<script lang="ts">
  import { Label, Select } from "flowbite-svelte";
  import { group_outros } from "svelte/internal";
  import {
    getTranslatedGroups,
    groupFilter,
  } from "../../components/groupsStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { getMembers, Group, type Member } from "./member.service";
  import MemberCard from "./MemberCard.svelte";

  let allGroups = [];
  getTranslatedGroups.subscribe((groups) => {
    if (groups) {
      allGroups = [{ name: "Alle", value: "all" }, ...groups];
    }
  });

  let groupFilterInternal;
  groupFilter.subscribe((v) => (groupFilterInternal = v));
  function updateGroupFilter() {
    groupFilter.set(groupFilterInternal);
  }

  function byGroup(m: Member, filter?: string) {
    if (!filter || filter === "all") {
      return true;
    }
    return m.group === filter;
  }

  let membersPromise = getMembers();
</script>

<Navigation />
<h1>Mitglieder</h1>
<Label class="my-4" size="md">
  <div class="mb-2">Gruppe</div>
  <Select
    size="lg"
    required
    items={allGroups}
    bind:value={groupFilterInternal}
    on:change={updateGroupFilter}
  />
</Label>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members.filter((m) => byGroup(m, groupFilterInternal)) as member}
      <MemberCard {member} columns={2} />
    {/each}
  </div>
{/await}
