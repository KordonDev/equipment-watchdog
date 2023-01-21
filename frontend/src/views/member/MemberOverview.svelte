<script lang="ts">
  import { Input, Label, Select } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import {
    getTranslatedGroups,
    groupFilter,
  } from "../../components/groupsStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { getMembers, type Member } from "./member.service";
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

  function byName(m: Member, search: string) {
    if (search === "") {
      return true;
    }
    return m.name.includes(search);
  }

  let search = "";

  let membersPromise = getMembers();
</script>

<Navigation />
<h1>Mitglieder</h1>
<a href={routes.AddMember.link} use:link>Mitglied hinzufügen</a>
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
<Label class="block mb-2">
  Suche
  <Input required class="mb-4" bind:value={search} />
</Label>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members
      .filter((m) => byGroup(m, groupFilterInternal))
      .filter((m) => byName(m, search)) as member}
      <MemberCard {member} columns={2} />
    {/each}
  </div>
{/await}
